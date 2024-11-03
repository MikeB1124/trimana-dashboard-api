package controllers

import (
	"encoding/json"

	"github.com/MikeB1124/trimana-dashboard-api/email"
	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

func createResponse(response payroll.Response) (events.APIGatewayProxyResponse, error) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		responseBody, _ = json.Marshal(payroll.Response{Result: err.Error()})
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(responseBody),
		}, nil
	}

	if response.StatusCode > 200 {
		if err := email.PayrollActivityEvent("Trimana Payroll Error", response.Result, ""); err != nil {
			responseBody, _ = json.Marshal(payroll.Response{Result: err.Error()})
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       string(responseBody),
			}, nil
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseBody),
	}, nil
}
