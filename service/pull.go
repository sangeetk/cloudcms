package service

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Pull request
func (s *Service) Pull(ctx context.Context, req *api.PullRequest) (*api.PullResponse, error) {
	var resp api.PullResponse
	var db *bolt.DB
	var err error

	// Open database in read-only mode
	options := bolt.Options{ReadOnly: true}
	db, err = bolt.Open(DBFile, 0644, &options)
	if err != nil {
		resp.Err = api.ErrorNotFound.Error()
		return &resp, nil
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		l := tx.Bucket([]byte("log"))
		if l == nil {
			return api.ErrorNotFound
		}
		c := l.Cursor()
		startKey := make([]byte, 8)
		binary.LittleEndian.PutUint64(startKey, req.Seq+1)

		count := req.Count
		for k, v := c.Seek(startKey); k != nil && count > 0; k, v = c.Next() {
			var e api.Event
			if err := json.Unmarshal(v, &e); err != nil {
				return err
			}
			resp.Events = append(resp.Events, e)
			count--
		}
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
	}

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
