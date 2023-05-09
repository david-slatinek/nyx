package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const QueueName = "summary-q"

func main() {
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		for {
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

			if err == nil {
				if len(msgResult.Messages) > 0 {
					msg := msgResult.Messages[0]
					log.Printf("message: %v", *msg.Body)

					_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
						QueueUrl:      urlResult.QueueUrl,
						ReceiptHandle: msg.ReceiptHandle,
					})

					if err != nil {
						log.Printf("failed to delete message: %v", err)
					}
				}
			} else {
				log.Printf("failed to receive message: %v", err)
			}

			time.Sleep(time.Second)
		}
	}()

	log.Println("waiting for messages...")
	<-c
	fmt.Println("exiting...")
}
