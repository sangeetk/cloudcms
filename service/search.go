package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/urantiatech/kit/endpoint"
)

// Search - creates a single item
func (Service) Search(ctx context.Context, req *api.Search) (*api.SearchResults, error) {
	var resp api.SearchResults
	var searchRequest *bleve.SearchRequest

	if _, ok := Index[req.Type]; !ok {
		resp.Err = "Invalid content type"
		return &resp, nil
	}

	// Add date-range if applicable
	if !req.StartDate.IsZero() || !req.EndDate.IsZero() {
		if req.EndDate.IsZero() {
			req.EndDate = time.Now()
		}
		dateRangeQuery := bleve.NewDateRangeQuery(req.StartDate, req.EndDate)
		stringQuery := bleve.NewQueryStringQuery(req.Query)
		conjunctionQuery := bleve.NewConjunctionQuery(dateRangeQuery, stringQuery)
		searchRequest = bleve.NewSearchRequest(conjunctionQuery)
	} else {
		stringQuery := bleve.NewQueryStringQuery(req.Query)
		searchRequest = bleve.NewSearchRequest(stringQuery)
	}

	searchRequest.Fields = []string{"*"}
	searchRequest.Highlight = bleve.NewHighlight()
	searchRequest.Size = req.Limit
	if searchRequest.Size <= 0 {
		searchRequest.Size = 10
	}
	searchRequest.From = req.Skip

	searchResult, err := Index[req.Type].Search(searchRequest)
	if err != nil {
		resp.Err = ErrorNotFound.Error()
		return &resp, nil
	}

	resp.Total = int(searchResult.Total)
	resp.Limit = len(searchResult.Hits)
	resp.Skip = req.Skip

	for _, hit := range searchResult.Hits {
		item := api.Item{
			Header: api.Header{
				ID:        uint64(hit.Fields["id"].(float64)),
				Title:     hit.Fields["title"].(string),
				Slug:      hit.Fields["slug"].(string),
				Status:    hit.Fields["status"].(string),
				CreatedAt: int64(hit.Fields["created_at"].(float64)),
				UpdatedAt: int64(hit.Fields["updated_at"].(float64)),
			},
		}
		resp.Results = append(resp.Results, item)
	}

	return &resp, nil
}

// SearchEndpoint - creates endpoint for Search service
func SearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Search)
		return svc.Search(ctx, &req)
	}
}

// DecodeSearchReq - decodes the incoming request
func DecodeSearchReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Search
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
