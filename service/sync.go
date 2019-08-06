package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/go-kit/kit/endpoint"
)

// LastSyncTimestamp in nanoseconds
var LastSyncTimestamp int64

// Sync request
func (s *Service) Sync(ctx context.Context, req *worker.SyncRequest) (*worker.SyncResponse, error) {
	var syncResp worker.SyncResponse
	var content = req.Response.Content
	var resp *api.Response

	// Call Create, Update or Delete service based on Operation
	switch req.Operation {
	case api.CreateOp:
		createReq := api.CreateRequest{
			Type:    req.Type,
			Content: content,
		}
		resp, _ = s.Create(ctx, &createReq, true)
	case api.UpdateOp:
		updateReq := api.UpdateRequest{
			Type:    req.Type,
			Slug:    req.Slug,
			Content: content,
		}
		resp, _ = s.Update(ctx, &updateReq, true)
	case api.DeleteOp:
		deleteReq := api.DeleteRequest{
			Type: req.Type,
			Slug: req.Slug,
		}
		resp, _ = s.Delete(ctx, &deleteReq, true)
	default:
		syncResp.Err = api.ErrorInvalidOperation.Error()
		return &syncResp, nil
	}

	syncResp.Operation = req.Operation
	syncResp.Timestamp = time.Now().UnixNano()

	if resp.Err != "" {
		syncResp.Err = api.ErrorSync.Error()
	} else {
		syncResp.Response = resp
	}
	return &syncResp, nil
}

// SyncEndpoint - creates endpoint for Sync service
func SyncEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(worker.SyncRequest)
		return svc.Sync(ctx, &req)
	}
}

// DecodeSyncReq - decodes the incoming request
func DecodeSyncReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request worker.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
