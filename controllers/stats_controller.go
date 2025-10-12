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

type StatsRequest struct {
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
}

func GetStats(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Get stats request")

	var startDate, endDate time.Time
	dateNow := time.Now().In(payroll.TimeZone)

	// Check if custom date range is provided
	startDateParam := event.QueryStringParameters["startDate"]
	endDateParam := event.QueryStringParameters["endDate"]

	if startDateParam != "" && endDateParam != "" {
		var err error
		startDate, err = time.Parse(time.RFC3339, startDateParam)
		if err != nil {
			return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid start date format"})
		}
		endDate, err = time.Parse(time.RFC3339, endDateParam)
		if err != nil {
			return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid end date format"})
		}
		startDate = startDate.In(payroll.TimeZone)
		endDate = endDate.In(payroll.TimeZone)
	} else {
		// Use current pay period
		endDate = calendar.GetPayPeriod(dateNow)
		if endDate.IsZero() {
			return createResponse(payroll.Response{StatusCode: 500, Result: "Could not determine current pay period"})
		}
		startDate = endDate.AddDate(0, 0, -13)
	}

	log.Printf("Getting stats for period: %s - %s", startDate, endDate)

	stats, err := payroll.GetPayrollStats(startDate, endDate)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error calculating stats: " + err.Error()})
	}

	// Get today's stats
	todayTimecards, todayHours, _ := payroll.GetTodayStats()
	
	// Add today's stats to response
	response := map[string]interface{}{
		"stats":          stats,
		"todayTimecards": todayTimecards,
		"todayHours":     todayHours,
		"startDate":      startDate,
		"endDate":        endDate,
	}

	responseBody, err := json.Marshal(response)
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

