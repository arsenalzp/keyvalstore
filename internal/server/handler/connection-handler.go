// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bufio"
	"context"
	"fmt"
	"gokeyval/internal/server/errors"
	strg "gokeyval/internal/server/storage"
	"io"
	"log"
	"net"
	"time"

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
		cmd := getCmd(buf[0:3])
		if cmd == "" {
			continue Loop
		}

		// select the command
		switch cmd {
		case "set":
			respBuf := make([]byte, 64) // create outcomming buffer
			key := buf[3:259]           // get key value from the buffer
			val := buf[259:512]         // get value from the buffer

			go ds.set(ctx, key, val, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("set operation error", errors.SetOpTimeout, ctx.Err())
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				log.Println(err)

			case err := <-errCh:
				err = errors.New("set operation error", errors.SettOpErr, err)
				log.Printf("%+v", err)
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

			case <-dataCh:
				respBuf[0] = OK // the first byte indicates type of the message
				respBuf = append(respBuf, EOT)
				_, err := writer.Write(respBuf[:])
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("set operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				continue Loop // exit
			}

		case "get":
			respBuf := make([]byte, 513)

			// get key value from the buffer
			key := buf[3:259]

			go ds.get(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("get operation error", errors.GetOpTimeout, ctx.Err())
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[512] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				err = errors.New("get operation error", errors.GetOpErr, err)
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[512] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				log.Printf("%+v", err)

			case data := <-dataCh:
				respBuf[0] = OK // the first byte indicates type of the message
				copy(respBuf[1:], data)
				respBuf = append(respBuf, EOT)
				mu.Lock()
				_, err := writer.Write(respBuf)
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("get operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)
					return
				}
				mu.Unlock()
				continue Loop
			}

		case "del":
			respBuf := make([]byte, 64) // create outcomming buffer

			// get key value from the buffer
			key := buf[3:259]

			go ds.del(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("del operation error", errors.DelOpTimeout, ctx.Err())
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				log.Printf("%+v", err)

			case err := <-errCh:
				err = errors.New("del operation error", errors.DelOpErr, err)
				respBuf[0] = NOK
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf[:]) // return error to a client
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				log.Printf("%+v", err)

			case <-dataCh:
				respBuf[0] = OK // the first byte indicates type of the message
				respBuf[63] = EOT
				_, err := writer.Write(respBuf)
				if err != nil {
					err = errors.New("del operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
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
				err := errors.New("export operation error", errors.ExpOpTimeout, ctx.Err())
				respBuf[0] = NOK
				copy(respBuf[1:63], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf) // return error to a client
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err := errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}
				log.Printf("%+v", err)

			case err := <-errCh:
				err = errors.New("export operation error", errors.ExpOpErr, err)
				log.Printf("%+v", err)
				respBuf[0] = NOK
				copy(respBuf[1:63], []byte(fmt.Sprint(err))) // return error to client
				respBuf[63] = EOT
				_, err = writer.Write(respBuf) // return error to a client
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

			case data := <-dataCh:
				respBuf[0] = OK // the first byte indicates type of the message
				respBuf = append(respBuf, data...)
				respBuf = append(respBuf, EOT)
				mu.Lock()
				_, err := writer.Write(respBuf)
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
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

			go ds.imp(ctx, buf[3:], dataCh, errCh)

			select {
			case <-ctx.Done():
				err = errors.New("import operation error", errors.ImpOpTimeout, ctx.Err())
				log.Printf("%+v", err)
				respBuf[0] = NOK
				copy(respBuf[1:63], []byte(fmt.Sprint(err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf) // return error to a client
				if err != nil {
					err = errors.New("export operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err := errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

			case err := <-errCh:
				err = errors.New("import operation error", errors.ImpOpErr, err)
				log.Printf("%+v", err)
				respBuf[0] = NOK
				copy(respBuf[1:63], []byte(fmt.Sprintf("%s", err)))
				respBuf[63] = EOT
				_, err = writer.Write(respBuf) // return error to a client
				if err != nil {
					err = errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
				if err != nil {
					err = errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

			case <-dataCh:
				respBuf[0] = OK // the first byte indicates type of the message
				respBuf = append(respBuf, EOT)
				respBuf[63] = EOT
				_, err := writer.Write(respBuf)
				if err != nil {
					err = errors.New("import operation error", errors.WriteClientErr, err)
					log.Printf("%+v", err)

					return
				}

				err = writer.Flush()
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
	cmd := string(buf)

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
