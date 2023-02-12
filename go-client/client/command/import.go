// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
)

// Import key=value pairs into a server
func Import(con net.Conn, dataChan chan<- struct{}, errChan chan<- error, dataBuffer []byte) {
	var buf []byte = make([]byte, 3)

	writer := bufio.NewWriter(con) // connection writer to send the data to the server

	copy(buf[0:3], []byte("imp"))

	buf = append(buf, dataBuffer...)
	buf = append(buf, EOT) // add EOT to signal the end of transmission

	_, err := writer.Write(buf) // send data to the server
	if err != nil {
		err = errors.New("import operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("import operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT)
	if err != nil {
		err = errors.New("import operation error", errors.WriteServerErr, err)
		errChan <- err
		return
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("import operation error", errors.ImpServerRespErr, err)
		errChan <- err
		return
	}

	dataChan <- struct{}{}
}
