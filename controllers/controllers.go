package controllers

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func GetPoyntSalesByDate(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %+v\n", event)

	eventQueryParam := event.QueryStringParameters
	log.Println("Query parameters: ", eventQueryParam)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Processing Slaes Event Has Benn Completed!",
	}, nil
}
