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
func (Service) Update(ctx context.Context, req *api.Update) (*api.Response, error) {
	var resp api.Response
	var err error

	if _, ok := Index[req.Type]; !ok {
		resp.Err = "Invalid content type"
		return &resp, nil
	}

	// Open database in read-write mode
	// It will be created if it doesn't exist.
	//options := bolt.Options{ReadOnly: false}
	DB, err = bolt.Open(dbFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	err = DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return ErrorNotFound
		}

		// Get the existing value
		val := b.Get([]byte(req.Slug))
		if val == nil {
			return ErrorNotFound
		}
		err := json.Unmarshal(val, &resp.Item)
		if err != nil {
			return err
		}

		// Update values
		resp.Item.Title = req.Title
		resp.Item.Status = req.Status
		resp.Item.Content = req.Content
		resp.Item.UpdatedAt = time.Now().Unix()

		// Commit to database
		itm, err := json.Marshal(resp.Item)
		if err != nil {
			return err
		}
		err = b.Put([]byte(resp.Item.Slug), itm)
		if err != nil {
			return err
		}

		// Update index
		err = Index[req.Type].Index(resp.Item.Slug, resp.Item)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
	}

	return &resp, nil
}

// UpdateEndpoint - creates endpoint for Update service
func UpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Update)
		return svc.Update(ctx, &req)
	}
}

// DecodeUpdateReq - decodes the incoming request
func DecodeUpdateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Update
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
