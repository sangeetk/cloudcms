package service

import (
	"encoding/json"
	"log"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/item"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
)

// Initialize function
func Initialize(dbFile, syncFile string, local, upstream *worker.Worker) error {
	var err error
	var db, syncDB *bolt.DB

	DBFile = dbFile
	SyncFile = syncFile
	LocalWorker = local
	Upstream = upstream

	Index = make(map[string]bleve.Index)

	// Create databse if it doesn't exist.
	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for t := range item.Types {
			_, err := tx.CreateBucketIfNotExists([]byte(t))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return err
	}
	db.Close()

	// Create databse if it doesn't exist.
	syncDB, err = bolt.Open(SyncFile, 0644, nil)
	if err != nil {
		return err
	}

	// Add the (IP:Port, timestamp) to the database
	err = syncDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("workers"))
		if err != nil {
			return err
		}
		err = b.Put([]byte(LocalWorker.String()), nil)
		return err
	})
	if err != nil {
		syncDB.Close()
		return err
	}
	syncDB.Close()

	// Open database again in read-only mode to allow multipe readers
	options := bolt.Options{ReadOnly: true}
	db, err = bolt.Open(DBFile, 0644, &options)
	if err != nil {
		return err
	}
	defer db.Close()

	// Lock the index
	IndexLock.Lock()
	defer IndexLock.Unlock()

	// Initialize index for all Content Types
	for contentType := range item.Types {
		// Create index
		mapping := bleve.NewIndexMapping()
		Index[contentType], err = bleve.NewMemOnly(mapping)
		if err != nil {
			return err
		}

		// Index all available items
		// Access data from within a read-only transactional block.
		if err := db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket([]byte(contentType)).Cursor()
			if c == nil {
				return nil
			}
			// Iterate over the cursor and index the value
			for k, v := c.First(); k != nil; k, v = c.Next() {
				var resp api.Response
				slug := string(k[:])
				err = json.Unmarshal(v, &resp.Content)
				if err != nil {
					return err
				}
				err = Index[contentType].Index(slug, resp.Content)
				if err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	return nil
}
