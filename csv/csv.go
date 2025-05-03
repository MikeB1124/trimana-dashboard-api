package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
)

func WriteCSV(payrollRecords []payroll.EmployeePayrollRecord) (*bytes.Buffer, error) {
	// Write CSV data to a buffer instead of a file
	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	// if err := writer.Write([]string{"##GENERIC## V1.0"}); err != nil {
	// 	return nil, fmt.Errorf("Error writing heaversion headerder: %v", err)
	// }

	// Write header
	// header := []string{
	// 	"IID", "Pay Frequency", "Pay Period Start",
	// 	"Pay Period End", "Employee Id", "Earnings Code",
	// 	"Pay Hours", "Dollars", "Separate Check", "Worked In Dept, Rate Code",
	// }
	// if err := writer.Write(header); err != nil {
	// 	return nil, fmt.Errorf("Error writing header: %v", err)
	// }

	// for _, record := range payrollRecords {
	// 	row := []string{
	// 		record.EmployeeInfo.BranchID,
	// 		record.EmployeeInfo.PayFrequency,
	// 		fmt.Sprintf("%02d/%02d/%d", record.StartDate.Month(), record.StartDate.Day(), record.StartDate.Year()),
	// 		fmt.Sprintf("%02d/%02d/%d", record.EndDate.Month(), record.EndDate.Day(), record.EndDate.Year()),
	// 		record.EmployeeInfo.EmployeeID,
	// 		record.EmployeeInfo.EarningsCode,
	// 		strconv.FormatFloat(record.TotalPayPeriodHours, 'f', 2, 64),
	// 		strconv.FormatFloat(record.Total, 'f', 2, 64),
	// 		"0",
	// 		"",
	// 		record.EmployeeInfo.RateCode,
	// 	}
	// 	if err := writer.Write(row); err != nil {
	// 		return nil, fmt.Errorf("Error writing record: %v", err)
	// 	}
	// }

	// Write header
	header := []string{
		"Employee Name", "Pay Period Start", "Pay Period End", "Pay Rate", "Standard Hours", "Overtime Hours", "Total Standard Pay", "Total Overtime Pay", "Total Pay",
	}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("Error writing header: %v", err)
	}

	for _, record := range payrollRecords {
		row := []string{
			record.EmployeeInfo.Name,
			fmt.Sprintf("%02d/%02d/%d", record.StartDate.Month(), record.StartDate.Day(), record.StartDate.Year()),
			fmt.Sprintf("%02d/%02d/%d", record.EndDate.Month(), record.EndDate.Day(), record.EndDate.Year()),
			strconv.FormatFloat(record.EmployeeInfo.HourlyRate, 'f', 2, 64),
			strconv.FormatFloat(record.TotalStandardHours, 'f', 2, 64),
			strconv.FormatFloat(record.TotalOverTimeHours, 'f', 2, 64),
			strconv.FormatFloat(record.TotalStandardPay, 'f', 2, 64),
			strconv.FormatFloat(record.TotalOvertimePay, 'f', 2, 64),
			strconv.FormatFloat(record.TotalPay, 'f', 2, 64),
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("Error writing record: %v", err)
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("Error flushing writer: %v", err)
	}

	return &buffer, nil
}
