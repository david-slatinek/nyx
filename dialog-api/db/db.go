package db

import (
	"context"
	"github.com/go-kivik/kivik/v3"
	"main/model"
	"os"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	couchdb := client.DB(ctx, os.Getenv("DB_NAME"))
	if couchdb.Err() != nil {
		return CouchDB{}, couchdb.Err()
	}

	ctx, cancel = context.WithCancel(context.Background())
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

func (receiver CouchDB) AddDocument(chat model.Chat) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := receiver.couchDB.Put(ctx, chat.ID, chat)
	if err != nil {
		return err
	}
	return nil
}
