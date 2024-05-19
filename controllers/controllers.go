package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MikeB1124/trimana-dashboard-api/poynt"
	"github.com/aws/aws-lambda-go/events"
)

func GetPoyntTransactions(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %+v\n", event)

	if event.QueryStringParameters["period"] != "" {
		startAt, endAt, err := findDateRange(event.QueryStringParameters["period"])
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       err.Error(),
			}, nil
		}

		transactions, err := poynt.GetPoyntTransactions(startAt, endAt)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		// Marshal the transactions into JSON
		jsonTransactions, err := json.Marshal(transactions)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       err.Error(),
			}, nil
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(jsonTransactions),
		}, nil

	} else {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Processing Slaes Event Has Benn Completed!",
		}, nil
	}
}
