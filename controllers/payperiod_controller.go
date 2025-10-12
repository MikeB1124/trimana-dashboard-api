package controllers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/calendar.go"
	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

type PayPeriod struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Label     string    `json:"label"`
}

func GetPayPeriods(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Get pay periods request")

	// Get biweekly pay schedule
	endDate := time.Now().Year() + 1
	startDate := time.Date(2023, time.January, 2, 0, 0, 0, 0, payroll.TimeZone)

	var payPeriods []PayPeriod
	currentPayDate := startDate.AddDate(0, 0, 13)

	for currentPayDate.Year() <= endDate {
		periodStart := currentPayDate.AddDate(0, 0, -13)
		periodEnd := currentPayDate

		label := periodStart.Format("01/02/2006") + " - " + periodEnd.Format("01/02/2006")

		payPeriods = append(payPeriods, PayPeriod{
			StartDate: periodStart,
			EndDate:   periodEnd,
			Label:     label,
		})

		currentPayDate = currentPayDate.AddDate(0, 0, 14)
	}

	// Reverse to show most recent first
	for i, j := 0, len(payPeriods)-1; i < j; i, j = i+1, j-1 {
		payPeriods[i], payPeriods[j] = payPeriods[j], payPeriods[i]
	}

	responseBody, err := json.Marshal(payPeriods)
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

func GetCurrentPayPeriod(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Get current pay period request")

	dateNow := time.Now().In(payroll.TimeZone)
	endPayPeriod := calendar.GetPayPeriod(dateNow)
	if endPayPeriod.IsZero() {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Could not determine current pay period"})
	}

	startPayPeriod := endPayPeriod.AddDate(0, 0, -13)
	label := startPayPeriod.Format("01/02/2006") + " - " + endPayPeriod.Format("01/02/2006")

	payPeriod := PayPeriod{
		StartDate: startPayPeriod,
		EndDate:   endPayPeriod,
		Label:     label,
	}

	responseBody, err := json.Marshal(payPeriod)
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
