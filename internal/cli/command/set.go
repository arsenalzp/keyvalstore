// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/arsenalzp/keyvalstore/internal/cli/errors"
	"github.com/arsenalzp/keyvalstore/internal/cli/util"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&serverAddress, "server", "s", "", "use server and port for connection")
	setCmd.Flags().StringVarP(&client_cert, "cert", "c", "", "path to certificate file")
	setCmd.Flags().StringVarP(&privkey_cert, "key", "k", "", "path to private key file")
	setCmd.Flags().StringVarP(&rootca_cert, "CAcert", "r", "", "path to CA certificate file")
}

var setCmd = &cobra.Command{
	Use:   "set [--server] key=val",
	Short: "Set key=value",
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := CreateConnection()
		if err != nil {
			return err
		}

		if err := Set(conn, cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func Set(conn net.Conn, cmd *cobra.Command, args []string) error {
	var key, value []byte
	var buf [MESSAGE_SIZE]byte // command 3B, key 256B, value 511B

	defer conn.Close()

	// read data from arguments
	// else read data from stdin
	if len(args) > 0 {
		key, value = readArgs(args)
	} else {
		var err error

		reader := bufio.NewReader(os.Stdin)

		key, value, err = readStdin(reader)
		if err != nil {
			err := errors.New("set command error", errors.ReadStdinErr, err)
			return err
		}
	}

	// sanitize the key data
	key = sanitizeData(key)

	// sanitize the value data
	value = sanitizeData(value)

	// validate the key and the value data parameters
	err := util.ValidateInput(key, value)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(conn)

	copy(buf[0:3], []byte("set")) // copy the command data
	copy(buf[3:259], key)         // copy the key data
	copy(buf[259:], value)        // copy the value data
	buf[771] = EOT

	_, err = writer.Write(buf[:]) // write command, key and val
	if err != nil {
		err = errors.New("set command error", errors.WriteServerErr, err)
		return err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("set command error", errors.WriteServerErr, err)
		return err
	}

	reader := bufio.NewReader(conn)
	respBuf, err := reader.ReadBytes(EOT) // waiting for server response
	if err != nil {
		err = errors.New("set command error", errors.ReadServerErr, err)
		return err
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("set command error", errors.SetResponseError, err)
		return err
	}

	return nil
}
