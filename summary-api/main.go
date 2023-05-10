package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/go-kivik/couchdb/v3"
	"log"
	"main/client"
	"main/db"
	"main/env"
	"main/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const QueueName = "summary-q"
const RecommendQueueName = "recommend-q"

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	if err != nil {
		log.Fatalf("failed to create aws session: %v", err)
	}

	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(QueueName),
	})
	if err != nil {
		log.Fatalf("failed to get %s queue url: %v", QueueName, err)
	}

	recommendUrlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(RecommendQueueName),
	})
	if err != nil {
		log.Fatalf("failed to get %s queue url: %v", RecommendQueueName, err)
	}

	grpcClient, err := client.NewClient()
	if err != nil {
		log.Fatalf("failed to connect to grpc: %v", err)
	}
	defer func(client client.Client) {
		if err := client.Close(); err != nil {
			log.Printf("failed to close grpc: %v", err)
		}
	}(grpcClient)

	couchDB, err := db.NewCouchDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer func(client db.CouchDB) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Close(ctx); err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}(couchDB)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
		loop:
			select {
			case <-c:
				return
			default:
				break
			}

			msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            urlResult.QueueUrl,
				MaxNumberOfMessages: aws.Int64(1),
			})

			if err != nil {
				log.Printf("failed to receive message: %v", err)
				goto loop
			}

			if len(msgResult.Messages) > 0 {
				msg := msgResult.Messages[0]
				log.Printf("message: %v", *msg.Body)

				dialogMap := map[string]string{}
				if err := json.Unmarshal([]byte(*msg.Body), &dialogMap); err != nil {
					log.Printf("failed to unmarshal message: %v", err)
					goto loop
				}

				var dialogID, ok = dialogMap["dialogID"]
				if !ok {
					log.Printf("failed to get dialogID: %v", err)
					goto loop
				}

				text, err := util.GetDialogs(dialogID)
				if err != nil {
					log.Printf("failed to get dialogs: %v", err)
					goto loop
				}

				summary, err := grpcClient.GetSummary(text)
				if err != nil {
					log.Printf("failed to get summary: %v", err)
					goto loop
				}

				err = couchDB.AddSummary(dialogID, summary)
				if err != nil {
					log.Printf("failed to add summary: %v", err)
					goto loop
				}
				log.Printf("add summary to db\n")

				jsonString, err := json.Marshal(map[string]string{
					"dialogID": dialogID,
					"summary":  summary,
				})
				if err != nil {
					log.Printf("failed to marshal recommend content: %v", err)
				}

				_, err = svc.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String(string(jsonString)),
					QueueUrl:    recommendUrlResult.QueueUrl,
				})
				if err != nil {
					log.Printf("failed to send message: %v", err)
					goto loop
				}

				log.Printf("send message to %s queue\n", RecommendQueueName)

				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      urlResult.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					log.Printf("failed to delete message: %v", err)
				}
				log.Printf("delete message from %s queue\n", QueueName)
			}

			time.Sleep(time.Second)
		}
	}()

	log.Println("waiting for messages...")
	<-c
	fmt.Println("exiting...")
}
