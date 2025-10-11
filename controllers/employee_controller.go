package controllers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"github.com/aws/aws-lambda-go/events"
)

func GetEmployees(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Get employees request")

	employees, err := payroll.GetEmployees()
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error getting employees: " + err.Error()})
	}

	responseBody, err := json.Marshal(employees)
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

func GetEmployeeByID(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Get employee by ID: %s", id)

	employee, err := payroll.GetEmployeeByID(id)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 404, Result: "Employee not found"})
	}

	responseBody, err := json.Marshal(employee)
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

func CreateEmployee(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Create employee request")

	var employee payroll.Employee
	err := json.Unmarshal([]byte(event.Body), &employee)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid request body"})
	}

	createdEmployee, err := payroll.CreateEmployee(employee)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error creating employee: " + err.Error()})
	}

	responseBody, err := json.Marshal(createdEmployee)
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

func UpdateEmployee(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Update employee ID: %s", id)

	var employee payroll.Employee
	err := json.Unmarshal([]byte(event.Body), &employee)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 400, Result: "Invalid request body"})
	}

	updatedEmployee, err := payroll.UpdateEmployee(id, employee)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error updating employee: " + err.Error()})
	}

	responseBody, err := json.Marshal(updatedEmployee)
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

func DeleteEmployee(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := event.PathParameters["id"]
	log.Printf("Delete employee ID: %s", id)

	err := payroll.DeleteEmployee(id)
	if err != nil {
		return createResponse(payroll.Response{StatusCode: 500, Result: "Error deleting employee: " + err.Error()})
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"result": "Employee deleted successfully"}`,
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}
