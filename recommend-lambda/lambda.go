package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"main/email"
	"main/env"
	"main/model"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		var recommendation []model.Recommend
		err := json.Unmarshal([]byte(message.Body), &recommendation)
		if err != nil {
			log.Printf("failed to unmarshal message body: %v", err)
			continue
		}

		err = email.SendEmail(recommendation)
		if err != nil {
			log.Printf("failed to send email: %v", err)
			continue
		}
	}
	return nil
}

func main() {
	err := env.Load("env/.env")
	if err != nil {
		log.Fatalf("failed to load env: %v", err)
	}
	lambda.Start(handler)
}
