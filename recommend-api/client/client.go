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
	client     pb.RecommendServiceClient
}

func NewClient() (Client, error) {
	conn, err := grpc.Dial(os.Getenv("RECOMMEND_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return Client{}, err
	}
	return Client{
		connection: conn,
		client:     pb.NewRecommendServiceClient(conn),
	}, nil
}

func (receiver Client) Close() error {
	return receiver.connection.Close()
}

func (receiver Client) GetRecommendationDialogs(dialogs []model.Dialog, categoriesText []string) ([]model.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dialogsMap = make(map[string]string, len(dialogs))
	for _, dialog := range dialogs {
		dialogsMap[dialog.Text] = dialog.ID
	}

	var dialogsText = make([]string, 0, len(dialogs))
	for key := range dialogsMap {
		dialogsText = append(dialogsText, key)
	}

	recommendResponse, err := receiver.client.RecommendDialog(ctx, &pb.RecommendRequestDialog{
		Dialogs:    dialogsText,
		Categories: categoriesText,
	})
	if err != nil {
		return nil, err
	}

	var recommendResult []model.RecommendResult
	for _, value := range recommendResponse.Responses {
		recommendResult = append(recommendResult, model.RecommendResult{
			ID:     dialogsMap[value.Text],
			Dialog: value.Text,
			Labels: value.Labels,
			Scores: value.Scores,
		})
	}
	return recommendResult, nil
}

func (receiver Client) GetRecommendationSummary(summary model.Recommend, categoriesText []string) (model.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	recommendResponse, err := receiver.client.RecommendSummary(ctx, &pb.RecommendRequestSummary{
		Summary:    summary.Summary,
		Categories: categoriesText,
	})
	if err != nil {
		return model.RecommendResult{}, err
	}

	return model.RecommendResult{
		ID:     summary.DialogID,
		Dialog: recommendResponse.Text,
		Labels: recommendResponse.Labels,
		Scores: recommendResponse.Scores,
	}, nil
}
