package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/calendar.go"
	"github.com/MikeB1124/trimana-dashboard-api/csv"
	"github.com/MikeB1124/trimana-dashboard-api/email"
	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

type PayrollEventRequest struct {
	CardID string `json:"cardID"`
}

func PayrollEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Payroll event %+v\n", event)
	dateNow := time.Now().In(payroll.TimeZone)
	dateNow = time.Date(dateNow.Year(), dateNow.Month(), dateNow.Day(), 0, 0, 0, 0, payroll.TimeZone)

	// Unmarshal the request
	request := PayrollEventRequest{}
	err := json.Unmarshal([]byte(event.Body), &request)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Error unmarshalling request: " + err.Error()})
	}
	log.Printf("Payroll event request: %+v\n", request)

	//Get employee by card ID scanned
	employee, err := payroll.GetEmployeeByCardID(request.CardID)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: fmt.Sprintf("Error getting employee by card ID (%s): %s", request.CardID, err.Error())})
	}
	log.Printf("Employee Found: %+v\n", employee)

	// Get time card for today
	timeCardToday, err := payroll.GetTimeCardByDate(employee.EmployeeID, dateNow)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error getting time card: " + err.Error()})
	}

	// Create a new time card if one does not exist for today
	if timeCardToday == nil {
		if err := payroll.CreateNewTimeCard(employee.EmployeeID, dateNow); err != nil {
			log.Printf("Error creating time card: %v\n", err)
			return createResponse(payroll.Response{StatusCode: 500, Result: "Error creating time card: " + err.Error()})
		}
		log.Printf("New Time Card Created for %s at %s\n", employee.Name, dateNow)
		if err := email.PayrollActivityEvent(fmt.Sprintf("%s Has Checked In", employee.Name), "New Time Card Created For the Day", employee.Email); err != nil {
			log.Printf("Error sending email: %v\n", err)
			return createResponse(payroll.Response{StatusCode: 500, Result: "Error sending email: " + err.Error()})
		}
		return createResponse(payroll.Response{StatusCode: 200, Result: fmt.Sprintf("New Time Card Created for %s at %s", employee.Name, dateNow)})
	} else {
		result, err := payroll.StampTimeCard(timeCardToday)
		if err != nil {
			log.Printf("Error stamping time card: %v\n", err)
			return createResponse(payroll.Response{StatusCode: 500, Result: "Error stamping time card: " + err.Error()})
		}
		activity := fmt.Sprintf("%s Has %s", employee.Name, result)
		log.Println(activity)
		if err := email.PayrollActivityEvent(activity, "", employee.Email); err != nil {
			log.Printf("Error sending email: %v\n", err)
			return createResponse(payroll.Response{StatusCode: 500, Result: "Error sending email: " + err.Error()})
		}
		return createResponse(payroll.Response{StatusCode: 200, Result: fmt.Sprintf("%s has %s", employee.Name, result)})
	}
}

func PayrollReport(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Creating payroll report\n")
	dateNow := time.Now().In(payroll.TimeZone)
	endPayPeriod := calendar.GetPayPeriod(dateNow)
	if endPayPeriod.IsZero() {
		log.Printf("Could not find pay date for %s\n", dateNow)
		return createResponse(payroll.Response{StatusCode: 500, Result: fmt.Sprintf("Could not find pay date for %s", dateNow)})
	}
	startPayPeriod := endPayPeriod.AddDate(0, 0, -13)
	log.Printf("Pay Period: %s - %s\n", startPayPeriod, endPayPeriod)

	employees, err := payroll.GetEmployees()
	if err != nil {
		log.Printf("Error getting employees: %v\n", err)
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error getting employees: " + err.Error()})
	}
	log.Printf("Employees: %+v\n", employees)

	var employeePayrollsRecords []payroll.EmployeePayrollRecord
	for _, employee := range employees {
		log.Printf("Processing payroll for %s\n", employee.Name)
		// Get time cards for the current pay period
		timeCards, err := payroll.GetTimeCardsByPayPeriod(employee.EmployeeID, startPayPeriod, endPayPeriod)
		if err != nil {
			log.Printf("Error getting time cards: %v\n", err)
			return createResponse(payroll.Response{StatusCode: 500, Result: "Error getting time cards: " + err.Error()})
		}
		log.Printf("Time Cards for %s: %+v\n", employee.Name, timeCards)

		// Calculate total hours and total pay
		totalHours := 0.0
		for _, timeCard := range timeCards {
			for _, block := range timeCard.TimeBlocks {
				if !block.CheckOut.IsZero() {
					totalHours += block.CheckOut.Sub(block.CheckIn).Hours()
				}
			}
		}

		if totalHours == 0.0 {
			log.Printf("%s did not work any hours for pay period %s - %s\n", employee.Name, startPayPeriod, endPayPeriod)
			continue
		}

		totalPay := math.Round((totalHours*employee.HourlyRate)*100) / 100

		employeePayrollsRecords = append(employeePayrollsRecords, payroll.EmployeePayrollRecord{
			Name:       employee.Name,
			EmployeeID: employee.EmployeeID,
			Hours:      math.Round(totalHours*100) / 100,
			HourlyRate: employee.HourlyRate,
			Total:      totalPay,
			StartDate:  startPayPeriod,
			EndDate:    endPayPeriod,
		})
		log.Printf("Payroll for %s: Hours: %.2f, Total: %.2f\n", employee.Name, totalHours, totalPay)
	}

	if len(employeePayrollsRecords) == 0 {
		log.Printf("No payroll records to process\n")
		return createResponse(payroll.Response{StatusCode: 200, Result: "No payroll records to process"})
	}

	// Write payroll records to CSV buffer
	csvBuffer, err := csv.WriteCSV(employeePayrollsRecords)
	if err != nil {
		log.Printf("Error writing CSV: %v\n", err)
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error writing CSV: " + err.Error()})
	}
	log.Printf("CSV buffer created\n")

	csvFileName := fmt.Sprintf("Biweekly-%02d%02d%d-%02d%02d%d.csv",
		startPayPeriod.Month(), startPayPeriod.Day(), startPayPeriod.Year(),
		endPayPeriod.Month(), endPayPeriod.Day(), endPayPeriod.Year(),
	)

	if err := email.EmailCSVPayrollReport(csvBuffer, csvFileName, employeePayrollsRecords, startPayPeriod, endPayPeriod); err != nil {
		log.Printf("Error sending email: %v\n", err)
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error sending email: " + err.Error()})
	}
	log.Printf("Email sent\n")

	return createResponse(payroll.Response{StatusCode: 200, Result: "OK"})
}
