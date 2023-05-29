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

	couchdb := client.DB(ctx, os.Getenv("DB_NAME_DIALOG"))
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

func (receiver CouchDB) AddDialog(dialog model.Dialog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := receiver.couchDB.Put(ctx, dialog.ID, dialog)
	return err
}

func (receiver CouchDB) GetByDialogID(dialogID string) ([]model.Dialog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	selector := map[string]interface{}{
		"dialogID": dialogID,
	}

	rows, err := receiver.couchDB.Find(ctx, kivik.Options{
		"selector": selector,
	})
	if err != nil {
		return nil, err
	}
	defer func(rows *kivik.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var dialogs = make([]model.Dialog, 0, 10)
	for rows.Next() {
		var dialog model.Dialog
		if err := rows.ScanDoc(&dialog); err != nil {
			log.Printf("failed to scan doc: %v", err)
			continue
		}
		dialogs = append(dialogs, dialog)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dialogs, nil
}

func (receiver CouchDB) GetAll() ([]model.Dialog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := receiver.couchDB.AllDocs(ctx, kivik.Options{"include_docs": true})
	if err != nil {
		return nil, err
	}
	defer func(rows *kivik.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	var dialogs = make([]model.Dialog, 0, 10)
	for rows.Next() {
		var dialog model.Dialog
		if err := rows.ScanDoc(&dialog); err != nil {
			log.Printf("failed to scan doc: %v", err)
			continue
		}
		dialogs = append(dialogs, dialog)
	}

	if rows.Err() != nil {
		return nil, err
	}
	return dialogs, nil
}
