package queue

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"main/model"
)

type Queue struct {
	queueName string
	queueUrl  string
	svc       *sqs.SQS
}

func NewQueue(queueName string) (Queue, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	})
	if err != nil {
		return Queue{}, err
	}

	svc := sqs.New(sess)

	urlResult, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return Queue{}, err
	}

	return Queue{
		queueName: queueName,
		queueUrl:  *urlResult.QueueUrl,
		svc:       svc,
	}, nil
}

func (q *Queue) Send(recommendation []model.Recommendation) error {
	message, err := json.Marshal(recommendation)
	if err != nil {
		return err
	}

	_, err = q.svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String(string(message)),
		QueueUrl:    &q.queueUrl,
	})
	return err
}
