// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"gokeyval/internal/cli/errors"
	"io"
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
		var respBuf [256]byte
		var buf [MESSAGE_SIZE]byte

		con, err := createConnection()
		if err != nil {
			return err
		}

		defer con.Close()

		key := args[0]

		writer := bufio.NewWriter(con)

		copy(buf[0:3], []byte("get"))
		copy(buf[3:], []byte(key))

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
	},
}
