package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Create - creates a single item
func (s *Service) Create(ctx context.Context, req *api.CreateRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}
	var db *bolt.DB
	var err error
	log.Println("Create()", "Type:", req.Type, "Sync:", sync)

	// Validate the content type
	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Forward request to Upstream Server
	if !sync && Upstream.Host != "" {
		log.Println("Forwarding to ", Upstream)
		return LocalWorker.Forward("create", req, Upstream)
	}

	// Sync message
	if sync {
		// Simply index the content and return
		log.Println("Received sync message by ", LocalWorker)

		var item = (req.Content).(map[string]interface{})
		err = Index[req.Type].Index(item["slug"].(string), item)
		if err != nil {
			return &resp, nil
		}
		resp.Content = req.Content
		return &resp, nil
	}

	// Normal request
	// Open database in read-write mode
	log.Println("Normal create request")
	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(req.Type))
		if b == nil {
			return api.ErrorInvalidContentType
		}

		nextSeq, err := b.NextSequence()
		if err != nil {
			return err
		}

		var item = (req.Content).(map[string]interface{})
		item["id"] = nextSeq
		title := item["title"].(string)
		slug := stringToSlug(title)
		item["slug"] = slug
		item["created_at"] = time.Now().Unix()
		item["updated_at"] = time.Now().Unix()
		item["deleted_at"] = 0

		newSlug := slug
		var i int
		// Find a unique slug
		for i = 2; b.Get([]byte(newSlug)) != nil; i++ {
			newSlug = fmt.Sprintf("%s-%d", slug, i)
		}
		if i > 2 {
			item["slug"] = newSlug
		}

		j, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = b.Put([]byte(newSlug), j)
		if err != nil {
			return err
		}

		resp.Content = item

		// Create index
		log.Println("Index the item ", newSlug)

		err = Index[req.Type].Index(newSlug, item)
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
		Operation: "create",
		Timestamp: time.Now().Unix(),
		Source:    LocalWorker.String(),
		Response:  &resp,
	}
	LocalWorker.SyncPeers(SyncFile, &sreq)
	LocalWorker.SyncChilds(SyncFile, &sreq)

	return &resp, nil
}

// CreateEndpoint - creates endpoint for Create service
func CreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(api.CreateRequest)
		return svc.Create(ctx, &req, false)
	}
}

// DecodeCreateReq - decodes the incoming request
func DecodeCreateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var request api.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
