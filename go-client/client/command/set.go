// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

// Set key and value pair
func Set(con net.Conn, dataChan chan<- struct{}, errChan chan<- error, key string, value string) {
	var buf [MESSAGE_SIZE]byte // command 3B, key 256B, value 512B

	// Check input, neither key nor value should be empty
	if key == "" {
		err := errors.New("set command error", errors.KeyEmptyErr, nil)
		errChan <- err
	} else if value == "" {
		err := errors.New("set command error", errors.ValueEmptyErr, nil)
		errChan <- err
		return
	}

	writer := bufio.NewWriter(con) // connection writer to send the data to the server

	copy(buf[0:3], []byte(SET))
	copy(buf[3:259], []byte(key))
	copy(buf[259:], []byte(value))
	buf[771] = EOT

	_, err := writer.Write(buf[:]) // write command, key and val
	if err != nil {
		err = errors.New("set operation error", errors.WriteServerErr, err)
		errChan <- err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("set operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT)
	if err != nil {
		err = errors.New("set operation error", errors.ReadServerErr, err)
		errChan <- err
		return
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:]) // retrieve error value from the server response
		err = errors.New("set operation error", errors.SetServerRespErr, err)
		errChan <- err
		return
	}

	dataChan <- struct{}{}
}
