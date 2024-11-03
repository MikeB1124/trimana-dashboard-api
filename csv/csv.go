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

	// Write header
	header := []string{
		"Branch Code", "Pay Frequency", "Pay Period Start Date",
		"Pay Period End Date", "Employee Id", "Earnings",
		"Hours", "Dollars", "Separate", "Rate Code",
	}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("Error writing header: %v", err)
	}

	for _, record := range payrollRecords {
		row := []string{
			"ABCDEFG",
			"B",
			fmt.Sprintf("%02d/%02d/%d", record.StartDate.Month(), record.StartDate.Day(), record.StartDate.Year()),
			fmt.Sprintf("%02d/%02d/%d", record.EndDate.Month(), record.EndDate.Day(), record.EndDate.Year()),
			record.EmployeeID,
			"REG",
			strconv.FormatFloat(record.Hours, 'f', 2, 64),
			strconv.FormatFloat(record.Total, 'f', 2, 64),
			"0",
			"BASE",
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
