package main

import (
	"context"
	"log"

	"github.com/MikeB1124/trimana-dashboard-api/controllers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%+v", event)

	var response events.APIGatewayProxyResponse
	switch event.Path {
	case "/hello":
		response = controllers.HelloController()
	case "/goodbye":
		response = controllers.GoodbyeController()
	case "/custom":
		response = controllers.CustomController()
	default:
		response = events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "\"Not found\"",
		}
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
