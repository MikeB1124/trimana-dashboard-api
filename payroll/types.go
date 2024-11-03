package payroll

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID         string  `json:"id" bson:"_id"`
	Name       string  `json:"name" bson:"name"`
	EmployeeID string  `json:"employeeID" bson:"employeeID"`
	HourlyRate float64 `json:"hourlyRate" bson:"hourlyRate"`
	CardID     string  `json:"cardID" bson:"cardID"`
	Email      string  `json:"email" bson:"email"`
}

type TimeCard struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	EmployeeID string             `json:"employeeID" bson:"employeeID"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	TimeBlocks []TimeBlock        `json:"timeBlocks" bson:"timeBlocks"`
}

type TimeBlock struct {
	CheckIn  time.Time `json:"checkIn" bson:"checkIn"`
	CheckOut time.Time `json:"checkOut,omitempty" bson:"checkOut,omitempty"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Result     string `json:"result"`
}

type EmployeePayrollRecord struct {
	Name       string    `json:"name"`
	EmployeeID string    `json:"employeeID"`
	Hours      float64   `json:"hours"`
	Total      float64   `json:"total"`
	HourlyRate float64   `json:"hourlyRate"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
}
