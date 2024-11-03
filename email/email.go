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

	subject := fmt.Sprintf("Trimana Payroll Report %02d/%02d/%d - %02d/%02d/%d",
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

func PayrollActivityEvent(subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", configuration.Config.GmailConfig.FromAddress)
	m.SetHeader("To", "trimanaucla@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Send the email
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
	body := "Employee Payroll Records\n\n"
	for _, record := range payrollRecords {
		body += fmt.Sprintf("Name: %s\n", record.Name)
		body += fmt.Sprintf("Employee ID: %s\n", record.EmployeeID)
		body += fmt.Sprintf("Hours: %.2f\n", record.Hours)
		body += fmt.Sprintf("Hourly Rate: $%.2f\n", record.HourlyRate)
		body += fmt.Sprintf("Total: $%.2f\n", record.Total)
		body += "\n"
	}
	return body
}
