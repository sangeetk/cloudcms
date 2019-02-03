package api

import (
	"time"
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
	Size     int           `json:"size"`
	Skip     int           `json:"skip"`
	Facets   FacetsRequest `json:"facets"`
}

type SearchStatus struct {
	Total      int `json:"total"`
	Failed     int `json:"failed"`
	Successful int `json:"successful"`
}

type TermFacet struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}
type TermFacets []*TermFacet

type FacetResult struct {
	Field          string           `json:"field"`
	Total          int              `json:"total"`
	Missing        int              `json:"missing"`
	Other          int              `json:"other"`
	Terms          TermFacets       `json:"terms,omitempty"`
	NumericRanges  []*NumericRange  `json:"numeric_ranges,omitempty"`
	DateTimeRanges []*DateTimeRange `json:"datetime_ranges,omitempty"`
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
	Total   uint64               `json:"total_hits"`
	Took    time.Duration        `json:"took"`
	Facets  FacetResults         `json:"facets"`
	Err     string               `json:"err,omitempty"`
}
