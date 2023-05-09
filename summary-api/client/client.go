package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "main/schema"
	"os"
	"time"
)

type Client struct {
	connection *grpc.ClientConn
	client     pb.SummaryServiceClient
}

func NewClient() (Client, error) {
	conn, err := grpc.Dial(os.Getenv("SUMMARY_URL")+":9050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Client{}, err
	}
	return Client{
		connection: conn,
		client:     pb.NewSummaryServiceClient(conn),
	}, nil
}

func (receiver Client) Close() error {
	return receiver.connection.Close()
}

func (receiver Client) GetSummary(text string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	answer, err := receiver.client.Summary(ctx, &pb.SummaryRequest{
		Text: text,
	})
	if err != nil {
		return "", err
	}
	return answer.Summary, nil
}
