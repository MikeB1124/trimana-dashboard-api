package main

import (
	"context"
	"log"
	"net/http"

	"github.com/MikeB1124/trimana-dashboard-api/controllers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("%+v", event)

	var response events.APIGatewayProxyResponse
	switch event.Path {
	case "/hello":
		if event.HTTPMethod == "GET" {
			response = controllers.HelloController(event)
		}
	case "/goodbye":
		if event.HTTPMethod == "GET" {
			response = controllers.GoodbyeController(event)
		}
	case "/custom":
		if event.HTTPMethod == "GET" {
			response = controllers.CustomController(event)
		}
	default:
		response = events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "\"Not found\"",
		}
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
