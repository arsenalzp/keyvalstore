// Package implements CLI commands.

package command

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gokeyval/internal/cli/errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	MESSAGE_SIZE = 772
	EOT          = '\u0004' // End-Of-Trasmission character
	DELETE       = "del"
	GET          = "get"
	SET          = "set"
	EXPORT       = "exp"
	IMPORT       = "imp"
)

var serverAddress string
var client_cert, privkey_cert, rootca_cert string
var tlsConf tls.Config

var rootCmd = &cobra.Command{
	Use: `
	keyval get [--server] [--key] [--cert] [--CAcert] key | 
	set [--server] [--key] [--cert] [--CAcert] key=val | 
	del [--server] [--key] [--cert] [--CAcert] key | 
	export [--server] [--key] [--cert] [--CAcert] |
	import [--server] [--key] [--cert] [--CAcert] JSON
	`,
	Short: "Keyval is fast Unix-style key=val storage",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("main command")
	},
	SilenceUsage: true,
}

func CreateConnection() (*tls.Conn, error) {
	//var ip string
	var port string

	addrArg := strings.Trim(serverAddress, "\n\t") // argument from stdin

	addr, port, err := net.SplitHostPort(addrArg)
	if err != nil {
		err = errors.New("connection error", errors.InvalidAddrErr, err)
		return nil, err
	}

	ip, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		err = errors.New("connection error", errors.InvalidAddrErr, err)
		return nil, err
	}

	if len(port) == 0 {
		port = "6842"
	}

	err = initTLS() // initialize TLS configuration
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", ip.IP.String()+":"+port)
	if err != nil {
		err = errors.New("connection error", errors.NetworkErr, err)
		return nil, err
	}

	tlsConn := tls.Client(conn, &tlsConf)

	// set timeout
	timeOut := time.Now()
	tlsConn.SetWriteDeadline(timeOut.Add(time.Second * 20))
	tlsConn.SetReadDeadline(timeOut.Add(time.Second * 20))

	err = tlsConn.Handshake()
	if err != nil {
		log.Printf("TLS handshake error: %s\n", err)
		conn.Close()
		return nil, err
	}

	return tlsConn, nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initTLS() error {
	crt, err := tls.LoadX509KeyPair(client_cert, privkey_cert)
	if err != nil {
		return errors.New("unabel to load certificate or key", errors.NetworkErr, err)
	}

	rootCA, err := ioutil.ReadFile(rootca_cert)
	if err != nil {
		return errors.New("unabel to load root CA file", errors.NetworkErr, err)
	}

	caPool := x509.NewCertPool()

	if ok := caPool.AppendCertsFromPEM(rootCA); !ok {
		log.Println()
		return errors.New("unabel to add root CA into the certificate pool", errors.NetworkErr, nil)
	}

	tlsConf = tls.Config{
		MinVersion:         tls.VersionTLS13,
		Certificates:       []tls.Certificate{crt},
		ClientCAs:          caPool,
		InsecureSkipVerify: true, // FOR TEST PURPOSES ONLY, CONSIDER NEW CERTIFICATE CRETION WITH THE !!SAN!!
	}

	return nil
}
