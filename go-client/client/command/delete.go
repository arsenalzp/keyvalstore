// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/arsenalzp/keyvalstore/go-client/internal/errors"

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
		var respBuf [64]byte
		var buf [MESSAGE_SIZE]byte

		con, err := createConnection()
		if err != nil {
			return err
		}

		defer con.Close()

		key := args[0]

		writer := bufio.NewWriter(con)

		copy(buf[0:3], []byte("del"))
		copy(buf[3:], []byte(key))

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

		_, err = io.ReadFull(con, respBuf[:]) // waiting for server response
		if err != nil {
			err = errors.New("del command failed", errors.WriteServerErr, err)
			fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
			return err
		}

		respCode := respBuf[:1]
		if respCode[0] != 'O' {
			err = fmt.Errorf("%s", respBuf[1:])
			err = errors.New("del command failed", errors.WriteServerErr, err)
			fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
			return err
		}

		return nil
	},
}
