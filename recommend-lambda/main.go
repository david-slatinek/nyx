package main

import (
	"log"
	"main/email"
	"main/env"
	"main/model"
	"main/queue"
	"os"
	"os/signal"
	"syscall"
)

const QueueEmail = "email-q"

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	emailQueue, err := queue.NewQueue(QueueEmail)
	if err != nil {
		log.Fatalf("failed to get %s queue url: %v", QueueEmail, err)
	}
	recommendChannel := make(chan model.Result, 1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-c:
				return
			default:
				recommend, deleteM, err := emailQueue.Receive()
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
				break
			case recommend := <-recommendChannel:
				log.Printf("recommend: %v", recommend.Recommend)

				err = email.SendEmail(recommend.Recommend)
				if err != nil {
					log.Printf("failed to send email: %v", err)
					continue
				}

				err = recommend.Delete()
				if err != nil {
					log.Printf("failed to delete message: %v", err)
				}
				break
			default:
				break
			}
		}
	}()

	log.Printf("waiting for messages from %s\n", QueueEmail)
	<-c
	log.Println("exiting...")
}
