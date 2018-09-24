package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Delete - creates a single item
func (s *Service) Delete(ctx context.Context, req *api.DeleteRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}
	var db *bolt.DB
	var err error
	log.Println("Delete()", "Type:", req.Type, "Slug:", req.Slug, "Sync:", sync)

	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Forward request to Upstream Server
	if !sync && Upstream.Host != "" {
		log.Println("Forwarding to ", Upstream)
		return LocalWorker.Forward("delete", req, Upstream)
	}

	// Update request as sync msg contains full information
	// Simply index the content and return
	if sync {
		log.Println("Received sync message by ", LocalWorker)

		readReq := api.ReadRequest{Type: req.Type, Slug: req.Slug}
		item, err := s.Read(ctx, &readReq)
		if err != nil {
			resp.Err = api.ErrorNotFound.Error()
			return &resp, nil
		}
		err = Index[req.Type].Delete(req.Slug)
		if err != nil {
			return &resp, nil
		}
		resp.Content = item.Content
		return &resp, nil
	}

	// Open database in read-write mode
	// It will be created if it doesn't exist.
	//options := bolt.Options{ReadOnly: false}
	log.Println("Normal delete request")

	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return api.ErrorNotFound
		}

		// Get the existing value
		val := b.Get([]byte(req.Slug))
		if val == nil {
			return api.ErrorNotFound
		}

		err := json.Unmarshal(val, &resp.Content)
		if err != nil {
			return err
		}

		err = b.Delete([]byte(req.Slug))
		if err != nil {
			return err
		}

		err = Index[req.Type].Delete(req.Slug)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
		return &resp, nil
	}

	// Sync other workers
	sreq := worker.SyncRequest{
		Type:      req.Type,
		Operation: "delete",
		Slug:      req.Slug,
		Timestamp: time.Now().Unix(),
		Source:    LocalWorker.String(),
		Response:  &resp,
	}
	LocalWorker.SyncPeers(SyncFile, &sreq)
	LocalWorker.SyncChilds(SyncFile, &sreq)

	return &resp, nil
}

// DeleteEndpoint - creates endpoint for Delete service
func DeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.DeleteRequest)
		return svc.Delete(ctx, &req, false)
	}
}

// DecodeDeleteReq - decodes the incoming request
func DecodeDeleteReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
