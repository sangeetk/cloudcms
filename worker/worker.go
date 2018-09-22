package worker

import (
	"fmt"
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
	return fmt.Sprintf("%s:%d", w.Host, w.Port)
}
