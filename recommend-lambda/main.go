package main

import (
	"html/template"
	"log"
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

	f, err := os.Create("output.html")
	if err != nil {
		log.Printf("failed to create file: %v", err)
	}

	go func() {
		for {
			select {
			case <-c:
				break
			case recommend := <-recommendChannel:
				log.Printf("recommend: %v", recommend.Recommend)

				t, err := template.New("template.html").ParseFiles("template.html")
				if err != nil {
					log.Printf("failed to parse template: %v", err)
					break
				}

				err = t.Execute(f, recommend.Recommend)
				if err != nil {
					log.Printf("failed to execute template: %v", err)
				}
				break
			default:
				break
			}
		}
	}()

	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(f)

	log.Printf("waiting for messages from %s\n", QueueEmail)
	<-c
	log.Println("exiting...")
}
