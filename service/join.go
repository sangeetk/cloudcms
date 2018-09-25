package service

import (
	"context"
	"encoding/json"
	// "log"
	"net/http"
	// "time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	// "git.urantiatech.com/cloudcms/cloudcms/worker"
	// "github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Join - allow other servers/clusters to join the master
func (s *Service) Join(ctx context.Context, req *api.JoinRequest, sync bool) (*api.JoinResponse, error) {
	var resp = api.JoinResponse{}

	return &resp, nil
}

// JoinEndpoint - creates endpoint for Join service
func JoinEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.JoinRequest)
		return svc.Join(ctx, &req, false)
	}
}

// DecodeJoinReq - decodes the incoming request
func DecodeJoinReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
