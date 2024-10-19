// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/arsenalzp/keyvalstore/internal/server/errors"
	strg "github.com/arsenalzp/keyvalstore/internal/server/storage"

	"sync"
)

const (
	timeoutOp = 10 * time.Second // timeout for a storage operations
	OK        = 'O'
	NOK       = 'N'
	EOT       = '\u0004'
)

type Cmd = string

type Handler interface {
	HandleCon(context.Context, net.Conn, strg.Storage)
}

type dataStruct struct {
	strg.Storage
}

// Handle connection from a cli
func HandleCon(pCtx context.Context, con net.Conn, storage strg.Storage) {
	var mu sync.Mutex
	var ds = &dataStruct{storage}              // init new data structure
	var errCh chan error = make(chan error)    // channel to send errors
	var dataCh chan []byte = make(chan []byte) // channel to send a date
	var reader *bufio.Reader
	var writer *bufio.Writer

	ctx, cancel := context.WithCancel(pCtx) // create context from the parent context

	defer func() {
		if err := recover(); err != nil {
			log.Printf("%+v", err)
		}
	}()
	defer cancel()
	defer con.Close()

	// Handle different requests withing a single connection
Loop:
	// continiously reading a data from the connection
	for {
		reader = bufio.NewReader(con)
		writer = bufio.NewWriter(con)

		// read a data from the connection, until EOT reached
		buf, err := reader.ReadBytes(EOT)
		if err == io.EOF { // probably, connection was closed by remote peer
			return
		}

		if err != nil && err != io.ErrUnexpectedEOF {
			err = errors.New("handler error", errors.ReadClientErr, err)
			log.Printf("%+v", err)
			return
		}

		// get command from the buffer
		cmd := getCmd(buf)
		if cmd == "" {
			continue Loop
		}

		// select the command
		switch cmd {
		case "set":
			respBuf := make([]byte, 63) // create outcomming buffer
			key := readKey(buf)         // get key value from the buffer
			val := readValue(buf)       // get value from the buffer

			go ds.set(ctx, key, val, dataCh, errCh)

			select {
			case <-ctx.Done():
				respBuf = writeStatus(respBuf, NOK)

				err := errors.New("set operation error", errors.SetOpTimeout, ctx.Err())
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("set operation error", errors.SettOpErr, err)
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case <-dataCh:
				respBuf = writeStatus(respBuf, OK)
				respBuf = writeEOT(respBuf)

				err := sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				continue Loop // exit
			}

		case "get":
			respBuf := make([]byte, 512)

			// get key value from the buffer
			key := readKey(buf)

			go ds.get(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				respBuf = writeStatus(respBuf, NOK)

				err := errors.New("get operation error", errors.GetOpTimeout, ctx.Err())
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("get operation error", errors.GetOpErr, err)
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case data := <-dataCh:
				respBuf = writeStatus(respBuf, OK)
				respBuf = writeValue(respBuf, data)
				respBuf = writeEOT(respBuf)

				mu.Lock()
				err := sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}
				mu.Unlock()

				continue Loop
			}

		case "del":
			respBuf := make([]byte, 63) // create outcomming buffer
			// get key value from the buffer
			key := readKey(buf)

			go ds.del(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				respBuf = writeStatus(respBuf, NOK)

				err := errors.New("del operation error", errors.DelOpTimeout, ctx.Err())
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("del operation error", errors.DelOpErr, err)
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case <-dataCh:
				respBuf = writeStatus(respBuf, OK)
				respBuf = writeEOT(respBuf)

				err := sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				continue Loop
			}
		case "exp":
			respBuf := make([]byte, 1)

			go ds.exp(ctx, dataCh, errCh)

			select {
			case <-ctx.Done():
				respBuf = writeStatus(respBuf, NOK)

				err := errors.New("export operation error", errors.ExpOpTimeout, ctx.Err())
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("export operation error", errors.ExpOpErr, err)
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case data := <-dataCh:
				respBuf = writeStatus(respBuf, OK)

				respBuf = writeExport(respBuf, data)
				respBuf = writeEOT(respBuf)

				mu.Lock()
				err := sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}
				mu.Unlock()

				continue Loop
			}

		case "imp":
			respBuf := make([]byte, 64)

			importData := readImport(buf[3:])

			go ds.imp(ctx, importData, dataCh, errCh)

			select {
			case <-ctx.Done():
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("import operation error", errors.ImpOpTimeout, ctx.Err())
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				respBuf = writeStatus(respBuf, NOK)

				err = errors.New("import operation error", errors.ImpOpErr, err)
				respBuf = writeError(respBuf, err)

				err = sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				log.Printf("%+v", err)

			case <-dataCh:
				respBuf = writeStatus(respBuf, OK)
				respBuf = writeEOT(respBuf)

				err := sendData(respBuf, *writer)
				if err != nil {
					err = errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}

				continue Loop
			}
		}
	}
}

func getCmd(buf []byte) Cmd {
	cmd := string(buf[0:3])

	switch cmd {
	case "set":
		return cmd
	case "get":
		return cmd
	case "del":
		return cmd
	case "imp":
		return cmd
	case "exp":
		return cmd
	default:
		return ""
	}
}

func readKey(buf []byte) []byte {
	return buf[3:259]
}

func readValue(buf []byte) []byte {
	return buf[259:512]
}

func readImport(buf []byte) []byte {
	return trimEOT(buf)
}

func writeValue(respBuf, data []byte) []byte {
	copy(respBuf[1:], data)
	return respBuf
}

func writeExport(respBuf, export []byte) []byte {
	return append(respBuf, export...)

}

func writeError(respBuf []byte, err error) []byte {
	copy(respBuf[1:], []byte(fmt.Sprint(err)))
	respBuf[63] = EOT
	return respBuf[0:64]
}

func writeStatus(respBuf []byte, status rune) []byte {
	respBuf[0] = byte(status) // the first byte indicates type of the message
	return respBuf
}

func writeEOT(respBuf []byte) []byte {
	return append(respBuf, EOT)
}

func sendData(respBuf []byte, writer bufio.Writer) error {
	_, err := writer.Write(respBuf[:])
	if err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func trimEOT(trimData []byte) []byte {
	return bytes.Trim(bytes.Trim(trimData, "\x00"), string(EOT))
}
