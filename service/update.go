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
	"github.com/urantiatech/kit/endpoint"
)

// Update - creates a single item
func (s *Service) Update(ctx context.Context, req *api.UpdateRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type, Language: req.Language}
	var db *bolt.DB
	var err error

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Update request as sync msg contains full information
	// Simply update the index the content and return
	if sync {
		index, err := getIndex(req.Type, req.Language)
		if err != nil {
			return &resp, nil
		}

		err = index.Index(req.Slug, req.Content)
		if err != nil {
			return &resp, nil
		}
		resp.Content = req.Content
		return &resp, nil
	}

	// Normal update request
	// Open database in read-write mode
	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return api.ErrorNotFound
		}
		bb := b.Bucket([]byte(req.Language))
		if bb == nil {
			return api.ErrorNotFound
		}

		var content map[string]interface{}
		// Get the existing value
		val := bb.Get([]byte(req.Slug))
		if val == nil {
			return api.ErrorNotFound
		}
		err = json.Unmarshal(val, &content)
		if err != nil {
			return err
		}

		// Update values
		if req.Content == nil {
			return api.ErrorNullContent
		}
		var fields = (req.Content).(map[string]interface{})
		for k, v := range fields {
			content[k] = v
		}
		content["updated_at"] = time.Now().Unix()

		// Commit to database
		j, err := json.Marshal(content)
		if err != nil {
			return err
		}
		err = bb.Put([]byte(req.Slug), j)
		if err != nil {
			return err
		}

		resp.Content = content

		index, err := getIndex(req.Type, req.Language)
		if err != nil {
			return err
		}
		err = index.Index(req.Slug, content)
		if err != nil {
			return err
		}

		// Add the request to LOG bucket
		l := tx.Bucket([]byte("log"))
		seq, err := l.NextSequence()
		if err != nil {
			return err
		}
		event := api.Event{Seq: seq, Op: api.UpdateOp, Request: req}
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
		Operation: "update",
		Slug:      req.Slug,
		Timestamp: time.Now().Unix(),
		Source:    LocalWorker.String(),
		Response:  &resp,
	}
	LocalWorker.SyncPeers(SyncFile, &sreq)

	return &resp, nil
}

// UpdateEndpoint - creates endpoint for Update service
func UpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.UpdateRequest)
		return svc.Update(ctx, &req, false)
	}
}

// DecodeUpdateReq - decodes the incoming request
func DecodeUpdateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
