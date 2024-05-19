package controllers

import "fmt"

func findDateRange(period string) (string, string, error) {
	var startDate, endDate string
	switch period {
	case "DAILY":
		startDate = "2024-05-17T00:00:00Z"
		endDate = "2024-05-17T23:59:59Z"
	case "WEEKLY":
		startDate = "2024-05-17T00:00:00Z"
		endDate = "2024-05-23T23:59:59Z"
	case "MONTHLY":
		startDate = "2024-05-01T00:00:00Z"
		endDate = "2024-05-31T23:59:59Z"
	default:
		return "", "", fmt.Errorf("Invalid period: '%s', Allowed period values (DAILY, WEEKLY, MONTHLY).", period)
	}

	return startDate, endDate, nil
}
