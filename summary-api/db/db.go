package db

import (
	"context"
	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"main/model"
	"os"
	"time"
)

type CouchDB struct {
	client  *kivik.Client
	couchDB *kivik.DB
}

func NewCouchDB() (CouchDB, error) {
	client, err := kivik.New("couch", os.Getenv("COUCHDB_URL"))
	if err != nil {
		return CouchDB{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	couchdb := client.DB(ctx, os.Getenv("DB_NAME_SUMMARY"))
	if couchdb.Err() != nil {
		return CouchDB{}, couchdb.Err()
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Ping(ctx)
	if err != nil {
		return CouchDB{}, err
	}

	return CouchDB{
		client:  client,
		couchDB: couchdb,
	}, err
}

func (receiver CouchDB) Close(ctx context.Context) error {
	return receiver.client.Close(ctx)
}

func (receiver CouchDB) AddSummary(dialogID, summary string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := uuid.New().String()
	sm := model.Summary{
		ID:        id,
		DialogID:  dialogID,
		Timestamp: time.Now(),
		Summary:   summary,
	}

	_, err := receiver.couchDB.Put(ctx, sm.ID, sm)
	return id, err
}
