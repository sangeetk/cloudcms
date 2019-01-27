package service

import (
	"context"
	"encoding/json"
	"net/http"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
)

// Interface definition
type Interface interface {
	// Normal DB operations
	Create(context.Context, *api.CreateRequest, bool) (*api.Response, error)
	Read(context.Context, *api.ReadRequest) (*api.Response, error)
	Update(context.Context, *api.UpdateRequest, bool) (*api.Response, error)
	Delete(context.Context, *api.DeleteRequest, bool) (*api.Response, error)
	Search(context.Context, *api.SearchRequest) (*api.SearchResults, error)
	List(context.Context, *api.ListRequest) (*api.ListResults, error)

	// Schema request from admin interface
	Schema(context.Context, *api.SchemaRequest) (*api.SchemaResponse, error)

	// Pull request from downstream connections
	Pull(context.Context, *api.PullRequest) (*api.PullResponse, error)
	// Push request from upstream server
	Push(context.Context, *api.PushRequest) (*api.PushResponse, error)

	// Only between peer-to-peer communication
	Ping(context.Context, *worker.PingRequest) (*worker.PingResponse, error)
	Sync(context.Context, *worker.SyncRequest) (*worker.SyncResponse, error)
}

// Service struct for accessing services
type Service struct{}

// Encode the response
func Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
