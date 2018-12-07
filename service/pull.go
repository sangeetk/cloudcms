package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/urantiatech/kit/endpoint"
)

// Pull request
func (s *Service) Pull(context.Context, *api.PullRequest) (*api.PullResponse, error) {
	var resp api.PullResponse

	return &resp, nil
}

// PullEndpoint - creates endpoint for Pull service
func PullEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.PullRequest)
		return svc.Pull(ctx, &req)
	}
}

// DecodePullReq - decodes the incoming request
func DecodePullReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
