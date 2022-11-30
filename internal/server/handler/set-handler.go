// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
)

// handle SET command
func (ds *dataStruct) set(ctx context.Context, key, val []byte, dataCh chan<- []byte, errCh chan<- error) {
	key = bytes.Trim(key, "\x00")
	val = bytes.Trim(val, "\x00")
	_, err := ds.Insert(ctx, string(key), string(val))
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- []byte{}
}
