package service

import (
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/blevesearch/bleve"
)

// DBFile is the path of database file
var DBFile string

// SyncFile is the path of sync file
var SyncFile string

// LocalWorker is the current worker process
var LocalWorker *worker.Worker

// DefaultBucket name
const DefaultBucket = "default"

// Index Map
var Index map[string]bleve.Index
