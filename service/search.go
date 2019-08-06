package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	q "github.com/blevesearch/bleve/search/query"
	"github.com/go-kit/kit/endpoint"
)

// Search - searches for query
func (s *Service) Search(ctx context.Context, req *api.SearchRequest) (*api.SearchResults, error) {
	var resp = api.SearchResults{Type: req.Type, Request: req}
	var searchRequest *bleve.SearchRequest
	var query q.Query

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	if req.Query == "" {
		query = bleve.NewMatchAllQuery()
	} else if req.Fuzzy {
		query = bleve.NewFuzzyQuery(req.Query)
	} else {
		query = bleve.NewQueryStringQuery(req.Query)
	}

	// Create a new search request
	searchRequest = bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlight()
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
	resp.Took = searchResult.Took

	for _, hit := range searchResult.Hits {
		resp.Hits = append(resp.Hits, hit.Fields)
	}

	return &resp, nil
}

// SearchEndpoint - creates endpoint for Search service
func SearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.SearchRequest)
		return svc.Search(ctx, &req)
	}
}

// DecodeSearchReq - decodes the incoming request
func DecodeSearchReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
