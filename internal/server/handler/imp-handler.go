// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"context"
	"encoding/json"

	entity "github.com/arsenalzp/keyvalstore/internal/server/storage/entity"
)

// handle IMPORT command
func (ds *dataStruct) imp(ctx context.Context, data []byte, dataCh chan<- []byte, errCh chan<- error) {
	var importItems []entity.ImportData

	// deserialize client's JSON
	err := json.Unmarshal(data, &importItems)
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
