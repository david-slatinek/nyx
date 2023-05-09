package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"main/model"
	pb "main/schema"
	"os"
	"time"
)

type Client struct {
	connection *grpc.ClientConn
	client     pb.DialogServiceClient
}

func NewClient() (Client, error) {
	conn, err := grpc.Dial(os.Getenv("DIALOG_URL")+":9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Client{}, err
	}
	return Client{
		connection: conn,
		client:     pb.NewDialogServiceClient(conn),
	}, nil
}

func (receiver Client) Close() error {
	return receiver.connection.Close()
}

func (receiver Client) GetAnswer(dialog model.Dialog) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	answer, err := receiver.client.Dialog(ctx, &pb.DialogRequest{
		Text: dialog.Text,
	})
	if err != nil {
		return "", err
	}
	return answer.Answer, nil
}
