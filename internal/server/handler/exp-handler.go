// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"context"
	"encoding/json"
)

// export EXPORT command
func (ds *dataStruct) exp(ctx context.Context, dataCh chan<- []byte, errCh chan<- error) {
	exports, err := ds.Export(ctx)
	if err != nil {
		errCh <- err
		return
	}

	data, err := json.Marshal(exports)
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- data
}
