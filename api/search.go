package api

import (
	"time"
)

// SearchRequest - search request
type SearchRequest struct {
	Type      string    `json:"type"`
	Query     string    `json:"query"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
	Skip      int       `json:"skip"`
}

// SearchResults - search results
type SearchResults struct {
	Type    string        `json:"type"`
	Results []interface{} `json:"results"`
	Total   int           `json:"total"`
	Limit   int           `json:"limit"`
	Skip    int           `json:"skip"`
	Err     string        `json:"err,omitempty"`
}