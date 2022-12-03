package server

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"gokeyval/internal/server/errors"
	"gokeyval/internal/server/storage"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	tlsConf     *tls.Config
	CrlPath     string
	CrtPath     string
	PrivkeyPath string
	RootcaPath  string
	Nic         string
	IP          net.IP
	Address     string
	Port        string
	Storage     *storage.Storage
	lsnr        net.Listener
}

func (s *Server) Start() (net.Listener, error) {
	err := s.initNetwork()
	if err != nil {
		err = errors.New("unable to start server", errors.SrvStartErr, err)
		return nil, err
	}

	if s.CrtPath == "" || s.PrivkeyPath == "" || s.RootcaPath == "" || s.CrlPath == "" {
		lsnr, err := net.Listen("tcp", s.Address)
		if err != nil {
			err = errors.New("unable to start server", errors.SrvStartErr, err)
			return nil, err
		}
		s.lsnr = lsnr
		return s.lsnr, nil
	}

	lsnr, err := net.Listen("tcp", s.Address)
	if err != nil {
		err = errors.New("unable to start server", errors.SrvStartErr, err)
		return nil, err
	}
	s.lsnr = lsnr

	err = s.initTLS()
	if err != nil {
		err = errors.New("unable to start server", errors.SrvStartErr, err)
		return nil, err
	}

	log.Printf("starting server on port %s\n", s.Port)

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

	if s.Port != "" {
		s.Address = s.IP.String() + ":" + s.Port
	} else {
		s.Port = "6842"
		s.Address = ":" + s.Port
	}

	return nil
}

// initTLS initialize TLS configration by loading server's certificate,
// key, CA certificate and CRL
func (s *Server) initTLS() error {
	crt, err := tls.LoadX509KeyPair(s.CrtPath, s.PrivkeyPath)
	if err != nil {
		return errors.New("unabel to load certificate or key file", errors.KeyCertLoadErr, err)
	}

	rootCA, err := ioutil.ReadFile(s.RootcaPath)
	if err != nil {
		return errors.New("unabel to load root CA file", errors.CAcertLoadErr, err)
	}

	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(rootCA); !ok {
		return errors.New("unabel to add root CA into the pool", errors.CAPoolLoadErr, nil)
	}

	s.tlsConf = &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{crt},
		ClientCAs:    caPool,
		ServerName:   "server.example.com",
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			err := checkCertWithCRL(verifiedChains[0][0], s.CrlPath)
			if err != nil {
				return err
			}
			return nil
		},
	}

	return nil
}

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
