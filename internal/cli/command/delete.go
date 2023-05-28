// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"gokeyval/internal/cli/errors"
	"gokeyval/internal/cli/util"
	"net"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(delCmd)
	delCmd.Flags().StringVarP(&serverAddress, "server", "s", "", "use server and port for connection")
	delCmd.Flags().StringVarP(&client_cert, "cert", "c", "", "path to certificate file")
	delCmd.Flags().StringVarP(&privkey_cert, "key", "k", "", "path to private key file")
	delCmd.Flags().StringVarP(&rootca_cert, "CAcert", "r", "", "path to CA certificate file")
}

var delCmd = &cobra.Command{
	Use:   "del [--server] key",
	Short: "Delete a key",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := Del(nil, cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func Del(externalConn net.Conn, cmd *cobra.Command, args []string) error {
	var buf [MESSAGE_SIZE]byte
	var con net.Conn

	if externalConn == nil {
		newCon, err := CreateConnection()
		if err != nil {
			return err
		}
		con = newCon
	} else {
		con = externalConn
	}

	defer con.Close()

	// validate the key data parameter
	err := util.ValidateInput(args[0], "")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(con)

	copy(buf[0:3], []byte("del"))  // copy the command data
	copy(buf[3:], []byte(args[0])) // copy the key data
	buf[771] = EOT

	_, err = writer.Write(buf[:])
	if err != nil {
		err = errors.New("del command error", errors.WriteServerErr, err)
		return err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("del command error", errors.WriteServerErr, err)
		return err
	}

	reader := bufio.NewReader(con)
	respBuf, err := reader.ReadBytes(EOT) // waiting for server response
	if err != nil {
		err = errors.New("del command failed", errors.ReadServerErr, err)
		return err
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("del command failed", errors.DelResponseError, err)
		return err
	}

	return nil
}
