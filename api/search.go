package api

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"

	ht "github.com/urantiatech/kit/transport/http"
)

// SearchRequest - search request
type SearchRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Query    string `json:"query"`
	Fuzzy    bool   `json:fuzzy"`
	Size     int    `json:"size"`
	Skip     int    `json:"skip"`
}

// SearchResults - search results
type SearchResults struct {
	Type    string         `json:"type"`
	Request *SearchRequest `json:"request"`
	Hits    []interface{}  `json:"hits"`
	Total   uint64         `json:"total_hits"`
	Took    time.Duration  `json:"took"`
	Err     string         `json:"err,omitempty"`
}

// Search - searches for query
func Search(contentType, language, query string, fuzzy bool, size, skip int, dns string) (
	[]interface{}, uint64, time.Duration, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/search")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeSearchResults).Endpoint()
	req := SearchRequest{
		Type:     contentType,
		Language: language,
		Query:    query,
		Fuzzy:    fuzzy,
		Size:     size,
		Skip:     skip,
	}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, 0, 0, err
	}
	if resp.(SearchResults).Err != "" {
		return nil, 0, 0, errors.New(resp.(SearchResults).Err)
	}
	r := resp.(SearchResults)
	return r.Hits, r.Total, r.Took, nil
}
