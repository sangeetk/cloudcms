package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/urantiatech/kit/endpoint"
)

// List - list all items
func (s *Service) List(ctx context.Context, req *api.ListRequest) (*api.ListResults, error) {
	var resp = api.ListResults{Type: req.Type}
	var searchRequest *bleve.SearchRequest

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	stringQuery := bleve.NewMatchAllQuery()
	searchRequest = bleve.NewSearchRequest(stringQuery)
	searchRequest.SortBy([]string{"id"})

	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Size = req.Limit
	if searchRequest.Size <= 0 {
		searchRequest.Size = 10
	}
	searchRequest.From = req.Skip

	searchResult, err := Index[req.Type].Search(searchRequest)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}

	resp.Total = int(searchResult.Total)
	resp.Limit = len(searchResult.Hits)
	resp.Skip = req.Skip

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
