package service

import (
	"context"
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
	var resp = api.Response{Type: req.Type}
	var db *bolt.DB
	var err error
	log.Println("Update()", "Type:", req.Type, "Slug:", req.Slug, "Sync:", sync)

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Forward request to Upstream Server
	if !sync && Upstream.Host != "" {
		log.Println("Forwarding to ", Upstream)
		return LocalWorker.Forward("update", req, Upstream)
	}

	// Update request as sync msg contains full information
	// Simply update the index the content and return
	log.Println("Update sync: ", req.Slug, sync)
	if sync {
		log.Println("Received sync message by ", LocalWorker)

		err = Index[req.Type].Index(req.Slug, req.Content)
		if err != nil {
			return &resp, nil
		}
		resp.Content = req.Content
		return &resp, nil
	}

	// Normal update request
	// Open database in read-write mode
	log.Println("Normal update request")
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

		var content map[string]interface{}
		// Get the existing value
		val := b.Get([]byte(req.Slug))
		if val == nil {
			return api.ErrorNotFound
		}
		err := json.Unmarshal(val, &content)
		if err != nil {
			return err
		}

		// Update values
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
		err = b.Put([]byte(req.Slug), j)
		if err != nil {
			return err
		}

		resp.Content = content

		err = Index[req.Type].Index(req.Slug, content)
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
	LocalWorker.SyncChilds(SyncFile, &sreq)

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
