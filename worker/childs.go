package worker

import (
	"context"
	"log"
	"net/url"

	"github.com/boltdb/bolt"
	ht "github.com/urantiatech/kit/transport/http"
)

// SyncChilds syncs all other childs within the same cluster
func (w *Worker) SyncChilds(syncFile string, sreq *SyncRequest) error {
	// Open the sync database
	syncDB, err := bolt.Open(syncFile, 0644, nil)
	if err != nil {
		return err
	}
	defer syncDB.Close()

	// Add the (IP:Port, timestamp) to the database
	err = syncDB.Update(func(tx *bolt.Tx) error {
		var childs *bolt.Bucket

		childs = tx.Bucket([]byte("childs"))
		if childs == nil {
			return nil
		}
		c := childs.Cursor()
		for child, _ := c.First(); child != nil; child, _ = c.Next() {

			if w.String() == string(child[:]) {
				// Ignore sync message when it was sent by itself
				continue
			}

			// Send sync request to the child
			ctx := context.Background()
			tgt, err := url.Parse("http://" + string(child[:]) + "/sync")
			if err != nil {
				// remove the child
				continue
			}
			log.Println("Sending sync msg to ", tgt)
			endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
			_, err = endPoint(ctx, sreq)
			if err != nil {
				// remove the child
				continue
			}
		}
		return nil
	})
	return err
}
