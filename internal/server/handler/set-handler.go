// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
)

// handle SET command
func (ds *dataStruct) set(ctx context.Context, key, val []byte, dataCh chan<- []byte, errCh chan<- error) {
	clearKey := bytes.Trim(key, string(EOT))
	clearKey = bytes.Trim(clearKey, "\x00")

	clearValue := bytes.Trim(val, string(EOT))
	clearValue = bytes.Trim(clearValue, "\x00")

	_, err := ds.Insert(ctx, string(clearKey), string(clearValue))
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- []byte{}
}
