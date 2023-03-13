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
	"github.com/arsenalzp/keyvalstore/go-client/internal/util"
)

type Client struct {
	conn *tls.Conn
	mux  sync.Mutex
}

type ClientConfig struct {
	RootCAPath      string
	CertificatePath string
	PrivateKeyPath  string
	Port            uint16
	Address         string
}

// Connect to a server. Connect returns *Client structure
func (c *ClientConfig) Connect() (*Client, error) {
	tlsConfig, err := c.initTLS()
	if err != nil {
		err = errors.New("connection error", errors.NetworkErr, err)
		return nil, err
	}

	// Call to a server
	conn, err := net.Dial("tcp", c.Address+":"+fmt.Sprint(c.Port))
	if err != nil {
		err = errors.New("connection error", errors.NetworkErr, err)
		return nil, err
	}

	// Initialize TLS connection for the given connection and TLS Config
	tlsConn := tls.Client(conn, tlsConfig)

	err = tlsConn.Handshake()
	if err != nil {
		err = errors.New("TLS handshake error", errors.NetworkErr, err)
		conn.Close()
		return nil, err
	}

	clientConnection := &Client{
		conn: tlsConn,
	}
	return clientConnection, nil
}

// Close connection with a server
func (c *Client) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}

// Get a value for a given key. Get returns []byte or error in case of failure
func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	// validate the key data parameter
	err := util.ValidateInput(key, "")
	if err != nil {
		return nil, err
	}

	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	c.mux.Lock()
	defer c.mux.Unlock()

	go cmd.Get(c.conn, dataChan, errChan, key)

	select {
	case <-ctx.Done():
		err := errors.New("get command interrupted", errors.GetCancelErr, ctx.Err())
		return nil, err
	case val := <-dataChan:
		return val, nil
	case err := <-errChan:
		fmt.Printf("error to retrieve key %s: %s\n", key, err)
		return nil, err
	}
}

// Save key=value pair on a server. Set return error in case of failure
func (c *Client) Set(ctx context.Context, key, value string) error {
	// validate the key and the value data parameters
	err := util.ValidateInput(key, value)
	if err != nil {
		return err
	}

	dataChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)

	c.mux.Lock()
	defer c.mux.Unlock()

	go cmd.Set(c.conn, dataChan, errChan, key, value)

	select {
	case <-ctx.Done():
		err := errors.New("get command interrupted", errors.SetCancelErr, ctx.Err())
		return err
	case <-dataChan:
		return nil
	case err := <-errChan:
		return err
	}
}

// Delete key=value pair on a server. Del returns error in case of failure
func (c *Client) Del(ctx context.Context, key string) error {
	// validate the key and the value data parameters
	err := util.ValidateInput(key, "")
	if err != nil {
		return err
	}

	dataChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)

	c.mux.Lock()
	defer c.mux.Unlock()

	go cmd.Del(c.conn, dataChan, errChan, key)

	select {
	case <-ctx.Done():
		err := errors.New("get operation interrupted", errors.DelCancelErr, ctx.Err())
		return err
	case <-dataChan:
		return nil
	case err := <-errChan:
		return err
	}
}

// Import key=value pairs into a server. Import returns error in case of failure
func (c *Client) Import(ctx context.Context, data []byte) error {
	dataChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)

	err := util.ValidateData([]byte(data)) // validate data befor sending
	if err != nil {
		err = errors.New("import operation failed: validation of input failed", errors.ReadStdinErr, err)
		return err
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	go cmd.Import(c.conn, dataChan, errChan, data)

	select {
	case <-ctx.Done():
		err := errors.New("import operation interrupted", errors.DelCancelErr, ctx.Err())
		return err
	case <-dataChan:
		return nil
	case err := <-errChan:
		return err
	}
}

// Export key=value pairs from a server. Export returns []byte or error in case of failure
func (c *Client) Export(ctx context.Context) ([]byte, error) {
	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	c.mux.Lock()
	defer c.mux.Unlock()

	go cmd.Export(c.conn, dataChan, errChan)

	select {
	case <-ctx.Done():
		err := errors.New("get operation interrupted", errors.DelCancelErr, ctx.Err())
		return nil, err
	case data := <-dataChan:
		return data, nil
	case err := <-errChan:
		return nil, err
	}
}

// Initialize TLS Config. initTLS returns *tls.Config or error in case of failure
func (c *ClientConfig) initTLS() (*tls.Config, error) {
	crt, err := tls.LoadX509KeyPair(c.CertificatePath, c.PrivateKeyPath)
	if err != nil {
		return nil, errors.New("unabel to load certificate or key", errors.NetworkErr, err)
	}

	rootCA, err := ioutil.ReadFile(c.RootCAPath)
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
