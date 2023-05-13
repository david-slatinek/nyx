package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"main/model"
	pb "main/schema"
	"main/util"
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

func (receiver Client) GetRecommendation(summary string, dialogs []model.Dialog, categories []model.Category) ([]model.RecommendResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dialogsMap = make(map[string]bool, len(dialogs))
	for _, dialog := range dialogs {
		dialogsMap[dialog.Text] = true
	}

	var dialogsText = make([]string, 0, len(dialogs))
	for key := range dialogsMap {
		dialogsText = append(dialogsText, key)
	}

	var categoriesText = make([]string, 0, len(categories))
	util.GetCategoriesNames(categories, &categoriesText)

	recommendResponse, err := receiver.client.Recommend(ctx, &pb.RecommendRequest{
		Dialogs:    dialogsText,
		Summary:    summary,
		Categories: categoriesText,
	})
	if err != nil {
		return nil, err
	}

	var recommendResult []model.RecommendResult
	for key, value := range recommendResponse.Labels {
		recommendResult = append(recommendResult, model.RecommendResult{
			Label: value,
			Score: recommendResponse.Scores[key],
		})
	}
	return recommendResult, nil
}
