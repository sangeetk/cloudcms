package api

import (
	"context"
	"errors"
	"log"
	"net/url"

	"git.urantiatech.com/cloudcms/cloudcms/item"
	ht "github.com/urantiatech/kit/transport/http"
)

// SchemaRequest - schema request
type SchemaRequest struct {
	Type string `json:"type"`
}

// ContentType is the type of content stored in cms
type ContentType struct {
	Fields []item.Field `json:"fields"`
}

// SchemaResponse - schema response
type SchemaResponse struct {
	Languages []string               `json:"languages"`
	Schema    map[string]ContentType `json:"schema,omitempty"`
	Err       string                 `json:"err,omitempty"`
}

// Schema - fetches the info about content types
func Schema(dns string) ([]string, map[string]ContentType, error) {
	ctx := context.Background()
	tgt, err := url.Parse("http://" + dns + "/schema")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", tgt, encodeRequest, decodeSchemaResponse).Endpoint()
	req := SchemaRequest{}
	resp, err := endPoint(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if resp.(SchemaResponse).Err != "" {
		return nil, nil, errors.New(resp.(SchemaResponse).Err)
	}
	return resp.(SchemaResponse).Languages, resp.(SchemaResponse).Schema, nil
}
