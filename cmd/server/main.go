// Implements server-side of key-value storage with the import-export feature.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/arsenalzp/keyvalstore/cmd/server/server"
	"github.com/arsenalzp/keyvalstore/internal/server/errors"
	hndlr "github.com/arsenalzp/keyvalstore/internal/server/handler" // import handlers
	"github.com/arsenalzp/keyvalstore/internal/server/storage"
)

const helpMessage string = `
Usage of server:

The following environment variables are required:
CRL_PATH - path to a CRL file
SERVER_CERT - path to a server's certificate
SERVER_KEY - path to a server private key
ROOTCA_CERT - path to a root CA certificate
SERVICE_STORAGE - set an underlying storage (hash table or sqlite)
SERVICE_PORT - set TCP port to listen on
SERVICE_NIC - set NIC for binding
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	}
	flag.Parse()
}

func main() {
	// Read a server certificate
	serverCertData, err := os.ReadFile(os.Getenv("SERVER_CERT"))
	if err != nil {
		log.Fatal(err)
		return
	}

	// Read a server private key
	serverPrivKeyData, err := os.ReadFile(os.Getenv("SERVER_KEY"))
	if err != nil {
		log.Fatal(err)
		return
	}

	// Read CA certificate
	rootCACertData, err := os.ReadFile(os.Getenv("ROOTCA_CERT"))
	if err != nil {
		log.Fatal(err)
		return
	}

	var port int
	if stringPort, ok := os.LookupEnv("SERVICE_PORT"); ok {
		intPort, err := strconv.Atoi(stringPort)
		if err != nil {
			log.Fatal(err)
			return
		}
		port = intPort
	}

	srv := server.Server{
		CrlPath:        os.Getenv("CRL_PATH"),
		ServerCrtData:  serverCertData,
		ServerKeyData:  serverPrivKeyData,
		RootCACertData: rootCACertData,
		Nic:            os.Getenv("SERVICE_NIC"),
		Port:           port,
	}

	defer srv.Stop()

	strg, err := storage.NewStrg(os.Getenv("SERVICE_STORAGE"))
	if err != nil {
		log.Fatal(err) // followed by os.Exit(1)
	}

	lsnr, err := srv.Start()
	if err != nil {
		log.Fatal(err) // followed by os.Exit(1)
	}

	for {
		conn, err := lsnr.Accept()
		if err != nil {
			err = errors.New("network error", errors.NetworkCallErr, err)
			log.Println(err)
			continue
		}

		tlsConn := tls.Server(conn, srv.GetTlsConf())

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
