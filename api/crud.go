package api

import (
	"context"
	"errors"
	"log"
	"net/url"

	ht "github.com/urantiatech/kit/transport/http"
)

const (
	// CreateOp opration
	CreateOp = "create"
	// ReadOp operation
	ReadOp = "read"
	// UpdateOp operation
	UpdateOp = "update"
	// DeleteOp operation
	DeleteOp = "delete"
)

// CreateRequest structure
type CreateRequest struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Slug     string      `json:"slug"`
	SlugText string      `json:"slug_text"`
	Content  interface{} `json:"content"`
}

// UpdateRequest structure
type UpdateRequest struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Slug     string      `json:"slug"`
	Content  interface{} `json:"content"`
}

// DeleteRequest structure
type DeleteRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Slug     string `json:"slug"`
}

// ReadRequest structure
type ReadRequest struct {
	Type     string `json:"type"`
	Language string `json:"language"`
	Slug     string `json:"slug"`
}

// Response structure
type Response struct {
	Type     string      `json:"type"`
	Language string      `json:"language"`
	Content  interface{} `json:"content"`
	Err      string      `json:"err,omitempty"`
}

// Create - creates a new item
func Create(contentType, language, slug, slugtext string, content interface{}, dns string) (interface{}, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/create")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
	req := CreateRequest{Type: contentType, Language: language, Slug: slug, SlugText: slugtext, Content: content}
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

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
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

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
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

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeResponse).Endpoint()
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
