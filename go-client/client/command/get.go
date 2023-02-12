// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

// Get a value for a give key
func Get(con net.Conn, dataChan chan<- []byte, errChan chan<- error, key string) {
	var buf [MESSAGE_SIZE]byte

	writer := bufio.NewWriter(con) // connection writer to send the data to the server

	copy(buf[0:3], []byte(GET))
	copy(buf[3:259], []byte(key))
	buf[771] = EOT

	_, err := writer.Write(buf[:]) // write command, key and val
	if err != nil {
		err = errors.New("get operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("get operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT) // waiting for server response
	if err != nil {
		err = errors.New("get operation failed", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:]) // retrieve error value from the server response
		err = errors.New("get operation error", errors.GetServerRespErr, err)
		errChan <- err
		return
	}

	respBuf = bytes.TrimRight(respBuf[1:], string(EOT))
	respBuf = bytes.TrimRight(respBuf, "\x00")

	dataChan <- respBuf
}
