package controllers

import (
	"context"
	"fmt"
	"log"

	"github.com/MikeB1124/trimana-dashboard-api/poynt"
	"github.com/aws/aws-lambda-go/events"
)

func GetPoyntSalesByDate(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %+v\n", event)

	if event.QueryStringParameters["period"] != "" {
		if event.QueryStringParameters["period"] == "DAILY" {
			log.Println("Processing Daily Sales Event...")
			transactions, err := poynt.GetPoyntTransactions("2024-05-17T00:00:00Z", "2024-05-17T23:59:59Z")
			if err != nil {
				return events.APIGatewayProxyResponse{
					StatusCode: 500,
					Body:       fmt.Sprintf("Error fetching transactions: %v", err),
				}, nil
			}
			log.Printf("Transactions: %+v\n", transactions)
		} else if event.QueryStringParameters["period"] == "WEEKLY" {
			log.Println("Processing Weekly Sales Event...")
		} else if event.QueryStringParameters["period"] == "MONTHLY" {
			log.Println("Processing Monthly Sales Event...")
		} else {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       fmt.Sprintf("Invalid period: '%s', Allowed period values (DAILY, WEEKLY, MONTHLY).", event.QueryStringParameters["period"]),
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Processing Slaes Event Has Benn Completed!",
	}, nil
}
