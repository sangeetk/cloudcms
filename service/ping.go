package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/urantiatech/kit/endpoint"
)

// Ping request
func (s *Service) Ping(ctx context.Context, req *worker.PingRequest) (*worker.PingResponse, error) {
	var resp worker.PingResponse

	resp.Timestamp1 = req.Timestamp
	resp.Timestamp2 = time.Now()

	return &resp, nil
}

// PingEndpoint - creates endpoint for Ping service
func PingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(worker.PingRequest)
		return svc.Ping(ctx, &req)
	}
}

// DecodePingReq - decodes the incoming request
func DecodePingReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request worker.PingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
