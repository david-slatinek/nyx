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

	//mainID := util.GetMainCategoryID(util.Categories, "Dell")
	//log.Printf("mainID: %v", mainID)

	recommendDB, err := db.NewRecommendDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(recommendDB db.RecommendDB) {
		if err := recommendDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}(recommendDB)

	recommendModel := make([]model.RecommendDB, 2)
	recommendModel[0] = model.RecommendDB{
		UserID:        "1",
		FkCategory:    "1",
		CategoryName:  "pc",
		Score:         0.5,
		FkMainDialog:  "1",
		FkDialog:      "1",
		RecommendedAt: time.Now(),
	}

	recommendModel[1] = model.RecommendDB{
		UserID:        "1",
		FkCategory:    "1",
		CategoryName:  "pc",
		Score:         0.5,
		FkMainDialog:  "2",
		FkDialog:      "2",
		RecommendedAt: time.Now(),
	}

	err = recommendDB.Create(recommendModel)
	if err != nil {
		log.Fatalf("failed to create recommend: %v", err)
	}

	log.Printf("id: %v", recommendModel[0].ID)

	os.Exit(0)

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

	//recommendDB, err := db.NewRecommendDB()
	//if err != nil {
	//	log.Fatalf("failed to connect to db: %v", err)
	//}
	//defer func(recommendDB db.RecommendDB) {
	//	if err := recommendDB.Close(); err != nil {
	//		log.Printf("failed to close db: %v", err)
	//	}
	//}(recommendDB)
	//
	//recommendModel := make([]model.RecommendDB, 2)
	//recommendModel[0] = model.RecommendDB{
	//	UserID:       "1",
	//	FkCategory:   "1",
	//	CategoryName: "pc",
	//	Score:        0.5,
	//	FkMainDialog: "1",
	//	FkDialog:     "1",
	//}
	//
	//recommendModel[1] = model.RecommendDB{
	//	UserID:       "1",
	//	FkCategory:   "1",
	//	CategoryName: "pc",
	//	Score:        0.5,
	//	FkMainDialog: "2",
	//	FkDialog:     "2",
	//}

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

			recommendResult, err := rClient.GetRecommendationDialogs(dialogs, util.CategoriesText)
			if err != nil {
				log.Printf("failed to get recommendation for dialogs: %v", err)
				continue
			}

			//recommendModel := make([]model.RecommendDB, 0)
			//for key, r := range recommendResult {
			//	recommendModel[key] = model.RecommendDB{
			//		UserID:       recommend.Recommend.UserID,
			//		FkCategory:   util.GetMainCategoryID(util.Categories, r.Label),
			//		CategoryName: r.Label,
			//		Score:        r.Score,
			//		FkMainDialog: recommend.Recommend.DialogID,
			//		FkDialog:     dialogs[key].ID,
			//	}
			//}

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

			err = recommend.Delete()
			if err != nil {
				log.Printf("failed to delete message: %v", err)
				continue
			}

			recommendResultSummary, err := rClient.GetRecommendationSummary(recommend.Recommend, util.CategoriesText)
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
		}
	}()

	log.Println("Recommend API is running")
	<-c
	log.Println("Recommend API is shutting down")
}
