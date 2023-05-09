package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"main/client"
	"main/env"
	"main/util"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const QueueName = "summary-q"

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "david",
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
		log.Fatalf("failed to get queue url: %v", err)
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

				if _, ok := dialogMap["dialogID"]; !ok {
					log.Printf("failed to get dialogID: %v", err)
					goto loop
				}

				text, err := util.GetDialogs(dialogMap["dialogID"])
				if err != nil {
					log.Printf("failed to get dialogs: %v", err)
					goto loop
				}

				summary, err := grpcClient.GetSummary(text)
				if err != nil {
					log.Printf("failed to get summary: %v", err)
					goto loop
				}
				log.Printf("summary: %s", summary)

				_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      urlResult.QueueUrl,
					ReceiptHandle: msg.ReceiptHandle,
				})

				if err != nil {
					log.Printf("failed to delete message: %v", err)
				}
			}

			time.Sleep(time.Second)
		}
	}()

	log.Println("waiting for messages...")
	<-c
	fmt.Println("exiting...")
}
