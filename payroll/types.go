package payroll

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID           string  `json:"id" bson:"_id"`
	BranchID     string  `json:"branchID" bson:"branchID"`
	Name         string  `json:"name" bson:"name"`
	EmployeeID   string  `json:"employeeID" bson:"employeeID"`
	HourlyRate   float64 `json:"hourlyRate" bson:"hourlyRate"`
	PayFrequency string  `json:"payFrequency" bson:"payFrequency"`
	EarningsCode string  `json:"earningsCode" bson:"earningsCode"`
	RateCode     string  `json:"rateCode" bson:"rateCode"`
	CardID       string  `json:"cardID" bson:"cardID"`
	Email        string  `json:"email" bson:"email"`
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
	EmployeeInfo       Employee  `json:"employeeInfo"`
	TotalStandardHours float64   `json:"totalStandardHours"`
	TotalDayHours      float64   `json:"totalDayHours"`
	TotalOverTimeHours float64   `json:"totalOverTimeHours"`
	TotalOvertimePay   float64   `json:"totalOvertimePay"`
	TotalStandardPay   float64   `json:"totalStandardPay"`
	TotalPay           float64   `json:"totalPay"`
	StartDate          time.Time `json:"startDate"`
	EndDate            time.Time `json:"endDate"`
}

type PayrollTotals struct {
	TotalStandardHours float64 `json:"totalStandardHours"`
	TotalDayHours      float64 `json:"totalDayHours"`
	TotalOverTimeHours float64 `json:"totalOverTimeHours"`
	AverageRate        float64 `json:"averageRate"`
	TotalOvertimePay   float64 `json:"totalOvertimePay"`
	TotalStandardPay   float64 `json:"totalStandardPay"`
	TotalPay           float64 `json:"totalPay"`
}
