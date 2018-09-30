package worker

import (
	"context"
	"net/url"

	"github.com/boltdb/bolt"
	ht "github.com/urantiatech/kit/transport/http"
)

// SyncPeers syncs all other peers within the same cluster
func (w *Worker) SyncPeers(syncFile string, sreq *SyncRequest) error {
	// Open the sync database
	syncDB, err := bolt.Open(syncFile, 0644, nil)
	if err != nil {
		return err
	}
	defer syncDB.Close()

	// Add the (IP:Port, timestamp) to the database
	err = syncDB.Update(func(tx *bolt.Tx) error {
		var peers *bolt.Bucket
		var unreachable = make(map[string]bool)

		peers = tx.Bucket([]byte("peers"))
		if peers == nil {
			return nil
		}
		c := peers.Cursor()
		var peer string
		for p, _ := c.First(); p != nil; p, _ = c.Next() {
			peer = string(p[:])
			if w.String() == peer {
				// Ignore sync message when it was sent by itself
				continue
			}

			// Send sync request to the peer
			ctx := context.Background()
			tgt, err := url.Parse("http://" + peer + "/sync")
			if err != nil {
				// remove the junk url
				continue
			}
			endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
			_, err = endPoint(ctx, sreq)
			if err != nil {
				// remove the peer
				unreachable[peer] = true
				continue
			}
		}
		for peer = range unreachable {
			peers.Delete([]byte(peer))
		}
		return nil
	})
	return err
}
