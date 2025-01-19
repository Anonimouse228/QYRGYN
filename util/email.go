package util

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
)

var (
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = 587
	smtpUser = os.Getenv("SMTP_USER")
	smtpPass = os.Getenv("SMTP_PASS")
)

// SendEmail sends an email with optional attachments
func SendEmail(to, subject, body string, attachments []string) error {
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

func SendVerificationEmail(email, token string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", "your-email@example.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Email Verification")
	link := fmt.Sprintf("https://qyrgyn.onrender.com/verify?token=%s", token) // Replace with your domain
	mail.SetBody("text/html", fmt.Sprintf(`<h1>Verify Your Email</h1>
    <p>Click the link below to verify your email address:</p>
    <a href="%s">%s</a>`, link, link))

	// SMTP Config
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	return dialer.DialAndSend(mail)
}
