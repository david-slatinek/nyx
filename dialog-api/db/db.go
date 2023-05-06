package db

import (
	"context"
	"github.com/go-kivik/kivik/v3"
	"os"
)

type CouchDB struct {
	client *kivik.Client
}

func NewCouchDB() (CouchDB, error) {
	client, err := kivik.New("couch", os.Getenv("COUCHDB_URL"))
	if err != nil {
		return CouchDB{}, err
	}
	return CouchDB{client: client}, err
}

func (receiver CouchDB) Close(ctx context.Context) error {
	return receiver.client.Close(ctx)
}
