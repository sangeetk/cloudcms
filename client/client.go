package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	ht "github.com/urantiatech/kit/transport/http"
)

// Create - creates a new item
func Create(item *api.Item, dns string) (*api.Response, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/create")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	resp, err := endPoint(ctx, item)
	if err != nil {
		return nil, err
	}
	return &resp.(api.Response), nil
}

// Read - retreives an item from the DB
func Read(item *api.Item, dns string) (*api.Response, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/read")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	resp, err := endPoint(ctx, item)
	if err != nil {
		return nil, err
	}
	return &resp.(api.Response), nil
}

// Update - updated an existing item
func Update(item *api.Item, dns string) (*api.Response, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/update")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	resp, err := endPoint(ctx, item)
	if err != nil {
		return nil, err
	}
	return &resp.(api.Response), nil
}

// Delete - deletes an item from DB
func Delete(item *api.Item, dns string) (*api.Response, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/delete")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	resp, err := endPoint(ctx, item)
	if err != nil {
		return nil, err
	}
	return &resp.(api.Response), nil
}

// Search - searches for query
func Search(item *api.Item, dns string) (*api.SearchResults, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/search")
	if err != nil {
		log.Fatal(err.Error())
	}
	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	resp, err := endPoint(ctx, item)
	if err != nil {
		return nil, err
	}
	return &resp.(api.SearchResults), nil
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
