package server

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/arsenalzp/keyvalstore/internal/server/errors"
	"github.com/arsenalzp/keyvalstore/internal/server/storage"
)

type Server struct {
	tlsConf        *tls.Config
	CrlPath        string
	ServerCrtData  []byte
	ServerKeyData  []byte
	RootCACertData []byte
	Nic            string
	IP             net.IP
	Address        string
	Port           int
	Storage        *storage.Storage
	lsnr           net.Listener
}

func (s *Server) Start() (net.Listener, error) {
	err := s.initNetwork()
	if err != nil {
		err = errors.New("unable to initialize network config", errors.SrvStartErr, err)
		return nil, err
	}

	if s.CrlPath == "" || s.ServerCrtData == nil || s.ServerKeyData == nil || s.RootCACertData == nil {
		err = errors.New("configuration error", errors.SrvStartErr, err)
		return nil, err
	}

	lsnr, err := net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.New("unable to create listener", errors.SrvStartErr, err)
		return nil, err
	}
	s.lsnr = lsnr

	s.tlsConf, err = s.initTLS()
	if err != nil {
		err = errors.New("unable to create SSL config", errors.SrvStartErr, err)
		return nil, err
	}

	log.Printf("starting server on port %d\n", s.Port)

	return s.lsnr, nil
}

func (s *Server) Stop() {
	if s.lsnr == nil {
		log.Println("stopping the server...")
		return
	}

	err := s.lsnr.Close()
	if err != nil {
		err = errors.New("unable to start server", errors.SrvStopErr, err)
		log.Fatalf("failed to stop the server: %s\n", err)
		return
	}

	log.Println("stopping the server...")
}

// initNetwork initialize network configuration
// based either on env variables or CLI arguments.
func (s *Server) initNetwork() error {
	// check NIC argument from command line
	if s.Nic != "" {
		ip, err := s.getIP(s.Nic)
		if err != nil {
			return errors.New("unable to initialize network", errors.NetworkInitErr, err)
		}
		s.IP = ip
	} else {
		s.IP = net.IPv4(0, 0, 0, 0)
	}

	// get SERVICE_PORT from env variable
	// else - use a port defined by command line argument
	if s.Port != 0 {
		s.Address = s.IP.String() + ":" + fmt.Sprint(s.Port)
	} else {
		s.Port = 6842
		s.Address = s.IP.String() + ":" + fmt.Sprint(s.Port)
	}

	return nil
}

// initTLS initialize TLS configration by loading server's certificate,
// key, CA certificate and CRL
func (s *Server) initTLS() (*tls.Config, error) {
	crt, err := tls.X509KeyPair(s.ServerCrtData, s.ServerKeyData)
	if err != nil {
		return nil, errors.New("unabel to load certificate or key file", errors.KeyCertLoadErr, err)
	}

	caPool, err := createCAPool(s.RootCACertData)
	if err != nil {
		return nil, errors.New("unabel to load root CA file", errors.CAcertLoadErr, err)
	}

	tlsConf := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{crt},
		ClientCAs:    caPool,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			crl, err := parseCRL(s.CrlPath)
			if err != nil {
				return err
			}

			err = checkCertWithCRL(verifiedChains[0][0], crl)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return tlsConf, err
}

// Check provided certificate cert against CRL
func checkCertWithCRL(cert *x509.Certificate, crl *pkix.CertificateList) error {
	for _, revokedCertificate := range crl.TBSCertList.RevokedCertificates {
		if revokedCertificate.SerialNumber.Cmp(cert.SerialNumber) == 0 {
			err := errors.New("certificate was revoked", errors.CRLCertRevokErr, nil)
			return err
		}
	}
	return nil
}

// Parse CRL file which is accessible by crlPath
func parseCRL(crlPath string) (*pkix.CertificateList, error) {
	crlData, err := os.ReadFile(crlPath)
	if err != nil {
		err := errors.New("unabel to read CRL file", errors.CRLLoadErr, err)
		return nil, err
	}

	crlList, err := x509.ParseCRL(crlData)
	if err != nil {
		err := errors.New("unabel to parse CRL data", errors.CRLParseErr, err)
		return nil, err
	}

	if crlList.TBSCertList.NextUpdate.Before(time.Now()) {
		err := errors.New("CRL is outdated", errors.CRLExpiredErr, nil)
		return nil, err
	}

	return crlList, nil
}

// retrieve IP address from the give NIC
func (s *Server) getIP(nic string) (net.IP, error) {
	iface, err := net.InterfaceByName(s.Nic) // does NIC exist ?
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return net.IP{}, err
	}

	ipaddrs, err := iface.Addrs() // get IP-address assigned to NIC
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return net.IP{}, err
	}

	if len(ipaddrs) == 0 {
		err = errors.New("network error, no configured IP address for the give NIC", errors.NetworkErrIpaddrs, nil)
		return net.IP{}, err
	}

	ip, _, err := net.ParseCIDR(ipaddrs[0].String())
	if err != nil {
		err = errors.New("network error", errors.NetworkErr, err)
		return net.IP{}, err
	}

	return ip, nil
}

func (s *Server) GetTlsConf() *tls.Config {
	return s.tlsConf
}

func createCAPool(rootCA []byte) (*x509.CertPool, error) {
	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(rootCA); !ok {
		return nil, errors.New("unable to append CAcert into CAPool", errors.CAPoolLoadErr, nil)
	}

	return caPool, nil
}
