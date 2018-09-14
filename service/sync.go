package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/urantiatech/kit/endpoint"
)

// LastSyncTimestamp in nanoseconds
var LastSyncTimestamp int64

// ErrorInvalidOperation - invalid operation
var ErrorInvalidOperation = errors.New("Invalid Operation")

// SyncRequest sync request
type SyncRequest struct {
	Operation string             `json:"operation"`
	Timestamp int64              `json:"timestamp"`
	Source    string             `json:"source"`
	Create    *api.CreateRequest `json:"create"`
	Update    *api.UpdateRequest `json:"update"`
	Delete    *api.DeleteRequest `json:"delete"`
}

// SyncResponse sync response
type SyncResponse struct {
	Operation string        `json:"operation"`
	Timestamp int64         `json:"timestamp"`
	Response  *api.Response `json:"response"`
	Err       string        `json:"err,omitempty"`
}

// Sync request
func (s *Service) Sync(ctx context.Context, req *SyncRequest) (*SyncResponse, error) {
	var syncResp SyncResponse
	var resp *api.Response

	// Call Create, Update or Delete service based on Operation
	switch req.Operation {
	case "create":
		resp, _ = s.Create(ctx, req.Create, true)
	case "update":
		resp, _ = s.Update(ctx, req.Update, true)
	case "delete":
		resp, _ = s.Delete(ctx, req.Delete, true)
	default:
		syncResp.Err = ErrorInvalidOperation.Error()
		return &syncResp, nil
	}

	syncResp.Operation = req.Operation
	syncResp.Timestamp = time.Now().UnixNano()

	if resp.Err != "" {
		syncResp.Err = resp.Err
	} else {
		syncResp.Response = resp
	}
	return &syncResp, nil
}

// SyncEndpoint - creates endpoint for Sync service
func SyncEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SyncRequest)
		return svc.Sync(ctx, &req)
	}
}

// DecodeSyncReq - decodes the incoming request
func DecodeSyncReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
