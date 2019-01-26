package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Push request
func (s *Service) Push(ctx context.Context, req *api.PushRequest) (*api.PushResponse, error) {
	var resp api.PushResponse
	var r *api.Response
	var seqNum uint64

	// Check for out of sync events
	options := bolt.Options{ReadOnly: true}
	db, err := bolt.Open(DBFile, 0644, &options)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}

	err = db.View(func(tx *bolt.Tx) error {
		l := tx.Bucket([]byte("log"))
		if l == nil {
			return api.ErrorNotFound
		}
		seqNum = l.Sequence()
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
		return &resp, nil
	}
	db.Close()

	log.Println("Seqnum :", seqNum)

	if req.Seq != seqNum+1 {
		resp.Err = api.ErrorOutOfSync.Error()
		return &resp, nil
	}

	// Call Create, Update or Delete service based on Operation
	switch req.Op {
	case api.CreateOp:
		c := req.Request.(api.CreateRequest)
		r, _ = s.Create(ctx, &c, false)

	case api.UpdateOp:
		u := req.Request.(api.UpdateRequest)
		r, _ = s.Update(ctx, &u, false)

	case api.DeleteOp:
		d := req.Request.(api.DeleteRequest)
		r, _ = s.Delete(ctx, &d, false)

	default:
		resp.Err = api.ErrorInvalidOperation.Error()
	}
	resp.Err = r.Err
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
