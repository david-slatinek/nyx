package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"main/db"
	"main/model"
	pb "main/schema"
	"os"
	"time"
)

type Client struct {
	connection *grpc.ClientConn
	client     pb.DialogServiceClient
	db         db.CouchDB
}

func NewClient(db db.CouchDB) (Client, error) {
	conn, err := grpc.Dial(os.Getenv("DIALOG_URL")+":9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Client{}, nil
	}
	return Client{
		connection: conn,
		client:     pb.NewDialogServiceClient(conn),
		db:         db,
	}, nil
}

func (receiver Client) Close() error {
	return receiver.connection.Close()
}

func (receiver Client) GetAnswer(dialog model.Dialog) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	answer, err := receiver.client.Dialog(ctx, &pb.DialogRequest{Text: dialog.Text})
	return answer.Answer, err
}
