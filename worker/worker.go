package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"git.urantiatech.com/cloudcms/cloudcms/api"
)

// Worker details
type Worker struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// PingRequest request
type PingRequest struct {
	Timestamp time.Time `json:"timestamp"`
}

// PingResponse response
type PingResponse struct {
	Timestamp1 time.Time `json:"timestamp1"`
	Timestamp2 time.Time `json:"timestamp2"`
}

// SyncRequest sync request
type SyncRequest struct {
	Type      string        `json:"type"`
	Slug      string        `json:"slug"`
	Operation string        `json:"operation"`
	Timestamp int64         `json:"timestamp"`
	Source    string        `json:"source"`
	Response  *api.Response `json:"response"`
}

// SyncResponse sync response
type SyncResponse struct {
	Type      string        `json:"type"`
	Operation string        `json:"operation"`
	Timestamp int64         `json:"timestamp"`
	Response  *api.Response `json:"response"`
	Err       string        `json:"err,omitempty"`
}

func (w *Worker) String() string {
	return fmt.Sprintf("%s:%d", w.Host, w.Port+1)
}

// encodeRequest encodes the request as JSON
func encodeRequest(ctx context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// decodeResponse decodes the response from the service
func decodeResponse(ctx context.Context, r *http.Response) (interface{}, error) {
	var response SyncResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}
