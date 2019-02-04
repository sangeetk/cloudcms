package api

import (
	"context"
	"errors"
	"log"
	"net/url"

	ht "github.com/urantiatech/kit/transport/http"
)

// ListRequest - list request
type ListRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Status   string `json:"status"`
	SortBy   string `json:"sortby"`
	Size     int    `json:"size"`
	Skip     int    `json:"skip"`
}

// ListResults - list results
type ListResults struct {
	Type    string        `json:"type"`
	Request *ListRequest  `json:"request"`
	List    []interface{} `json:"list"`
	Total   uint64        `json:"total"`
	Err     string        `json:"err,omitempty"`
}

// List - list all items
func List(contentType, language, status, sortby string, size, skip int, dns string) (
	[]interface{}, uint64, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/list")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeListResults).Endpoint()
	req := ListRequest{
		Type:     contentType,
		Language: language,
		Status:   status,
		SortBy:   sortby,
		Size:     size,
		Skip:     skip,
	}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.(ListResults).Err != "" {
		return nil, 0, errors.New(resp.(ListResults).Err)
	}
	r := resp.(ListResults)
	return r.List, r.Total, nil
}
