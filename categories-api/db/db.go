package db

import (
	"context"
	"github.com/go-kivik/kivik/v3"
	"log"
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

	couchdb := client.DB(ctx, os.Getenv("DB_NAME_CATEGORIES"))
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

func (receiver CouchDB) AddCategory(category model.Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := receiver.couchDB.Put(ctx, category.ID, category)
	return err
}

func (receiver CouchDB) GetCategory(id string) (model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := receiver.couchDB.Get(ctx, id)

	if row.Err != nil {
		return model.Category{}, row.Err
	}

	var category model.Category
	if err := row.ScanDoc(&category); err != nil {
		return model.Category{}, err
	}
	return category, nil
}

func (receiver CouchDB) GetCategories() ([]model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := receiver.couchDB.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer func(rows *kivik.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		if err := rows.ScanDoc(&category); err != nil {
			continue
		}
		categories = append(categories, category)
	}
	return categories, nil
}
