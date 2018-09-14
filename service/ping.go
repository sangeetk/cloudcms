package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/urantiatech/kit/endpoint"
)

// PingRequest request
type PingRequest struct {
	Timestamp time.Time `json:"timestamp"`
}

// PingResponse response
type PingResponse struct {
	Timestamp1 time.Time `json:"timestamp1"`
	Timestamp2 time.Time `json:"timestamp2"`
}

// Ping request
func (s *Service) Ping(ctx context.Context, req *PingRequest) (*PingResponse, error) {
	var resp PingResponse

	resp.Timestamp1 = req.Timestamp
	resp.Timestamp2 = time.Now()

	return &resp, nil
}

// PingEndpoint - creates endpoint for Ping service
func PingEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PingRequest)
		return svc.Ping(ctx, &req)
	}
}

// DecodePingReq - decodes the incoming request
func DecodePingReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request PingRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
