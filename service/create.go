package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	i "git.urantiatech.com/cloudcms/cloudcms/item"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
	"github.com/urantiatech/kit/endpoint"
	"golang.org/x/text/language"
)

// Create - creates a single item
func (s *Service) Create(ctx context.Context, req *api.CreateRequest, sync bool) (*api.Response, error) {
	var resp = api.Response{Type: req.Type, Language: req.Language}
	var db *bolt.DB
	var err error

	// Set the language code
	if req.Language == "" {
		req.Language = language.English.String()
	}
	// Validate the content type
	if _, ok := Index[req.Type]; !ok {
		resp.Err = api.ErrorInvalidContentType.Error()
		return &resp, nil
	}

	// Sync message from peer
	if sync {
		// Simply index the content and return
		var item = (req.Content).(map[string]interface{})
		index, err := getIndex(req.Type, req.Language)
		if err != nil {
			return &resp, nil
		}
		err = index.Index(item["slug"].(string), item)
		if err != nil {
			return &resp, nil
		}
		resp.Content = req.Content
		return &resp, nil
	}

	// Normal request
	// Open database in read-write mode
	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		bb, err := getBucket(tx, req.Type, req.Language)
		if err != nil {
			return err
		}

		nextSeq, err := bb.NextSequence()
		if err != nil {
			return err
		}

		if req.Content == nil {
			return api.ErrorNullContent
		}

		var item = (req.Content).(map[string]interface{})
		item["language"] = req.Language
		item["id"] = nextSeq
		slug := stringToSlug(req.Slug)
		item["slug"] = slug
		item["created_at"] = time.Now().Unix()
		item["updated_at"] = time.Now().Unix()
		item["deleted_at"] = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

		// Copy file(s)
		for k, v := range item {
			if strings.HasPrefix(k, "file:") {
				var file i.File

				if b, err := json.Marshal(v); err != nil {
					return err
				} else if err := json.Unmarshal(b, &file); err != nil {
					return err
				}

				file.URI = fmt.Sprintf("/drive/%s/%s/%d/%s", req.Type, req.Language, nextSeq, file.Name)

				filemap := v.(map[string]interface{})
				filemap["uri"] = file.URI
				filemap["bytes"] = nil

				// Create path
				path := fmt.Sprintf("drive/%s/%s/%d", req.Type, req.Language, nextSeq)
				if err := os.MkdirAll(path, os.ModeDir|os.ModePerm); err != nil {
					return err
				}

				// Create file
				dst, err := os.Create(path + "/" + file.Name)
				if err != nil {
					return err
				}
				defer dst.Close()

				// Copy the uploaded file to the destination file
				buff := bytes.NewBuffer(file.Bytes)
				if _, err := io.Copy(dst, buff); err != nil {
					return err
				}
			}
		}

		// Assign empty status if not provided
		if _, ok := item["status"]; !ok {
			item["status"] = ""
		}

		newSlug := slug
		var i int
		// Find a unique slug
		for i = 2; bb.Get([]byte(newSlug)) != nil; i++ {
			newSlug = fmt.Sprintf("%s-%d", slug, i)
		}
		if i > 2 {
			item["slug"] = newSlug
		}

		j, err := json.Marshal(item)
		if err != nil {
			return err
		}
		err = bb.Put([]byte(newSlug), j)
		if err != nil {
			return err
		}

		resp.Content = item

		// Create index
		var index bleve.Index
		index, err = getIndex(req.Type, req.Language)
		if err != nil {
			return err
		}
		err = index.Index(newSlug, item)
		if err != nil {
			return err
		}

		// Add the request to LOG bucket
		l := tx.Bucket([]byte("log"))
		seq, err := l.NextSequence()
		if err != nil {
			return err
		}
		event := api.Event{Seq: seq, Op: api.CreateOp, Request: req}
		e, err := json.Marshal(event)
		if err != nil {
			return err
		}

		var key = make([]byte, 8)
		binary.LittleEndian.PutUint64(key, seq)
		err = l.Put(key, e)
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
