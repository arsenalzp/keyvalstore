// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

// Export key=value pairs from a server
func Export(con net.Conn, dataChan chan<- []byte, errChan chan<- error) {
	var buf []byte = make([]byte, 3)

	writer := bufio.NewWriter(con) // connection writer to send the data to the server

	copy(buf[0:3], []byte("exp"))
	buf = append(buf, EOT) // add EOT to signal the end of transmission

	_, err := writer.Write(buf)
	if err != nil {
		err = errors.New("export operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("export operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT)
	if err != nil {
		err = errors.New("export operation error", errors.ReadServerErr, err)
		errChan <- err
		return
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:]) // retrieve error value from the server response
		err = errors.New("export operation error", errors.ExpServerRespErr, err)
		errChan <- err
		return
	}

	// Trim response buffer: delete NULL and EOT bytes
	respBuf = bytes.TrimRight(respBuf[1:], "\x00")
	respBuf = bytes.TrimRight(respBuf, string(EOT))

	dataChan <- respBuf
}
