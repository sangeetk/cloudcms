package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/urantiatech/kit/endpoint"
)

// Ping request
func (s *Service) Ping(ctx context.Context, req *api.Ping) (*api.Pong, error) {
	var resp api.Pong

	resp.Timestamp1 = req.Timestamp
	resp.Timestamp2 = time.Now()

	return &resp, nil
}

// PingEndpoint - creates endpoint for Ping service
func PingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.Ping)
		return svc.Ping(ctx, &req)
	}
}

// DecodePingReq - decodes the incoming request
func DecodePingReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.Ping
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
