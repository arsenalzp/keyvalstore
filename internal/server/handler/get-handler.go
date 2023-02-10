// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
)

// handle GET commahd
func (ds *dataStruct) get(ctx context.Context, key []byte, dataCh chan<- []byte, errCh chan<- error) {
	clearKey := bytes.Trim(key, string(EOT))
	clearKey = bytes.Trim(clearKey, "\x00")

	val, err := ds.Search(ctx, string(clearKey))
	if err != nil {
		errCh <- err
		return
	}

	data := []byte(val)
	dataCh <- data
}
