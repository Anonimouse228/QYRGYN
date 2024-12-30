package util

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
)

var (
	smtpHost = "smtp.gmail.com"
	smtpPort = 587
	smtpUser = "suhansun13@gmail.com"
	smtpPass = "rimc awza ulyx gwfc"
)

// SendEmail sends an email with optional attachments
func SendEmail(to, subject, body string, attachments []string) error {
	fmt.Println("SMTP_HOST:", smtpHost)
	fmt.Println("SMTP_USER:", smtpUser)

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Attach files if provided
	for _, file := range attachments {
		m.Attach(file)
	}

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Email send error:", err)
		return err
	}
	return nil
}
