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
)

type Cmd = string

type dataStruct struct {
	strg.Storage
}

// Handle connection from a cli
func HandleCon(pCtx context.Context, con net.Conn, storage strg.Storage) {
	var mu sync.Mutex
	var buf []byte
	var ds = &dataStruct{storage}              // init new data structure
	var errCh chan error = make(chan error)    // channel to send errors
	var dataCh chan []byte = make(chan []byte) // channel to send a date

	var reader *bufio.Reader = bufio.NewReader(con)
	var writer *bufio.Writer = bufio.NewWriter(con)

	ctx, cancel := context.WithCancel(pCtx) // create context from the parent context

	defer cancel()
	defer con.Close()

	// Handle different requests withing a single connection
Loop:
	// continiously reading a data from the connection
	for {
		buf = make([]byte, 771) // initialize a buffer for incoming massage

		// read incoming bytes into initialized buffer
		i, err := io.ReadFull(con, buf)

		// probably, connection was closed by remote peer
		if err == io.EOF {
			return
		}

		if err != nil && err != io.ErrUnexpectedEOF {
			err = errors.New("handler error", errors.ReadClientErr, err)
			log.Printf("%+v", err)
			return
		}

		if i == 0 {
			continue Loop
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

			key := buf[3:256] // get key value from the buffer
			val := buf[256:]  // get value from the buffer

			go ds.set(ctx, key, val, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("set operation error", errors.SetOpTimeout, ctx.Err())
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
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

				copy(respBuf[1:], []byte(fmt.Sprint(err)))

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
			respBuf := make([]byte, 512) // create outcomming buffer

			key := buf[3:256]

			go ds.get(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("get operation error", errors.GetOpTimeout, ctx.Err())
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
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
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
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

			key := buf[3:256]

			go ds.del(ctx, key, dataCh, errCh)

			select {
			case <-ctx.Done():
				err := errors.New("del operation error", errors.DelOpTimeout, ctx.Err())
				copy(respBuf[1:], []byte(fmt.Sprint(err)))

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
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
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
				_, err := writer.Write(respBuf[:])
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
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				_, err = writer.Write(respBuf[:]) // return error to a client
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
				copy(respBuf[1:], []byte(fmt.Sprint(err))) // return error to client
				_, err = writer.Write(respBuf[:])          // return error to a client
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
				respBuf = append(respBuf, '\x00')
				mu.Lock()
				_, err := writer.Write(respBuf[:])
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

			// if the last byte equals '\x00' then no need to read more bytes
			// else read until reach byte equals '\x00'
			if buf[770] == '\x00' {
				data := buf[3:]

				go ds.imp(ctx, data, dataCh, errCh)

			} else {
				data, err := reader.ReadBytes('\x00')
				if err != nil {
					err := errors.New("import operation error", errors.ReadClientErr, err)
					copy(respBuf[1:], []byte(fmt.Sprint(err)))
					writer.Write(respBuf[:])

					return
				}

				buf = append(buf[3:], data...)

				go ds.imp(ctx, buf, dataCh, errCh)
			}

			select {
			case <-ctx.Done():
				err = errors.New("import operation error", errors.ImpOpTimeout, ctx.Err())
				log.Printf("%+v", err)
				copy(respBuf[1:], []byte(fmt.Sprint(err)))
				_, err = writer.Write(respBuf[:]) // return error to a client
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
				copy(respBuf[1:], []byte(fmt.Sprintf("%s", err)))
				_, err = writer.Write(respBuf[:]) // return error to a client
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
				_, err := writer.Write(respBuf[:])
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
