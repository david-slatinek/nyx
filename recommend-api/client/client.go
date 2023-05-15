package client

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"main/model"
	pb "main/schema"
	"os"
	"sort"
	"time"
)

type Client struct {
	connection *grpc.ClientConn
	client     pb.RecommendServiceClient
}

func NewClient() (Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, os.Getenv("RECOMMEND_MODEL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (receiver Client) GetRecommendationDialogs(dialogs []model.Dialog, categoriesText []string) ([]model.Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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

	var recommendResult []model.Recommendation
	for _, value := range recommendResponse.Responses {
		for i := range value.Scores {
			if value.Scores[i] > 0.1 {
				recommendResult = append(recommendResult, model.Recommendation{
					ID:    dialogsMap[value.Text],
					Label: value.Labels[i],
					Score: value.Scores[i],
				})
			}
		}
	}

	sort.Slice(recommendResult, func(i, j int) bool {
		return recommendResult[i].Score < recommendResult[j].Score
	})

	var mapRecommendResult = make(map[string]model.Recommendation, len(recommendResult))
	for _, value := range recommendResult {
		mapRecommendResult[value.Label] = value
	}

	recommendResult = make([]model.Recommendation, 0, len(mapRecommendResult))
	for _, value := range mapRecommendResult {
		recommendResult = append(recommendResult, value)
	}

	sort.Slice(recommendResult, func(i, j int) bool {
		return recommendResult[i].Score > recommendResult[j].Score
	})
	return recommendResult, nil
}

func (receiver Client) GetRecommendationSummary(summary model.Recommend, categoriesText []string) ([]model.Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	recommendResponse, err := receiver.client.RecommendSummary(ctx, &pb.RecommendRequestSummary{
		Summary:    summary.Summary,
		Categories: categoriesText,
	})
	if err != nil {
		return nil, err
	}

	var recommendResult []model.Recommendation
	for i := range recommendResponse.Scores {
		if recommendResponse.Scores[i] > 0.1 {
			recommendResult = append(recommendResult, model.Recommendation{
				ID:    summary.DialogID,
				Label: recommendResponse.Labels[i],
				Score: recommendResponse.Scores[i],
			})
		}
	}

	sort.Slice(recommendResult, func(i, j int) bool {
		return recommendResult[i].Score > recommendResult[j].Score
	})
	return recommendResult, nil
}
