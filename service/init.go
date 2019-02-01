package service

import (
	"encoding/json"
	"errors"
	"log"

	"git.urantiatech.com/cloudcms/cloudcms/api"
	"git.urantiatech.com/cloudcms/cloudcms/item"
	"git.urantiatech.com/cloudcms/cloudcms/worker"
	"github.com/blevesearch/bleve"
	"github.com/boltdb/bolt"
)

// Initialize function
func Initialize(dbFile, syncFile string, local *worker.Worker) error {
	var err error
	var db, syncDB *bolt.DB

	DBFile = dbFile
	SyncFile = syncFile
	LocalWorker = local

	Index = make(map[string]map[string]bleve.Index)

	// Create databse if it doesn't exist.
	db, err = bolt.Open(DBFile, 0644, nil)
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		for t := range item.Types {
			// Create bucket for content type
			b, err := tx.CreateBucketIfNotExists([]byte(t))
			if err != nil {
				return err
			}

			// Create in-memory index for content type
			if _, ok := Index[t]; !ok {
				Index[t] = make(map[string]bleve.Index)
			}

			for _, l := range Languages {
				// Create nested bucket for each supported language
				_, err := b.CreateBucketIfNotExists([]byte(l.String()))
				if err != nil {
					return err
				}
				// Create in-memory index for each supported language
				if _, ok := Index[t][l.String()]; !ok {
					mapping := bleve.NewIndexMapping()
					Index[t][l.String()], err = bleve.NewMemOnly(mapping)
					if err != nil {
						return err
					}
				}
			}

		}
		_, err := tx.CreateBucketIfNotExists([]byte("log"))
		return err
	})
	db.Close()
	if err != nil {
		return err
	}

	// Create databse if it doesn't exist.
	syncDB, err = bolt.Open(SyncFile, 0644, nil)
	if err != nil {
		return err
	}

	// Add the (IP:Port, timestamp) to the database
	err = syncDB.Update(func(tx *bolt.Tx) error {
		var peers *bolt.Bucket

		if peers, err = tx.CreateBucketIfNotExists([]byte("peers")); err != nil {
			return err
		}
		if err = peers.Put([]byte(LocalWorker.String()), nil); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		syncDB.Close()
		return err
	}
	syncDB.Close()

	return nil
}

func RebuildIndex() error {
	// Open database in read-only mode to allow multipe readers
	options := bolt.Options{ReadOnly: true}
	db, err := bolt.Open(DBFile, 0644, &options)
	if err != nil {
		return err
	}
	defer db.Close()

	// Rebuild index for all Content Types
	for t := range item.Types {
		// Index all available items
		if err := db.View(func(tx *bolt.Tx) error {
			// Load all content of all supported languages
			for _, l := range Languages {
				bb, err := getBucket(tx, t, l.String())
				if err != nil {
					return err
				}
				c := bb.Cursor()
				if c == nil {
					return errors.New("Unknown Error")
				}
				// Iterate over the cursor and index the value
				for k, v := c.First(); k != nil; k, v = c.Next() {
					var resp api.Response
					slug := string(k[:])
					err = json.Unmarshal(v, &resp.Content)
					if err != nil {
						return err
					}

					index, err := getIndex(t, l.String())
					if err != nil {
						return err
					}
					item := resp.Content.(map[string]interface{})
					err = index.Index(slug, item)
					if err != nil {
						return err
					}
				}
			}

			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func getBucket(tx *bolt.Tx, contentType, language string) (*bolt.Bucket, error) {
	b := tx.Bucket([]byte(contentType))
	if b == nil {
		return nil, errors.New("Invalid content type")
	}
	bb := b.Bucket([]byte(language))
	if bb == nil {
		return nil, errors.New("Unsupported language")
	}
	return bb, nil
}

func getIndex(contentType, language string) (bleve.Index, error) {
	if _, ok := Index[contentType]; !ok {
		return nil, errors.New("Invalid content type")
	}
	if _, ok := Index[contentType][language]; !ok {
		return nil, errors.New("Unsupported Language")
	}
	return Index[contentType][language], nil
}
