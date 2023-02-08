// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"
	"io"
	"net"
	"os"
)

func Get(con *net.Conn, c chan<- byte[], k string) {
	var respBuf [256]byte
	var buf [MESSAGE_SIZE]byte

	writer := bufio.NewWriter(con)

	copy(buf[0:3], []byte("get"))
	copy(buf[3:], []byte(k))

	_, err = writer.Write(buf[:]) // write command, key and val
	if err != nil {
		err = errors.New("get command error", errors.WriteServerErr, err)
		return err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("get command error", errors.WriteServerErr, err)
		return err
	}

	_, err = io.ReadFull(con, respBuf[:]) // waiting for server response
	if err != nil {
		err = errors.New("get command failed", errors.WriteServerErr, err)
		fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
		return err
	}

	respCode := respBuf[:1]
	if respCode[0] != 'O' {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("get command failed", errors.WriteServerErr, err)
		fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
		return err
	}

	fmt.Fprintf(os.Stdout, "%s\n", bytes.TrimRight(respBuf[1:], "\x00"))

	return nil
}
