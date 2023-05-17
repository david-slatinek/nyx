package main

import (
	"github.com/roylee0704/gron"
	"log"
	"main/client"
	"main/db"
	"main/env"
	"main/model"
	"main/queue"
	"main/util"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	util.GetCategories()

	g := gron.New()
	g.AddFunc(gron.Every(30*time.Minute), util.GetCategories)
	g.Start()

	recommendDB, err := db.NewRecommendDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(recommendDB db.RecommendDB) {
		if err := recommendDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(recommendDB)

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
			case recommend := <-recommendChannel:
				log.Printf("Got recommend from queue: %v", recommend.Recommend.DialogID)

				dialogs, err := util.GetDialogs(recommend.Recommend.DialogID)
				if err != nil {
					log.Printf("failed to get dialogs: %v", err)
					continue
				}

				recommendResult, err := rClient.GetRecommendationDialogs(dialogs, util.CategoriesText)
				if err != nil {
					log.Printf("failed to get recommendation for dialogs: %v", err)
					continue
				}
				if len(recommendResult) == 0 {
					log.Printf("recommendResult is empty")
					continue
				}

				recommendResultSummary, err := rClient.GetRecommendationSummary(recommend.Recommend, util.CategoriesText)
				if err != nil {
					log.Printf("failed to get recommendation for summary: %v", err)
					continue
				}
				if len(recommendResultSummary) == 0 && len(recommendResult) == 0 {
					log.Printf("recommendResultSummary and recommendResult are empty")
					continue
				}

				length := int(math.Max(float64(len(recommendResult)), float64(len(recommendResultSummary))))
				objects := make(map[string]model.Recommendation, length)

				for _, value := range recommendResult {
					objects[value.Label] = value
				}

				for _, value := range recommendResultSummary {
					obj, ok := objects[value.Label]
					if ok {
						if obj.Score < value.Score {
							objects[value.Label] = value
						}
					} else {
						objects[value.Label] = value
					}
				}

				recommendModel := make([]model.RecommendDB, 0, length)
				for _, value := range objects {
					log.Printf("value: %v", value)

					recommendModel = append(recommendModel, model.RecommendDB{
						UserID:        recommend.Recommend.UserID,
						FkCategory:    util.CategoriesMap[value.Label],
						CategoryName:  value.Label,
						Score:         value.Score,
						FkMainDialog:  recommend.Recommend.DialogID,
						FkDialog:      value.ID,
						RecommendedAt: time.Now(),
					})
				}

				err = recommendDB.Create(recommendModel)
				if err != nil {
					log.Printf("failed to save recommend summary to db: %v", err)
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
				}
			default:
				break
			}
		}
	}()

	log.Println("Recommend API is running")
	<-c
	log.Println("Recommend API is shutting down")
}
