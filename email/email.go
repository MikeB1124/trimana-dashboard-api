package email

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/MikeB1124/trimana-dashboard-api/configuration"
	"github.com/MikeB1124/trimana-dashboard-api/payroll"
	"gopkg.in/gomail.v2"
)

func EmailCSVPayrollReport(csvBuffer *bytes.Buffer, filename string, payrollRecords []payroll.EmployeePayrollRecord, startPayPeriod time.Time, endPayPeriod time.Time) error {
	// Send email logic here
	m := gomail.NewMessage()

	subject := fmt.Sprintf("Wilshire Grill Payroll Report %02d/%02d/%d - %02d/%02d/%d",
		startPayPeriod.Month(), startPayPeriod.Day(), startPayPeriod.Year(),
		endPayPeriod.Month(), endPayPeriod.Day(), endPayPeriod.Year(),
	)

	// Set email sender, recipient, and subject
	m.SetHeader("From", configuration.Config.GmailConfig.FromAddress)
	m.SetHeader("To", "trimanaucla@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", buildPayrollRecordBody(payrollRecords))

	// Attach the CSV file from the buffer
	m.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(csvBuffer.Bytes())
		return err
	}))

	// Send the email
	if err := SendEmail(m); err != nil {
		return err
	}
	return nil
}

func PayrollActivityEvent(subject string, body string, employeeEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.GmailConfig.FromAddress)
	to := []string{"trimanaucla@gmail.com"}
	if employeeEmail != "" {
		to = append(to, employeeEmail)
	}
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Send the email
	if err := SendEmail(m); err != nil {
		return err
	}
	return nil

}

func DailyHoursLowEvent(employeeName string, hoursForDay float64, employeeEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.GmailConfig.FromAddress)
	to := []string{"trimanaucla@gmail.com"}
	if employeeEmail != "" {
		to = append(to, employeeEmail)
	}
	subject := fmt.Sprintf("Low Hours for %s < 5.0 Hours", employeeName)
	body := fmt.Sprintf("%s worked %.2f hours today", employeeName, hoursForDay)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	if err := SendEmail(m); err != nil {
		return err
	}
	return nil
}

func DailyHoursHighEvent(employeeName string, hoursForDay float64, employeeEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.GmailConfig.FromAddress)
	to := []string{"trimanaucla@gmail.com"}
	if employeeEmail != "" {
		to = append(to, employeeEmail)
	}
	subject := fmt.Sprintf("High Hours for %s > 8.0 Hours", employeeName)
	body := fmt.Sprintf("%s worked %.2f hours today", employeeName, hoursForDay)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	if err := SendEmail(m); err != nil {
		return err
	}
	return nil
}

func SendEmail(m *gomail.Message) error {
	gmailConfig := configuration.Config.GmailConfig
	// SMTP configuration
	d := gomail.NewDialer(gmailConfig.Host, gmailConfig.Port, gmailConfig.FromAddress, gmailConfig.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func buildPayrollRecordBody(payrollRecords []payroll.EmployeePayrollRecord) string {
	var payrollTotals payroll.PayrollTotals
	header := "Employee Payroll Records\n\n"
	body := ""
	for _, record := range payrollRecords {
		payrollTotals.TotalStandardHours += record.TotalStandardHours
		payrollTotals.TotalDayHours += record.TotalDayHours
		payrollTotals.TotalOverTimeHours += record.TotalOverTimeHours
		payrollTotals.TotalOvertimePay += record.TotalOvertimePay
		payrollTotals.TotalStandardPay += record.TotalStandardPay
		payrollTotals.TotalPay += record.TotalPay
		body += fmt.Sprintf("Name: %s\n", record.EmployeeInfo.Name)
		body += fmt.Sprintf("Employee ID: %s\n", record.EmployeeInfo.EmployeeID)
		body += fmt.Sprintf("Total Standard Hours: %.2f\n", record.TotalStandardHours)
		body += fmt.Sprintf("Total Overtime Hours: %.2f\n", record.TotalOverTimeHours)
		body += fmt.Sprintf("Hours for Today: %.2f\n", record.TotalDayHours)
		body += fmt.Sprintf("Hourly Rate: $%.2f\n", record.EmployeeInfo.HourlyRate)
		body += fmt.Sprintf("Total Standard Pay: $%.2f\n", record.TotalStandardPay)
		body += fmt.Sprintf("Total Overtime Pay: $%.2f\n", record.TotalOvertimePay)
		body += fmt.Sprintf("Total Pay: $%.2f\n", record.TotalPay)
		body += "\n"
	}
	payrollTotals.AverageRate = payrollTotals.TotalPay / (payrollTotals.TotalStandardHours + payrollTotals.TotalOverTimeHours)
	header += "-------------\n"
	header += "Payroll Totals\n"
	header += "-------------\n"
	header += fmt.Sprintf("Total Standard Hours: %.2f\n", payrollTotals.TotalStandardHours)
	header += fmt.Sprintf("Total Overtime Hours: %.2f\n", payrollTotals.TotalOverTimeHours)
	header += fmt.Sprintf("Total Hours for Today: %.2f\n", payrollTotals.TotalDayHours)
	header += fmt.Sprintf("Pay Period Average Hourly Rate: $%.2f\n", payrollTotals.AverageRate)
	header += fmt.Sprintf("Total Standard Pay: $%.2f\n", payrollTotals.TotalStandardPay)
	header += fmt.Sprintf("Total Overtime Pay: $%.2f\n", payrollTotals.TotalOvertimePay)
	header += fmt.Sprintf("Total Pay for Pay Period: $%.2f\n\n\n", payrollTotals.TotalPay)
	return header + body
}
