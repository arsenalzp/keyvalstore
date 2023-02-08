// Package implements CLI commands.

package command

import (
	"bufio"
	"fmt"
	"gokeyval/internal/cli/errors"
	"gokeyval/internal/cli/util"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringVarP(&serverAddress, "server", "s", "", "use server and port for connection")
	importCmd.Flags().StringVarP(&client_cert, "cert", "c", "", "path to certificate file")
	importCmd.Flags().StringVarP(&privkey_cert, "key", "k", "", "path to private key file")
	importCmd.Flags().StringVarP(&rootca_cert, "CAcert", "r", "", "path to CA certificate file")
}

var importCmd = &cobra.Command{
	Use:   "import [--server]",
	Short: "Import stringified key=value pairs from stdin ",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var buf []byte = make([]byte, 3)
		var respBuf [64]byte

		con, err := createConnection()
		if err != nil {
			return err
		}

		defer con.Close()

		writer := bufio.NewWriter(con) // create connection writer
		copy(buf[0:3], []byte("imp"))

		// read data from args
		// else read data from stdin
		if len(args) > 0 {
			data := args[0]
			err := util.ValidateData([]byte(data)) // validate data befor sendind
			if err != nil {
				err = errors.New("import command error", errors.ReadStdinErr, err)
				return err
			}

			buf = append(buf[:], data...)
			buf = append(buf[:], '\x00') // add delimiter to the end of the buffer

			// if a length of the buffer less than MESSAGE_SIZE - fill it up to 711 byte
			if len(buf) < MESSAGE_SIZE {
				diffBuf := make([]byte, MESSAGE_SIZE-len(buf))
				buf = append(buf, diffBuf...)
			}

			_, err = writer.Write(buf[:]) // send data to the server
			if err != nil {
				err = errors.New("import command error", errors.WriteServerErr, err)
				return err
			}
		} else {
			reader := bufio.NewReader(os.Stdin) // read data from stdin
			data, err := reader.ReadString(byte('\n'))
			if err != nil {
				err = errors.New("import command error", errors.ReadStdinErr, err)
				return err // error reading data from stdin
			}

			err = util.ValidateData([]byte(data)) // validate data befor sending
			if err != nil {
				err = errors.New("import command error, validation of input failed", errors.ReadStdinErr, err)
				return err
			}

			buf = append(buf[:], data...)
			buf = append(buf[:], '\x00') // add delimiter to the end of the buffer

			// if a length of the buffer less than 771 byte - fill it up to 711 byte
			if len(buf) < MESSAGE_SIZE {
				diffBuf := make([]byte, MESSAGE_SIZE-len(buf))
				buf = append(buf, diffBuf...)
			}

			_, err = writer.Write(buf[:]) // send data to the server
			if err != nil {
				err = errors.New("import command error", errors.WriteServerErr, err)
				return err
			}

		}

		err = writer.Flush()
		if err != nil {
			err = errors.New("import command error", errors.WriteServerErr, err)
			return err
		}

		_, err = io.ReadFull(con, respBuf[:]) // waiting for server response
		if err != nil {
			err = errors.New("import command error", errors.WriteServerErr, err)
			fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
			return err
		}

		respCode := respBuf[:1]
		if respCode[0] != 'O' {
			err = fmt.Errorf("%s", respBuf[1:])
			err = errors.New("import command error", errors.WriteServerErr, err)
			fmt.Fprintf(os.Stderr, "%s\n", err) // print server response
			return err
		}

		return nil
	},
}
