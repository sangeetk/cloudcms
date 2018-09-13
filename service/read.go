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
func (s *Service) Read(ctx context.Context, req *api.Read) (*api.Response, error) {
	var resp api.Response

	if _, ok := Index[req.Type]; !ok {
		resp.Err = "Invalid content type"
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

// ReadFromIndex - returns a single item from index
func (s *Service) ReadFromIndex(ctx context.Context, req *api.Read) (*api.Response, error) {
	var resp api.Response

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
			resp.Err = ErrorNotFound.Error()
			return &resp, nil
		}

		for _, hit := range searchResult.Hits {
			slug := hit.Fields["slug"].(string)
			if slug == req.Slug {
				resp.Item = api.Item{
					Header: api.Header{
						ID:        uint64(hit.Fields["id"].(float64)),
						Title:     hit.Fields["title"].(string),
						Slug:      slug,
						Status:    hit.Fields["status"].(string),
						CreatedAt: int64(hit.Fields["created_at"].(float64)),
						UpdatedAt: int64(hit.Fields["updated_at"].(float64)),
					},
				}
				return &resp, nil
			}
		}
		searchRequest.From += searchRequest.Size
		if searchRequest.From >= int(searchResult.Total) {
			resp.Err = ErrorNotFound.Error()
			return &resp, nil
		}
	}

	resp.Err = ErrorNotFound.Error()
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
