package api

import (
	"time"
)

// Header contains some common header fields for content type
type Header struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt int64  `json:"deleted_at"`
}

// Item content type
type Item struct {
	Header
	Content interface{} `json:"content"`
}

// Create - creates a new item
type Create struct {
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Status  string      `json:"status"`
	Content interface{} `json:"content"`
}

// Read - retreives the item from db
type Read struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// Update - updates the item
type Update struct {
	Type    string      `json:"type"`
	Slug    string      `json:"slug"`
	Title   string      `json:"title"`
	Status  string      `json:"status"`
	Content interface{} `json:"content"`
}

// Delete - deletes the item
type Delete struct {
	Type string `json:"type"`
	Slug string `json:"slug"`
}

// Response - response for CRUD requests
type Response struct {
	Item
	Err string `json:"err,omitempty"`
}

// Search - search
type Search struct {
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Query     string    `json:"query"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Limit     int       `json:"limit"`
	Skip      int       `json:"skip"`
}

// SearchResults - search results
type SearchResults struct {
	Results []Item `json:"results"`
	Total   int    `json:"total"`
	Limit   int    `json:"limit"`
	Skip    int    `json:"skip"`
	Err     string `json:"err,omitempty"`
}

// Ping request
type Ping struct {
	Timestamp time.Time `json:"timestamp"`
}

// Pong response
type Pong struct {
	Timestamp1 time.Time `json:"timestamp1"`
	Timestamp2 time.Time `json:"timestamp2"`
}
