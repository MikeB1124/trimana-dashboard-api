package controllers

import "github.com/aws/aws-lambda-go/events"

func HelloController() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Hello from Lambda!\"",
	}
}

func GoodbyeController() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Goodbye from Lambda!\"",
	}
}

func CustomController() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Custom response from Lambda!\"",
	}
}
