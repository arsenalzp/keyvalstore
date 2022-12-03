// Implements server-side of key-value storage with the import-export feature.

package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"gokeyval/cmd/server/server"
	"gokeyval/internal/server/errors"
	hndlr "gokeyval/internal/server/handler" // import handlers
	"gokeyval/internal/server/storage"
	"log"
	"os"
)

const helpMessage string = `
Usage of server:

The following environment variables are required:
CRL_PATH - path to a CRL file
SERVER_CERT - path to a server's certificate
SERVER_KEY - path to a server private key
ROOTCA_CERT - path to a root CA certificate
SERVICE_STORAGE - set an underlying storage (hash table or sqlite)
SERVICE_NIC - set NIC for binding
SERVICE_PORT - set TCP port to listen on
`

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpMessage)
	}
	flag.Parse()
}

func main() {
	srv := server.Server{
		CrlPath:     os.Getenv("CRL_PATH"),
		CrtPath:     os.Getenv("SERVER_CERT"),
		PrivkeyPath: os.Getenv("SERVER_KEY"),
		RootcaPath:  os.Getenv("ROOTCA_CERT"),
		Nic:         os.Getenv("SERVICE_NIC"),
		Port:        os.Getenv("SERVICE_PORT"),
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
