package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

func GetTimeCards(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Get timecards request")

	// Check if filtering by employee ID
	employeeID := event.QueryStringParameters["employeeID"]
	var timeCards []payroll.TimeCard
	var err error

	if employeeID != "" {
		log.Printf("Filtering timecards by employee ID: %s", employeeID)
		timeCards, err = payroll.GetTimeCardsByEmployeeID(employeeID)
	} else {
		timeCards, err = payroll.GetAllTimeCards()
	}

	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error getting timecards: " + err.Error()})
	}

	responseBody, err := json.Marshal(timeCards)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error marshalling response"})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func GetTimeCardByID(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Get timecard by ID: %s", id)

	timeCard, err := payroll.GetTimeCardByID(id)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 404, Result: "Timecard not found"})
	}

	responseBody, err := json.Marshal(timeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error marshalling response"})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func CreateTimeCardManual(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Create timecard request")

	var timeCard payroll.TimeCard
	err := json.Unmarshal([]byte(event.Body), &timeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid request body"})
	}

	// Normalize createdAt to midnight in Pacific time (same as card scan)
	createdDate := timeCard.CreatedAt.In(payroll.TimeZone)
	timeCard.CreatedAt = time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, payroll.TimeZone)

	createdTimeCard, err := payroll.CreateTimeCard(timeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error creating timecard: " + err.Error()})
	}

	responseBody, err := json.Marshal(createdTimeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error marshalling response"})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func UpdateTimeCardManual(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Update timecard ID: %s", id)

	var timeCard payroll.TimeCard
	err := json.Unmarshal([]byte(event.Body), &timeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid request body"})
	}

	// Normalize createdAt to midnight in Pacific time (same as card scan)
	createdDate := timeCard.CreatedAt.In(payroll.TimeZone)
	timeCard.CreatedAt = time.Date(createdDate.Year(), createdDate.Month(), createdDate.Day(), 0, 0, 0, 0, payroll.TimeZone)

	updatedTimeCard, err := payroll.UpdateTimeCard(id, timeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error updating timecard: " + err.Error()})
	}

	responseBody, err := json.Marshal(updatedTimeCard)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error marshalling response"})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func DeleteTimeCardManual(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Delete timecard ID: %s", id)

	err := payroll.DeleteTimeCard(id)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error deleting timecard: " + err.Error()})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"result": "Timecard deleted successfully"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
