// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

// Delete a value for a given key
func Del(con net.Conn, dataChan chan<- struct{}, errChan chan<- error, key string) {
	var buf [MESSAGE_SIZE]byte

	writer := bufio.NewWriter(con) // connection writer to send the data to the server

	copy(buf[0:3], []byte(DELETE))
	copy(buf[3:], []byte(key))
	buf[771] = EOT

	_, err := writer.Write(buf[:])
	if err != nil {
		err = errors.New("del operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("del operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT) // waiting for server response
	if err != nil {
		err = errors.New("del operation failed", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:]) // retrieve error value from the server response
		err = errors.New("del operation error", errors.DelServerRespErr, err)
		errChan <- err
		return
	}

	dataChan <- struct{}{}
}
