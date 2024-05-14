package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func HelloController(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
}

func GoodbyeController(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Goodbye from Lambda!\"",
	}
}

func CustomController(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Custom response from Lambda!\"",
	}
}

type Name struct {
	Name string `json:"name"`
}

func CustomPostController(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	var personName Name
	err := json.Unmarshal([]byte(event.Body), &personName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("\"Error parsing request: %s\"", err),
		}
	}

	personNameJson, err := json.Marshal(personName)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("\"Error marshalling response: %s\"", err),
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(personNameJson),
	}
}
