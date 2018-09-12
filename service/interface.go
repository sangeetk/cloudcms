package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
)

// DB handler for the bolt database file
var DB *bolt.DB

// DefaultBucket name
const DefaultBucket = "default"

// indexDir - Directory to store index files
var indexDir string

// Index Map
var Index map[string]bleve.Index

// Interface definition
type Interface interface {
	Create(context.Context, *api.Create) (*api.Response, error)
	Read(context.Context, *api.Read) (*api.Response, error)
	Update(context.Context, *api.Update) (*api.Response, error)
	Delete(context.Context, *api.Delete) (*api.Response, error)
	Search(context.Context, *api.Search) (*api.SearchResults, error)
	Ping(context.Context, *api.Ping) (*api.Pong, error)
}

// InitIndexMap initializes the map for storing indexes for content types
func InitIndexMap(dir string) error {
	indexDir = dir
	Index = make(map[string]bleve.Index)
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		return err
	}

	// Initialize index for all content types here
	createIndexIfNotPresent("example")
	return nil
}

func createIndexIfNotPresent(contentType string) {
	var err error
	if Index[contentType] != nil {
		return
	}
	// Initialze the index file
	Index[contentType], err = bleve.Open(indexDir + "/" + contentType + ".index")
	if err != nil {
		mapping := bleve.NewIndexMapping()
		Index[contentType], err = bleve.New(indexDir+"/"+contentType+".index", mapping)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Service struct for accessing services
type Service struct{}

// ErrorNotFound - 404 Not Found
var ErrorNotFound = errors.New("Not Found")

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
