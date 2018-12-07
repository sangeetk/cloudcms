package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/urantiatech/kit/endpoint"
)

// Push request
func (s *Service) Push(context.Context, *api.PushRequest) (*api.PushResponse, error) {
	var resp api.PushResponse

	return &resp, nil
}

// PushEndpoint - creates endpoint for Push service
func PushEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.PushRequest)
		return svc.Push(ctx, &req)
	}
}

// DecodePushReq - decodes the incoming request
func DecodePushReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.PushRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
