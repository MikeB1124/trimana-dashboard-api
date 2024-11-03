package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/MikeB1124/trimana-dashboard-api/payroll"
)

func WriteCSV(fileName string, payrollRecords []payroll.EmployeePayrollRecord) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating file: %v", err)
	}
	defer file.Close()

	// Initialize CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Branch Code", "Pay Frequency", "Pay Period Start Date",
		"Pay Period End Date", "Employee Id", "Earnings",
		"Hours", "Dollars", "Separate", "Rate Code",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("Error writing header: %v", err)
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
			return fmt.Errorf("Error writing record: %v", err)
		}
	}

	return nil
}
