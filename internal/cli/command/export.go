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
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&serverAddress, "server", "s", "", "use server and port for connection")
	exportCmd.Flags().StringVarP(&client_cert, "cert", "c", "", "path to certificate file")
	exportCmd.Flags().StringVarP(&privkey_cert, "key", "k", "", "path to private key file")
	exportCmd.Flags().StringVarP(&rootca_cert, "CAcert", "r", "", "path to CA certificate file")
}

var exportCmd = &cobra.Command{
	Use:   "export [--server]",
	Short: "Retrieve key=value pairs and print them into stdout ",
	Args:  cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := CreateConnection()
		if err != nil {
			return err
		}

		data, err := Export(conn, cmd)
		if err != nil {
			return err
		}

		fmt.Fprintf(os.Stdout, "%s\n", data)
		return nil
	},
}

func Export(conn net.Conn, cmd *cobra.Command) ([]byte, error) {
	var buf []byte = make([]byte, 4)

	defer conn.Close()

	writer := bufio.NewWriter(conn)

	copy(buf[0:3], []byte("exp")) // copy the command data
	buf[3] = EOT                  // add EOT to signal the end of transmission

	_, err := writer.Write(buf)
	if err != nil {
		err = errors.New("export command error", errors.WriteServerErr, err)
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		err = errors.New("export command error", errors.WriteServerErr, err)
		return nil, err
	}

	reader := bufio.NewReader(conn)
	respBuf, err := reader.ReadBytes(EOT)
	if err != nil {
		err = errors.New("export command error", errors.ReadServerErr, err)
		return nil, err
	}

	if respBuf[0] == errors.ServerResponseError {
		err = fmt.Errorf("%s", respBuf[1:])
		err = errors.New("export command error", errors.ExpResponseError, err)
		return nil, err
	}

	// Trim response buffer: delete NULL and EOT bytes
	respBuf = bytes.TrimRight(respBuf[1:], "\x00")
	respBuf = bytes.TrimRight(respBuf, string(EOT))

	err = util.ValidateData(respBuf)
	if err != nil {
		err = errors.New("export command error, validation of output failed", errors.InvalidExport, err)
		return nil, err
	}

	return respBuf, err
}
