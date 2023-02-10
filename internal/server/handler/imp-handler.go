// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	entity "gokeyval/internal/server/storage/entity"
)

// handle IMPORT command
func (ds *dataStruct) imp(ctx context.Context, data []byte, dataCh chan<- []byte, errCh chan<- error) {
	var importItems []entity.ImportData

	clearImport := bytes.Trim(data, string(EOT))
	clearImport = bytes.Trim(clearImport, "\x00")

	// deserialize client's JSON
	err := json.Unmarshal(clearImport, &importItems)
	if err != nil {
		errCh <- err
		return
	}

	_, err = ds.Import(ctx, importItems)
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- []byte{}
}
