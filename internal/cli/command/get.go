// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"

	"github.com/arsenalzp/keyvalstore/internal/cli/errors"
	"github.com/arsenalzp/keyvalstore/internal/cli/util"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&serverAddress, "server", "s", "", "use server and port for connection")
	getCmd.Flags().StringVarP(&client_cert, "cert", "c", "", "path to certificate file")
	getCmd.Flags().StringVarP(&privkey_cert, "key", "k", "", "path to private key file")
	getCmd.Flags().StringVarP(&rootca_cert, "CAcert", "r", "", "path to CA certificate file")
}

var getCmd = &cobra.Command{
	Use:   "get [--server] key",
	Short: "Get value of a key",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := CreateConnection()
		if err != nil {
			return err
		}

		data, err := Get(conn, cmd, args)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%s\n", data)
		return nil
	},
}

func Get(conn net.Conn, cmd *cobra.Command, args []string) ([]byte, error) {
	var key []byte
	var buf [MESSAGE_SIZE]byte

	defer conn.Close()

	// read data from arguments
	// else read data from stdin
	if len(args) > 0 {
		key, _ = readArgs(args)
	} else {
		var err error

		reader := bufio.NewReader(os.Stdin)

		key, _, err = readStdin(reader)
		if err != nil {
			err := errors.New("del command error", errors.ReadStdinErr, err)
			return nil, err
		}
	}

	// sanitize the key data
	key = sanitizeData(key)

	// validate the key data parameter
	err := util.ValidateInput(key, []byte{})
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(conn)

	copy(buf[0:3], []byte("get")) // copy the command data
	copy(buf[3:259], key)         // copy the key data
	buf[771] = EOT

	_, err = writer.Write(buf[:]) // write command, key and val
	if err != nil {
		err = errors.New("get command error", errors.WriteServerErr, err)
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("get command error", errors.WriteServerErr, err)
		return nil, err
	}

	reader := bufio.NewReader(conn)
	respBuf, err := reader.ReadBytes(EOT) // reading the server response
	if err != nil {
		err = errors.New("get command failed", errors.ReadServerErr, err)
		return nil, err
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("get command failed", errors.GetResponseError, err)
		return nil, err
	}

	// Trim response buffer: delete NULL and EOT bytes
	respBuf = bytes.TrimRight(respBuf, string(EOT))
	respBuf = bytes.TrimRight(respBuf[1:], "\x00")

	return respBuf, err
}
