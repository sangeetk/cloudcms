package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Create - creates a single item
func (Service) Create(ctx context.Context, req *api.Create) (*api.Response, error) {
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
		b, err := tx.CreateBucketIfNotExists([]byte(req.Type))
		if err != nil {
			return err
		}

		nextSeq, err := b.NextSequence()
		if err != nil {
			return err
		}

		item := api.Item{
			Header: api.Header{
				ID:        nextSeq,
				Title:     req.Title,
				Slug:      stringToSlug(req.Title),
				Status:    req.Status,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
				DeletedAt: 0,
			},
			// Custom content type fields
			Content: req.Content,
		}

		slug := item.Slug
		var i int
		// Find a unique slug
		for i = 2; b.Get([]byte(slug)) != nil; i++ {
			slug = fmt.Sprintf("%s-%d", item.Slug, i)
		}
		if i > 2 {
			item.Slug = slug
		}

		itm, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = b.Put([]byte(item.Slug), itm)
		if err != nil {
			return err
		}
		resp.Item = item

		// Create index
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

// CreateEndpoint - creates endpoint for Create service
func CreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Create)
		return svc.Create(ctx, &req)
	}
}

// DecodeCreateReq - decodes the incoming request
func DecodeCreateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Create
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
