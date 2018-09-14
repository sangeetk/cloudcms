package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	ht "github.com/urantiatech/kit/transport/http"
)

// Client structure
type Client struct {
	DNS string
}

// New - create a new Client
func (c *Client) New(dns string) {
	c.DNS = dns
}

// Create - creates a new item
func (c *Client) Create(contentType string, content interface{}) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + c.DNS + "/create")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := api.CreateRequest{Type: contentType, Content: content}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(api.Response).Content, nil
}

// Read - retreives an item from the DB
func (c *Client) Read(contentType, slug string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + c.DNS + "/read")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := api.ReadRequest{Type: contentType, Slug: slug}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(api.Response).Content, nil
}

// Update - updated an existing item
func (c *Client) Update(contentType, slug string, content interface{}) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + c.DNS + "/update")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := api.UpdateRequest{Type: contentType, Slug: slug, Content: content}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(api.Response).Content, nil
}

// Delete - deletes an item from DB
func (c *Client) Delete(contentType, slug string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + c.DNS + "/delete")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := api.DeleteRequest{Type: contentType, Slug: slug}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(api.Response).Content, nil
}

// Search - searches for query
func (c *Client) Search(contentType, query string, startDate, endDate time.Time, limit, skip int) (
	[]interface{}, int, int, int, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + c.DNS + "/search")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := api.SearchRequest{
		Type:      contentType,
		Query:     query,
		StartDate: startDate,
		EndDate:   endDate,
		Limit:     limit,
		Skip:      skip,
	}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	r := resp.(api.SearchResults)
	return r.Results, r.Total, r.Limit, r.Skip, nil
}

func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response api.Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func decodeSearchResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results api.SearchResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}
