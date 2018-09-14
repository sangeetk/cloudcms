package api

import (
	"time"
)

// CreateRequest structure
type CreateRequest struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

// UpdateRequest structure
type UpdateRequest struct {
	Type    string      `json:"type"`
	Slug    string      `json:"slug"`
	Content interface{} `json:"content"`
}

// DeleteRequest structure
type DeleteRequest struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// ReadRequest structure
type ReadRequest struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// Response structure
type Response struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
	Err     string      `json:"err,omitempty"`
}

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
