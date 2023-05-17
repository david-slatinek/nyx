package queue

import (
	"encoding/json"
	"errors"
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

func (q Queue) Send(recommendation []map[string]string) error {
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

func (q Queue) Receive() (model.Recommend, func() error, error) {
	result, err := q.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &q.queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		return model.Recommend{}, nil, err
	}

	if len(result.Messages) == 0 {
		return model.Recommend{}, nil, errors.New("no messages in queue")
	}

	var recommendation model.Recommend
	err = json.Unmarshal([]byte(*result.Messages[0].Body), &recommendation)
	if err != nil {
		return model.Recommend{}, nil, err
	}

	deleteMessage := func() error {
		_, err = q.svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &q.queueUrl,
			ReceiptHandle: result.Messages[0].ReceiptHandle,
		})
		return err
	}

	return recommendation, deleteMessage, nil
}
