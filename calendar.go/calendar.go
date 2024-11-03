package calendar

import (
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
)

func getBiWeeklyPaySchedule() []time.Time {
	endDate := time.Now().Year() + 1
	startDate := time.Date(2023, time.January, 2, 0, 0, 0, 0, payroll.TimeZone)

	var paySchedule []time.Time
	currentPayDate := startDate.AddDate(0, 0, 13)

	for currentPayDate.Year() <= endDate {
		paySchedule = append(paySchedule, currentPayDate)
		currentPayDate = currentPayDate.AddDate(0, 0, 14) // Add 14 days for biweekly pay
	}
	return paySchedule
}

func GetPayPeriod(currentDateTime time.Time) time.Time {
	paySchedule := getBiWeeklyPaySchedule()
	for i, payDate := range paySchedule {
		if currentDateTime.Year() == payDate.Year() && currentDateTime.Month() == payDate.Month() && currentDateTime.Day() == payDate.Day() {
			return payDate
		}
		if currentDateTime.After(payDate) && currentDateTime.Before(paySchedule[i+1]) {
			return paySchedule[i+1]
		}
	}
	return time.Time{}
}
