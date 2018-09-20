package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Update - creates a single item
func (s *Service) Update(ctx context.Context, req *api.UpdateRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}
	var db *bolt.DB
	var err error

	if _, ok := Index[req.Type]; !ok {
		resp.Err = ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Update request as sync msg contains full information
	// Simply update the index the content and return
	if sync {
		IndexLock.Lock()
		defer IndexLock.Unlock()

		err = Index[req.Type].Index(req.Slug, req.Content)
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
			return ErrorNotFound
		}

		var content map[string]interface{}
		// Get the existing value
		val := b.Get([]byte(req.Slug))
		if val == nil {
			return ErrorNotFound
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

		// Update index
		IndexLock.Lock()
		defer IndexLock.Unlock()

		err = Index[req.Type].Index(req.Slug, content)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
	}

	// Sync others

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
