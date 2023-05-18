package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"main/controller"
	"main/db"
	"main/env"
	"main/model"
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

	recommendFollowDB, err := db.NewRecommendFollow()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(recommendDB db.RecommendFollow) {
		if err := recommendDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(recommendFollowDB)

	recommendFollowController := controller.NewRecommendFollow(recommendFollowDB)

	router := gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, model.Error{Error: "not found"})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, model.Error{Error: "method not allowed"})
	})

	router.LoadHTMLFiles("index.html")

	router.GET("/recommend/:id", recommendFollowController.AddRecommendFollow)

	srv := &http.Server{
		Addr:         ":8020",
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
	log.Println("shutting down tracking-api")
}
