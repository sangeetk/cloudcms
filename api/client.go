package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	ht "github.com/urantiatech/kit/transport/http"
)

// Create - creates a new item
func Create(contentType, language, slug string, content interface{}, dns string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/create")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeResponse).Endpoint()
	req := CreateRequest{Type: contentType, Language: language, Slug: slug, Content: content}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.(Response).Err != "" {
		return nil, errors.New(resp.(Response).Err)
	}
	return resp.(Response).Content, nil
}

// Read - retreives an item from the DB
func Read(contentType, language, slug string, dns string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/read")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeResponse).Endpoint()
	req := ReadRequest{Type: contentType, Language: language, Slug: slug}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.(Response).Err != "" {
		return nil, errors.New(resp.(Response).Err)
	}
	return resp.(Response).Content, nil
}

// Update - updated an existing item
func Update(contentType, language, slug string, content interface{}, dns string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/update")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeResponse).Endpoint()
	req := UpdateRequest{Type: contentType, Language: language, Slug: slug, Content: content}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.(Response).Err != "" {
		return nil, errors.New(resp.(Response).Err)
	}
	return resp.(Response).Content, nil
}

// Delete - deletes an item from DB
func Delete(contentType, language, slug string, dns string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/delete")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeResponse).Endpoint()
	req := DeleteRequest{Type: contentType, Language: language, Slug: slug}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.(Response).Err != "" {
		return nil, errors.New(resp.(Response).Err)
	}
	return resp.(Response).Content, nil
}

// Search - searches for query
func Search(contentType, language, query string, startDate, endDate time.Time, limit, skip int, dns string) (
	[]interface{}, int, int, int, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/search")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeSearchResults).Endpoint()
	req := SearchRequest{
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
	if resp.(SearchResults).Err != "" {
		return nil, 0, 0, 0, errors.New(resp.(SearchResults).Err)
	}
	r := resp.(SearchResults)
	return r.Results, r.Total, r.Limit, r.Skip, nil
}

// List - list all items
func List(contentType, language, status string, limit, skip int, dns string) (
	[]interface{}, int, int, int, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/list")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, EncodeRequest, DecodeListResults).Endpoint()
	req := ListRequest{
		Type:   contentType,
		Status: status,
		Limit:  limit,
		Skip:   skip,
	}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	if resp.(ListResults).Err != "" {
		return nil, 0, 0, 0, errors.New(resp.(ListResults).Err)
	}
	r := resp.(ListResults)
	return r.List, r.Total, r.Limit, r.Skip, nil
}

// EncodeRequest encodes the request as JSON
func EncodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// DecodeResponse decodes the response from the service
func DecodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response Response
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

// DecodeSearchResults decodes search results
func DecodeSearchResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results SearchResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

// DecodeListResults decodes list results
func DecodeListResults(ctx context.Context, r *http.Response) (interface{}, error) {
	var results ListResults
	if err := json.NewDecoder(r.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}
