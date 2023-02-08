package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"

	cmd "github.com/arsenalzp/keyvalstore/go-client/client/command"
	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

type Client struct {
	conn *tls.Conn
	mux  sync.Mutex
}

type ClientConfig struct {
	RootcaPath  string
	CrtPath     string
	PrivkeyPath string
	IP          net.IP
	Port        uint16
	Address     string
	tlsConf     *tls.Config
}

func (c *ClientConfig) Connect() (*Client, error) {
	tlsConfig, err := c.initTLS()
	if err != nil {
		err = errors.New("connection error", errors.NetworkErr, err)
		return nil, err
	}

	conn, err := net.Dial("tcp", c.IP.String()+":"+fmt.Sprint(c.Port))
	if err != nil {
		err = errors.New("connection error", errors.NetworkErr, err)
		return nil, err
	}

	tlsConn := tls.Client(conn, tlsConfig)

	err = tlsConn.Handshake()
	if err != nil {
		err = errors.New("TLS handshake error", errors.NetworkErr, err)
		conn.Close()
		return nil, err
	}

	clientConnection := &Client{
		conn: tlsConn,
		mux:  sync.Mutex{},
	}
	return clientConnection, nil
}

func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, k string) (string, error) {
	dataCh := make(chan []byte, 1)
	c.mux.Lock()
	defer c.mux.Unlock()
	go cmd.Get(c.conn, dataCh, k)

	select {
	case <-ctx.Done():
		return "", nil
	case val := <-dataCh:
		return string(val), nil
	}
}

func (c *Client) Set(ctx context.Context, k, v string) {

}

func (c *Client) Del(ctx context.Context, k string) {

}

func (c *Client) Import(ctx context.Context) {

}

func (c *Client) Export(ctx context.Context) {

}

func (c *ClientConfig) initTLS() (*tls.Config, error) {
	crt, err := tls.LoadX509KeyPair(c.CrtPath, c.PrivkeyPath)
	if err != nil {
		return nil, errors.New("unabel to load certificate or key", errors.NetworkErr, err)
	}

	rootCA, err := ioutil.ReadFile(c.RootcaPath)
	if err != nil {
		return nil, errors.New("unabel to load root CA file", errors.NetworkErr, err)
	}

	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(rootCA); !ok {
		log.Println()
		return nil, errors.New("unabel to add root CA into the certificate pool", errors.NetworkErr, nil)
	}

	tlsConf := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		Certificates:       []tls.Certificate{crt},
		ClientCAs:          caPool,
		InsecureSkipVerify: true, // FOR TEST PURPOSES ONLY, CONSIDER NEW CERTIFICATE CRETION WITH THE !!SAN!!
	}

	return tlsConf, nil
}
