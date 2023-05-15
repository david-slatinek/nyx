package main

import (
	"log"
	"main/client"
	"main/env"
	"main/model"
	"main/queue"
	"main/util"
	"os"
	"os/signal"
	"syscall"
)

const QueueRecommend = "recommend-q"
const QueueEmail = "email-q"

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	recommendQueue, err := queue.NewQueue(QueueRecommend)
	if err != nil {
		log.Fatalf("failed to get %s queue url: %v", QueueRecommend, err)
	}

	emailQueue, err := queue.NewQueue(QueueEmail)
	if err != nil {
		log.Fatalf("failed to get %s queue url: %v", QueueEmail, err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	recommendChannel := make(chan model.Result, 1)

	rClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	defer func(rClient client.Client) {
		if err := rClient.Close(); err != nil {
			log.Printf("failed to close client: %v", err)
		}
	}(rClient)

	go func() {
		for {
			select {
			case <-c:
				return
			default:
				recommend, deleteM, err := recommendQueue.Receive()
				if err != nil {
					log.Printf("failed to receive from queue: %v", err)
					continue
				}
				recommendChannel <- model.Result{
					Recommend: recommend,
					Delete:    deleteM,
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-c:
				return
			default:
				break
			}

			recommend := <-recommendChannel

			log.Printf("Got recommend from queue: %v", recommend.Recommend.DialogID)

			dialogs, err := util.GetDialogs(recommend.Recommend.DialogID)
			if err != nil {
				log.Printf("failed to get dialogs: %v", err)
				continue
			}

			categories, err := util.GetCategories()
			if err != nil {
				log.Printf("failed to get categories: %v", err)
				continue
			}

			var categoriesText = make([]string, 0, len(categories))
			util.GetCategoriesNames(categories, &categoriesText)

			recommendResult, err := rClient.GetRecommendationDialogs(dialogs, categoriesText)
			if err != nil {
				log.Printf("failed to get recommendation for dialogs: %v", err)
				continue
			}

			if len(recommendResult) == 0 {
				log.Printf("recommendResult is empty")
				continue
			}

			err = emailQueue.Send(recommendResult)
			if err != nil {
				log.Printf("failed to send to queue: %v", err)
			} else {
				log.Printf("recommendResult sent")
			}

			recommendResultSummary, err := rClient.GetRecommendationSummary(recommend.Recommend, categoriesText)
			if err != nil {
				log.Printf("failed to get recommendation for summary: %v", err)
				continue
			}
			if len(recommendResultSummary) == 0 {
				log.Printf("recommendResultSummary is empty")
				continue
			}

			err = emailQueue.Send(recommendResultSummary)
			if err != nil {
				log.Printf("failed to send to queue: %v", err)
				continue
			} else {
				log.Printf("recommendResultSummary sent")
			}

			err = recommend.Delete()
			if err != nil {
				log.Printf("failed to delete message: %v", err)
				continue
			}
			log.Printf("recommendation sent")
		}
	}()

	log.Println("Recommend API is running")
	<-c
	log.Println("Recommend API is shutting down")
}
