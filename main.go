package main

import (
	"github.com/MikeB1124/trimana-dashboard-api/controllers"
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("")

	// Existing payroll endpoints
	router.Route("POST", "/payroll/event", controllers.PayrollEvent)
	router.Route("POST", "/payroll/report", controllers.PayrollReport)

	// Authentication endpoints
	router.Route("POST", "/auth/login", controllers.Login)

	// Employee endpoints
	router.Route("GET", "/employees", controllers.GetEmployees)
	router.Route("GET", "/employees/:id", controllers.GetEmployeeByID)
	router.Route("POST", "/employees", controllers.CreateEmployee)
	router.Route("PUT", "/employees/:id", controllers.UpdateEmployee)
	router.Route("DELETE", "/employees/:id", controllers.DeleteEmployee)

	// Timecard endpoints
	router.Route("GET", "/timecards", controllers.GetTimeCards)
	router.Route("GET", "/timecards/:id", controllers.GetTimeCardByID)
	router.Route("POST", "/timecards", controllers.CreateTimeCardManual)
	router.Route("PUT", "/timecards/:id", controllers.UpdateTimeCardManual)
	router.Route("DELETE", "/timecards/:id", controllers.DeleteTimeCardManual)
}

func main() {
	lambda.Start(router.Handler)
}
