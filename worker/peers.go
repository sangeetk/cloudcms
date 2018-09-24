package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

		peers = tx.Bucket([]byte("peers"))
		if peers == nil {
			return nil
		}
		c := peers.Cursor()
		for peer, _ := c.First(); peer != nil; peer, _ = c.Next() {

			if w.String() == string(peer[:]) {
				// Ignore sync message when it was sent by itself
				continue
			}

			// Send sync request to the peer
			ctx := context.Background()
			tgt, err := url.Parse("http://" + string(peer[:]) + "/sync")
			if err != nil {
				// remove the peer
				continue
			}
			log.Println("Sending sync msg to ", tgt)
			endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
			_, err = endPoint(ctx, sreq)
			if err != nil {
				// remove the peer
				continue
			}
		}
		return nil
	})
	return err
}

// encodeRequest encodes the request as JSON
func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeResponse decodes the response from the service
func decodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response SyncResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
