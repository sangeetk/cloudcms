package service

import (
	"context"
	"encoding/json"
	"errors"
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

var workDir string
var dbFile string

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

// Initialize function
func Initialize(file, dir string) error {
	var err error
	dbFile = file

	workDir = dir
	Index = make(map[string]bleve.Index)
	if err := os.MkdirAll(workDir, os.ModeDir|0755); err != nil {
		return err
	}

	// Create databse if it doesn't exist.
	DB, err = bolt.Open(workDir+"/"+dbFile, 0644, nil)
	if err != nil {
		return err
	}
	defer DB.Close()

	// Initialize index for all content types here
	if err := createIndexIfNotPresent("example"); err != nil {
		return err
	}
	return nil
}

func createIndexIfNotPresent(contentType string) error {
	var err error
	if Index[contentType] != nil {
		return nil
	}
	// Initialze the index file
	Index[contentType], err = bleve.Open(workDir + "/" + contentType + ".index")
	if err != nil {
		mapping := bleve.NewIndexMapping()
		Index[contentType], err = bleve.New(workDir+"/"+contentType+".index", mapping)
		if err != nil {
			return err
		}
	}
	return nil
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
