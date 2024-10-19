// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
)

// handle DEL command
func (ds *dataStruct) del(ctx context.Context, key []byte, dataCh chan<- []byte, errCh chan<- error) {
	clearKey := bytes.Trim(key, string(EOT))
	clearKey = bytes.Trim(clearKey, "\x00")

	_, err := ds.Delete(ctx, string(clearKey))
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- []byte{}
}
