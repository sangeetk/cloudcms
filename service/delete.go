package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Delete - creates a single item
func (Service) Delete(ctx context.Context, req *api.Delete) (*api.Response, error) {
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

		err = b.Delete([]byte(resp.Item.Slug))
		if err != nil {
			return err
		}

		// Delete index
		err = Index[req.Type].Delete(resp.Item.Slug)
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

// DeleteEndpoint - creates endpoint for Delete service
func DeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Delete)
		return svc.Delete(ctx, &req)
	}
}

// DecodeDeleteReq - decodes the incoming request
func DecodeDeleteReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Delete
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
