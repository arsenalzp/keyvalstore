// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"gokeyval/internal/cli/errors"
	"gokeyval/internal/cli/util"
	"net"
	"os"

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
		_, err := Get(nil, cmd, args)
		if err != nil {
			return err
		}
		return nil
	},
}

func Get(externalConn net.Conn, cmd *cobra.Command, args []string) ([]byte, error) {
	var buf [MESSAGE_SIZE]byte
	var con net.Conn

	if externalConn == nil {
		newCon, err := CreateConnection()
		if err != nil {
			return nil, err
		}
		con = newCon
	} else {
		con = externalConn
	}

	defer con.Close()

	// validate the key data parameter
	err := util.ValidateInput(args[0], "")
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(con)

	copy(buf[0:3], []byte("get"))     // copy the command data
	copy(buf[3:259], []byte(args[0])) // copy the key data
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

	reader := bufio.NewReader(con)
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
	fmt.Fprintf(os.Stdout, "%s\n", respBuf)

	return respBuf, err
}
