// Handle incoming connection by reading command from a connection
// then run related handler.

package handler

import (
	"bytes"
	"context"
)

// handle DEL command
func (ds *dataStruct) del(ctx context.Context, key *[]byte, dataCh chan<- *[]byte, errCh chan<- error) {
	*key = bytes.Trim(*key, "\x00")

	_, err := ds.Delete(ctx, string(*key))
	if err != nil {
		errCh <- err
		return
	}

	dataCh <- &[]byte{}
}
