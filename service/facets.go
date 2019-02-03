package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	q "github.com/blevesearch/bleve/search/query"
	"github.com/urantiatech/kit/endpoint"
)

// FacetsSearch - searches for query with multiple facets
func (s *Service) FacetsSearch(ctx context.Context, req *api.FacetsSearchRequest) (*api.FacetsSearchResults, error) {
	var resp = api.FacetsSearchResults{Type: req.Type}
	var searchRequest *bleve.SearchRequest
	var query q.Query

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	if req.Query == "" {
		query = bleve.NewMatchAllQuery()
	} else {
		query = bleve.NewQueryStringQuery(req.Query)
	}

	// Create a new search request
	searchRequest = bleve.NewSearchRequest(query)

	// Add each facet request to search
	for fname, f := range req.Facets {
		facet := bleve.NewFacetRequest(f.Field, f.Size)
		for tname, trange := range f.DateTimeRanges {
			if trange.Start.IsZero() && trange.End.IsZero() {
				trange.End = time.Now()
			}
			facet.AddDateTimeRange(tname, trange.Start, trange.End)
		}
		for nname, nrange := range f.NumericRanges {
			facet.AddNumericRange(nname, &nrange.Min, &nrange.Max)
		}
		searchRequest.AddFacet(fname, facet)
	}

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

	// for _, hit := range searchResult.Hits {
	//	resp.Hits = append(resp.Hits, hit.Fields)
	// }

	return &resp, nil
}

// FacetsSearchEndpoint - creates endpoint for Search service
func FacetsSearchEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.FacetsSearchRequest)
		return svc.FacetsSearch(ctx, &req)
	}
}

// DecodeFacetsSearchReq - decodes the incoming request
func DecodeFacetsSearchReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.FacetsSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
