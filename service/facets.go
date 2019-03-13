package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
	q "github.com/blevesearch/bleve/search/query"
	"github.com/urantiatech/kit/endpoint"
)

// FacetsSearch - searches for query with multiple facets
func (s *Service) FacetsSearch(ctx context.Context, req *api.FacetsSearchRequest) (*api.FacetsSearchResults, error) {

	if j, err := json.Marshal(req); err == nil {
		fmt.Println(string(j))
	}

	log.Printf("req.Query=[%s]\n", req.Query)
	var resp = api.FacetsSearchResults{Type: req.Type}
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
	if req.Query == "" {
		sf := &search.SortField{Field: "created_at", Desc: true}
		searchRequest.SortByCustom(search.SortOrder{sf})
	}

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
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}

	// Create response structure
	// Type    string               `json:"type"`
	resp.Type = req.Type

	// Status  *SearchStatus        `json:"status"`
	resp.Status = &api.SearchStatus{
		Total:      searchResults.Status.Total,
		Failed:     searchResults.Status.Failed,
		Successful: searchResults.Status.Successful,
	}

	// Request *FacetsSearchRequest `json:"request"`
	resp.Request = req

	// Hits    []interface{}        `json:"hits"`
	for _, hit := range searchResults.Hits {
		resp.Hits = append(resp.Hits, hit.Fields)
	}

	// Total   uint64               `json:"total"`
	resp.Total = searchResults.Total

	// Took    time.Duration        `json:"took"`
	resp.Took = searchResults.Took

	// Facets  FacetResults         `json:"facets"`
	resp.Facets = make(api.FacetResults)
	for fname, fresult := range searchResults.Facets {
		facetResult := api.FacetResult{
			Field:   fresult.Field,
			Total:   fresult.Total,
			Missing: fresult.Missing,
			Other:   fresult.Other,
		}

		for _, tfacet := range fresult.Terms {
			tf := &api.TermFacet{Term: tfacet.Term, Count: tfacet.Count}
			facetResult.Terms = append(facetResult.Terms, tf)
		}

		for n, nfacet := range fresult.NumericRanges {
			facetResult.NumericRanges[n].Name = nfacet.Name
			facetResult.NumericRanges[n].Min = nfacet.Min
			facetResult.NumericRanges[n].Max = nfacet.Max
			facetResult.NumericRanges[n].Count = nfacet.Count
		}

		for d, dfacet := range fresult.DateRanges {
			facetResult.DateRanges[d].Name = dfacet.Name
			facetResult.DateRanges[d].Start = dfacet.Start
			facetResult.DateRanges[d].End = dfacet.End
			facetResult.DateRanges[d].Count = dfacet.Count
		}

		resp.Facets[fname] = &facetResult
	}

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
