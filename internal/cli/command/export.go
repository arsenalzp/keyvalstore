// Package implements CLI commands.

package command

import (
	"bufio"
	"bytes"
	"fmt"
	"gokeyval/internal/cli/errors"
	"gokeyval/internal/cli/util"
	"os"

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
		var buf [MESSAGE_SIZE]byte

		con, err := createConnection()
		if err != nil {
			return err
		}

		defer con.Close()

		writer := bufio.NewWriter(con)
		reader := bufio.NewReader(con)

		copy(buf[0:3], []byte("exp"))

		_, err = writer.Write(buf[:])
		if err != nil {
			err = errors.New("export command error", errors.WriteServerErr, err)
			return err
		}

		err = writer.Flush()
		if err != nil {
			err = errors.New("export command error", errors.WriteServerErr, err)
			return err
		}

		export, err := reader.ReadBytes('\x00')
		if err != nil {
			err = errors.New("export command error", errors.ReadServerErr, err)
			return err
		}

		export = bytes.TrimRight(export[1:], "\x00")
		err = util.ValidateData(export)
		if err != nil {
			err = errors.New("export command error, validation of output failed", errors.InvalidExport, err)
			return err
		}

		fmt.Fprintf(os.Stdout, "%s\n", export)

		return nil
	},
}
