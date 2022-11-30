// Implements server-side of key-value storage with the import-export feature.

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"gokeyval/internal/server/errors"
	hndlr "gokeyval/internal/server/handler" // import handlers
	"gokeyval/internal/server/storage"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

const helpMessage string = `
Usage of server:
-i define network interface to bind to
-p define a port to listen on (default: 6842)
-s define underlying storage

The following environment variables are required:
CRL_PATH - path to a CRL file
SERVER_CERT - path to a server's certificate
SERVER_KEY - path to a server private key
ROOTCA_CERT - root CA certificate
`

var (
	// Variables are used for network initialization
	// either env variables or CLI arguments
	_ipaddr, _nic, _addr string
	// these variables can be set by environment variables
	_SERVICE_PORT, _SERVICE_STORAGE string
	// these variables can be set by environment variables
	_SERVER_CERT, _SERVER_KEY, _ROOTCA_CERT, _CRL_path string

	_tlsConf tls.Config
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	}
	flag.StringVar(&_nic, "i", "", "define network interface to bind to")
	flag.StringVar(&_SERVICE_PORT, "p", "6842", "define a port to listen on")
	flag.StringVar(&_SERVICE_STORAGE, "s", "", "define underlying storage (hash for hash table, sqlite for SQLite)")
	flag.Parse()

	initTLS()     // 1) initialize TLS network config
	initNetwork() // 2) initialize network config
	initStorage() // 3) initialize underlying storage
}

// initTLS initialize TLS configration by loading server's certificate,
// key, CA certificate and CRL.
func initTLS() {
	if crl_path, ok := os.LookupEnv("CRL_PATH"); !ok {
		_CRL_path = filepath.Join("cert", "list.crl")
	} else {
		_CRL_path = crl_path
	}

	if cert_path, ok := os.LookupEnv("SERVER_CERT"); !ok {
		_SERVER_CERT = filepath.Join("cert", "tls.crt")
	} else {
		_SERVER_CERT = cert_path
	}

	if key_path, ok := os.LookupEnv("SERVER_KEY"); !ok {
		_SERVER_KEY = filepath.Join("cert", "tls.key")
	} else {
		_SERVER_KEY = key_path
	}

	if rootca_path, ok := os.LookupEnv("ROOTCA_CERT"); !ok {
		_ROOTCA_CERT = filepath.Join("cert", "rootCA.crt")
	} else {
		_ROOTCA_CERT = rootca_path
	}

	crt, err := tls.LoadX509KeyPair(_SERVER_CERT, _SERVER_KEY)
	if err != nil {
		err := errors.New("unabel to load certificate or key file", errors.KeyCertLoadErr, err)
		log.Fatal(err)
	}

	rootCA, err := ioutil.ReadFile(_ROOTCA_CERT)
	if err != nil {
		err := errors.New("unabel to load root CA file", errors.CAcertLoadErr, err)
		log.Fatal(err)
	}

	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(rootCA); !ok {
		err := errors.New("unabel to add root CA into the pool", errors.CAPoolLoadErr, nil)
		log.Fatal(err)
	}

	_tlsConf = tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{crt},
		ClientCAs:    caPool,
		ServerName:   "server.example.com",
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			err := checkCertWithCRL(verifiedChains[0][0], _CRL_path)
			if err != nil {
				return err
			}
			return nil
		},
	}
}

// initNetwork initialize network configuration
// based either on env variables or CLI arguments.
func initNetwork() {
	var err error
	// check NIC argument from command line
	if _nic != "" {
		_ipaddr, err = getIP(_nic)
		if err != nil {
			log.Fatal(err)
		}
	}

	// get SERVICE_PORT from env variable
	// else - use a port defined by command line argument
	if value, isExist := os.LookupEnv("SERVICE_PORT"); isExist {
		_SERVICE_PORT = value
		_addr = _ipaddr + ":" + _SERVICE_PORT
		return
	} else {
		_addr = _ipaddr + ":" + _SERVICE_PORT
	}

}

// initStorage initilize underlaying storage
// use the following: either hash table or sqlite3
func initStorage() {
	if _SERVICE_STORAGE != "" {
		return
	}

	if value, isExist := os.LookupEnv("SERVICE_STORAGE"); isExist {
		_SERVICE_STORAGE = value
		return
	}

	err := errors.New("SERVICE_STORAGE env isn't defined", errors.StorageInitErr, nil)
	log.Fatal(err)
}

// retrieve IP address from the give NIC
func getIP(nic string) (string, error) {
	ifi, err := net.InterfaceByName(nic) // does NIC exist ?
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return "", err
	}

	ipaddrs, err := ifi.Addrs() // get IP-address assigned to NIC
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return "", err
	}

	if len(ipaddrs) == 0 {
		err = errors.New("network error, no configured IP address for the give NIC", errors.NetworkErrIpaddrs, nil)
		return "", err
	}

	ip, _, err := net.ParseCIDR(ipaddrs[0].String())
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return "", err
	}

	return ip.String(), nil
}

// Function checkCertWithCRL validate client's certificate:
// - was it signed by known CA or
// - was it revoked ?
func checkCertWithCRL(cert *x509.Certificate, crlPath string) error {
	// Parse CRL file which is accessible by crlPath
	crl, err := parseCRL(crlPath)
	if err != nil {
		err := errors.New("certificate validation error", errors.CRLValidErr, err)
		return err
	}

	// Check provided certificate cert against CRL
	for _, revokedCertificate := range crl.TBSCertList.RevokedCertificates {
		if revokedCertificate.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			err := errors.New("certificate was revoked", errors.CRLCertRevokErr, nil)
			return err
		}
	}
	return nil
}

func parseCRL(crlPath string) (*pkix.CertificateList, error) {
	var crl *pkix.CertificateList

	f, err := os.Open(crlPath)
	if err != nil {
		err := errors.New("unabel to load CRL file", errors.CRLOpenErr, err)
		return nil, err
	}

	defer f.Close()

	crlData, err := ioutil.ReadAll(f)
	if err != nil {
		err := errors.New("unabel to load CRL file", errors.CRLLoadErr, err)
		return nil, err
	}

	crl, err = x509.ParseCRL(crlData)
	if err != nil {
		err := errors.New("unabel to parse CRL data", errors.CRLParseErr, err)
		return nil, err
	}

	if crl.TBSCertList.NextUpdate.Before(time.Now()) {
		err := errors.New("CRL is outdated", errors.CRLExpiredErr, nil)
		return nil, err
	}

	return crl, nil
}

func main() {
	strg, err := storage.NewStrg(_SERVICE_STORAGE)
	if err != nil {
		err = errors.New("undefined storage type", errors.StorageKindUndef, err)
		log.Fatal(err) // followed by os.Exit(1)
	}

	ln, err := net.Listen("tcp", _addr)
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		log.Fatal(err) // followed by os.Exit(1)
	}

	log.Printf("service is listening on port %s...\n", _SERVICE_PORT)

	for {
		conn, err := ln.Accept()
		if err != nil {
			err = errors.New("network error", errors.NetworkCallErr, err)
			log.Println(err)
			continue
		}

		tlsConn := tls.Server(conn, &_tlsConf)
		//tlsConn.SetDeadline(time.Now().Add(timeoutIO))

		err = tlsConn.Handshake()
		if err != nil {
			err = errors.New("network error", errors.NetworkCallErr, err)
			log.Println(err)
			continue
		}

		ctx := context.Background()
		if err != nil {
			err = errors.New("network error", errors.NetworkCallErr, err)
			log.Println(err)
		}

		go hndlr.HandleCon(ctx, tlsConn, strg)
	}
}
