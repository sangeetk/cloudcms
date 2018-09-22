package worker

import (
	"context"
	"log"
	"net/url"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/client"
	ht "github.com/urantiatech/kit/transport/http"
)

// Forward forwards the request to another worker
func (w *Worker) Forward(op string, req interface{}, upstream *Worker) (*api.Response, error) {
	var resp interface{}
	var err error

	ctx := context.Background()
	up, err := url.Parse("https://" + upstream.String() + "/create")
	if err != nil {
		log.Fatal(err.Error())
	}

	endPoint := ht.NewClient("POST", up, client.EncodeRequest, client.DecodeResponse).Endpoint()

	switch op {
	case "create":
		creq := req.(api.CreateRequest)
		resp, err = endPoint(ctx, creq)
	case "update":
		ureq := req.(api.UpdateRequest)
		resp, err = endPoint(ctx, ureq)
	case "delete":
		dreq := req.(api.DeleteRequest)
		resp, err = endPoint(ctx, dreq)
	default:
		return nil, api.ErrorInvalidOperation
	}

	if err != nil {
		return nil, err
	}
	response := resp.(api.Response)
	return &response, nil
}
