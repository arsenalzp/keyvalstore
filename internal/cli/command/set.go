// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"gokeyval/internal/cli/errors"
	"gokeyval/internal/cli/util"
	"net"
	"os"
	"strings"

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
		if err := Set(nil, cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func Set(externalConn net.Conn, cmd *cobra.Command, args []string) error {
	var substr []string
	var buf [MESSAGE_SIZE]byte // command 3B, key 256B, value 511B
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

	// read data from arguments
	// else read data from stdin
	if len(args) > 0 {
		substr = strings.SplitN(args[0], "=", 2) // split key and value from args
	} else {
		r := bufio.NewReader(os.Stdin)
		input, err := r.ReadString(byte('\n')) // read data from stdin
		if err != nil {
			err = errors.New("set command error", errors.ReadStdinErr, err)
			return err
		}
		substr = strings.Split(input, "=") // split key and value from stdin
	}

	if len(substr) == 1 {
		return fmt.Errorf("key and val shouldn't be empty")
	}

	// validate the key and the value data parameters
	err := util.ValidateInput(args[0], substr[1])
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(con)

	copy(buf[0:3], []byte("set"))       // copy the command data
	copy(buf[3:259], []byte(substr[0])) // copy the key data
	copy(buf[259:], []byte(substr[1]))  // copy the value data
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

	reader := bufio.NewReader(con)
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
