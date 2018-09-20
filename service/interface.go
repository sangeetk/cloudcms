package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"sync"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
)

// DefaultBucket name
const DefaultBucket = "default"

// Index Map
var Index map[string]bleve.Index

// IndexLock mutex
var IndexLock sync.Mutex

// Interface definition
type Interface interface {
	Create(context.Context, *api.CreateRequest, bool) (*api.Response, error)
	Read(context.Context, *api.ReadRequest) (*api.Response, error)
	Update(context.Context, *api.UpdateRequest, bool) (*api.Response, error)
	Delete(context.Context, *api.DeleteRequest, bool) (*api.Response, error)
	Search(context.Context, *api.SearchRequest) (*api.SearchResults, error)

	// Only between peer-to-peer communication
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	Sync(context.Context, *SyncRequest) (*SyncResponse, error)
}

// Service struct for accessing services
type Service struct{}

// ErrorNotFound - 404 Not Found
var ErrorNotFound = errors.New("Not Found")

// ErrorInvalidContentType -
var ErrorInvalidContentType = errors.New("Invalid ContentType")

// Encode the response
func Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func stringToSlug(title string) string {
	// Filter and conver to lowercase
	slug := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r + 'a' - 'A'
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		}
		return ' '
	}

	// Convert whitespace to hyphen '-'
	str := strings.Map(slug, title)
	strarray := strings.Fields(str)
	str = strings.Join(strarray, "-")

	return str
}
