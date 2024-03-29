package main

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	_ "github.com/go-kivik/couchdb/v3"
	"log"
	"main/client"
	"main/controller"
	"main/db"
	"main/env"
	"main/model"
	"main/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	couchDB, err := db.NewCouchDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(client db.CouchDB) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(couchDB)

	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.Error{Error: "not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, model.Error{Error: "method not allowed"})
	})

	grpcClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to connect to grpc: %v", err)
	}
	defer func(client client.Client) {
		if err := client.Close(); err != nil {
			log.Printf("failed to close grpc: %v", err)
		}
	}(grpcClient)

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	if err != nil {
		log.Fatalf("failed to create aws session: %v", err)
	}

	dialogController, err := controller.NewDialogController(couchDB, grpcClient, sqs.New(sess))
	if err != nil {
		log.Fatalf("failed to create dialog controller: %v", err)
	}

	router.Use(util.CORS)

	router.GET("/user", dialogController.UserID)
	router.GET("/dialog", dialogController.DialogID)
	router.POST("/dialog", dialogController.AddDialog)
	router.POST("/end", dialogController.EndDialog)
	router.GET("/dialog/:id", dialogController.GetDialog)
	router.GET("/dialogs", dialogController.GetDialogs)
	router.GET("/user/:dialogID", dialogController.GetUserIDByDialogID)

	srv := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("ListenAndServe() error: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Shutdown() error: %s\n", err)
	}
	log.Println("shutting down dialog-api")
}
