package service

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/boltdb/bolt"
	"github.com/go-kit/kit/endpoint"
)

// Delete - creates a single item
func (s *Service) Delete(ctx context.Context, req *api.DeleteRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type, Language: req.Language}
	var db *bolt.DB
	var err error

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Update request as sync msg contains full information
	// Simply index the content and return
	if sync {
		readReq := api.ReadRequest{Type: req.Type, Language: req.Language, Slug: req.Slug}
		item, err := s.Read(ctx, &readReq)
		if err != nil {
			resp.Err = api.ErrorNotFound.Error()
			return &resp, nil
		}
		index, err := getIndex(req.Type, req.Language)
		if err != nil {
			return &resp, nil
		}
		err = index.Delete(req.Slug)
		if err != nil {
			return &resp, nil
		}
		resp.Content = item.Content
		return &resp, nil
	}

	// Open database in read-write mode
	// It will be created if it doesn't exist.
	//options := bolt.Options{ReadOnly: false}

	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bb, err := getBucket(tx, req.Type, req.Language)
		if err != nil {
			return err
		}

		// Get the existing value
		val := bb.Get([]byte(req.Slug))
		if val == nil {
			return api.ErrorNotFound
		}

		err = json.Unmarshal(val, &resp.Content)
		if err != nil {
			return err
		}

		err = bb.Delete([]byte(req.Slug))
		if err != nil {
			return err
		}

		index, err := getIndex(req.Type, req.Language)
		if err != nil {
			return err
		}
		err = index.Delete(req.Slug)
		if err != nil {
			return err
		}

		// Add the request to LOG bucket
		l := tx.Bucket([]byte("log"))
		seq, err := l.NextSequence()
		if err != nil {
			return err
		}
		event := api.Event{Seq: seq, Op: api.DeleteOp, Request: req}
		e, err := json.Marshal(event)
		if err != nil {
			return err
		}

		var key = make([]byte, 8)
		binary.LittleEndian.PutUint64(key, seq)
		err = l.Put(key, e)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		resp.Err = err.Error()
		return &resp, nil
	}

	// Sync other workers
	sreq := worker.SyncRequest{
		Type:      req.Type,
		Operation: "delete",
		Slug:      req.Slug,
		Timestamp: time.Now().Unix(),
		Source:    LocalWorker.String(),
		Response:  &resp,
	}
	LocalWorker.SyncPeers(SyncFile, &sreq)

	return &resp, nil
}

// DeleteEndpoint - creates endpoint for Delete service
func DeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.DeleteRequest)
		return svc.Delete(ctx, &req, false)
	}
}

// DecodeDeleteReq - decodes the incoming request
func DecodeDeleteReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
