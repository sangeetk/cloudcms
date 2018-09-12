package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Read - creates a single item
func (Service) Read(ctx context.Context, req *api.Read) (*api.Response, error) {
	var resp api.Response
	var err error

	if _, ok := Index[req.Type]; !ok {
		resp.Err = "Invalid content type"
		return &resp, nil
	}

	// Open database in read-only mode
	// It will be created if it doesn't exist.
	options := bolt.Options{ReadOnly: true}
	DB, err = bolt.Open(dbFile, 0644, &options)
	if err != nil {
		resp.Err = ErrorNotFound.Error()
		return &resp, nil
	}
	defer DB.Close()

	err = DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return ErrorNotFound
		}

		val := b.Get([]byte(req.Slug))
		if val == nil {
			return ErrorNotFound
		}
		err := json.Unmarshal(val, &resp.Item)
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

// ReadEndpoint - creates endpoint for Read service
func ReadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Read)
		return svc.Read(ctx, &req)
	}
}

// DecodeReadReq - decodes the incoming request
func DecodeReadReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Read
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
