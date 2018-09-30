package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Read - returns a single item
func (s *Service) Read(ctx context.Context, req *api.ReadRequest) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}
	var db *bolt.DB

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// First read from index
	r, err := s.ReadFromIndex(ctx, req)
	if err == nil && r.Err == "" {
		return r, nil
	}

	// Open database in read-only mode
	// It will be created if it doesn't exist.
	options := bolt.Options{ReadOnly: true}
	db, err = bolt.Open(DBFile, 0644, &options)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return api.ErrorNotFound
		}
		val := b.Get([]byte(req.Slug))
		if val == nil {
			return api.ErrorNotFound
		}

		err := json.Unmarshal(val, &resp.Content)
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

// ReadFromIndex - returns a single item from index
func (s *Service) ReadFromIndex(ctx context.Context, req *api.ReadRequest) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}

	if _, ok := Index[req.Type]; !ok {
		resp.Err = "Invalid content type"
		return &resp, nil
	}
	query := bleve.NewMatchAllQuery()
	searchRequest := bleve.NewSearchRequest(query)

	searchRequest.Fields = []string{"*"}
	searchRequest.Size = 10
	searchRequest.From = 0

	for {
		searchResult, err := Index[req.Type].Search(searchRequest)
		if err != nil {
			resp.Err = api.ErrorNotFound.Error()
			return &resp, nil
		}

		for _, hit := range searchResult.Hits {
			slug := hit.Fields["slug"].(string)
			if slug == req.Slug {
				resp.Content = hit.Fields
				return &resp, nil
			}
		}
		searchRequest.From += searchRequest.Size
		if searchRequest.From >= int(searchResult.Total) {
			resp.Err = api.ErrorNotFound.Error()
			return &resp, nil
		}
	}

	resp.Err = api.ErrorNotFound.Error()
	return &resp, nil
}

// ReadEndpoint - creates endpoint for Read service
func ReadEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.ReadRequest)
		return svc.Read(ctx, &req)
	}
}

// DecodeReadReq - decodes the incoming request
func DecodeReadReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.ReadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
