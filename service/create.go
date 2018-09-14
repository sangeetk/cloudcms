package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
)

// Create - creates a single item
func (s *Service) Create(ctx context.Context, req *api.CreateRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type}
	var err error

	if _, ok := Index[req.Type]; !ok {
		if err := createIndex(req.Type); err != nil {
			return &resp, nil
		}
	}

	// Create request as sync msg contains full information
	// Simply index the content and return
	if sync {
		IndexLock.Lock()
		defer IndexLock.Unlock()

		var item = (req.Content).(map[string]interface{})
		err = Index[req.Type].Index(item["slug"].(string), item)
		if err != nil {
			return &resp, nil
		}
		resp.Content = req.Content
		return &resp, nil
	}

	// Normal create request
	// Open database in read-write mode
	DB, err = bolt.Open(dbFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	err = DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(req.Type))
		if err != nil {
			return err
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
		IndexLock.Lock()
		defer IndexLock.Unlock()

		err = Index[req.Type].Index(newSlug, item)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		resp.Err = err.Error()
	}

	// Sync others

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
