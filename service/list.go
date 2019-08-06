package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/go-kit/kit/endpoint"
)

// List - list all items
func (s *Service) List(ctx context.Context, req *api.ListRequest) (*api.ListResults, error) {
	var resp = api.ListResults{Type: req.Type, Request: req}
	var searchRequest *bleve.SearchRequest

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	query := bleve.NewMatchAllQuery()
	searchRequest = bleve.NewSearchRequest(query)

	if req.SortBy == "" {
		req.SortBy = "id"
	}
	searchRequest.SortBy([]string{req.SortBy})
	searchRequest.Fields = []string{"*"}
	searchRequest.Size = req.Size
	if searchRequest.Size <= 0 {
		searchRequest.Size = 10
	}
	searchRequest.From = req.Skip

	index, err := getIndex(req.Type, req.Language)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}

	resp.Total = searchResult.Total

	for _, hit := range searchResult.Hits {
		resp.List = append(resp.List, hit.Fields)
	}

	return &resp, nil
}

// ListEndpoint - creates endpoint for List service
func ListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.ListRequest)
		return svc.List(ctx, &req)
	}
}

// DecodeListReq - decodes the incoming request
func DecodeListReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
