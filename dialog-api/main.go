package main

import (
	"context"
	"log"
	"main/db"
	"main/env"
	"os"

	_ "github.com/go-kivik/couchdb/v3" // The couchDB driver
)

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	uri := os.Getenv("COUCHDB_URL")
	if uri == "" {
		log.Fatal("COUCHDB_URL is not set")
	}

	client, err := db.NewCouchDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(client db.CouchDB) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := client.Close(ctx); err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}(client)
}
