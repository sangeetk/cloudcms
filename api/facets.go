package api

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	ht "github.com/go-kit/kit/transport/http"
)

// NumericRange
type NumericRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// DateTimeRange
type DateTimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// TermFacet
type TermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

type TermFacets []*TermFacet

// NumericRangeFacet
type NumericRangeFacet struct {
	Name  string   `json:"name"`
	Min   *float64 `json:"min,omitempty"`
	Max   *float64 `json:"max,omitempty"`
	Count int      `json:"count"`
}

type NumericRangeFacets []*NumericRangeFacet

// DateTimeRangeFacet
type DateRangeFacet struct {
	Name  string  `json:"name"`
	Start *string `json:"start,omitempty"`
	End   *string `json:"end,omitempty"`
	Count int     `json:"count"`
}

type DateRangeFacets []*DateRangeFacet

//  FacetRequest describes a facet or aggregation of the result document set you would like to be built.
type FacetRequest struct {
	Size           int                       `json:"size"`
	Field          string                    `json:"field"`
	NumericRanges  map[string]*NumericRange  `json:"numeric_ranges,omitempty"`
	DateTimeRanges map[string]*DateTimeRange `json:"datetime_ranges,omitempty"`
}

type FacetsRequest map[string]*FacetRequest

//
// FacetSearchRequest - facet search request with multiple facets
//
type FacetsSearchRequest struct {
	Type     string        `json:"type"`
	Language string        `json:"language"`
	Query    string        `json:"query"`
	Fuzzy    bool          `json:"fuzzy"`
	Size     int           `json:"size"`
	Skip     int           `json:"skip"`
	Facets   FacetsRequest `json:"facets"`
}

type SearchStatus struct {
	Total      int `json:"total"`
	Failed     int `json:"failed"`
	Successful int `json:"successful"`
}

type FacetResult struct {
	Field         string             `json:"field"`
	Total         int                `json:"total"`
	Missing       int                `json:"missing"`
	Other         int                `json:"other"`
	Terms         TermFacets         `json:"terms,omitempty"`
	NumericRanges NumericRangeFacets `json:"numeric_ranges,omitempty"`
	DateRanges    DateRangeFacets    `json:"date_ranges,omitempty"`
}

type FacetResults map[string]*FacetResult

//
// FacetSearchResults - facet search results
//
type FacetsSearchResults struct {
	Type    string               `json:"type"`
	Status  *SearchStatus        `json:"status"`
	Request *FacetsSearchRequest `json:"request"`
	Hits    []interface{}        `json:"hits"`
	Total   uint64               `json:"total"`
	Took    time.Duration        `json:"took"`
	Facets  FacetResults         `json:"facets"`
	Err     string               `json:"err,omitempty"`
}

// FacetsSearch - searches for query with multiple facets
func FacetsSearch(contentType, language string, req *FacetsSearchRequest, dns string) (*FacetsSearchResults, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/facets")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeFacetsSearchResults).Endpoint()
	req.Type = contentType
	req.Language = language

	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	results := resp.(FacetsSearchResults)
	if results.Err != "" {
		return nil, errors.New(results.Err)
	}
	return &results, nil
}
