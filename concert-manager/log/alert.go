package log

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailAlerter struct {
	Emailer EmailService
	To      string
}

const alertToEmailEnvVar = "CM_ALERT_EMAIL"
const gmailUserEnv = "CM_GMAIL_USER"
const gmailPassEnv = "CM_GMAIL_PASSWORD"

func NewGmailAlerter() (*EmailAlerter, error) {
	to := os.Getenv(alertToEmailEnvVar)
	if to == "" {
		return nil, fmt.Errorf("%s environment variable must be set", alertToEmailEnvVar)
	}
	user := os.Getenv(gmailUserEnv)
	if user == "" {
		return nil, fmt.Errorf("%s environment variable must be set", gmailUserEnv)
	}
	pass := os.Getenv(gmailPassEnv)
	if pass == "" {
		return nil, fmt.Errorf("%s environment variable must be set", gmailPassEnv)
	}

	emailer := EmailService{
		smtpServer:   "smtp.gmail.com",
		smtpPort:     587,
		smtpUsername: user,
		smtpPassword: pass,
		fromEmail:    "cm-alert@gmail.com",
	}
	return &EmailAlerter{emailer, to}, nil
}

func (a EmailAlerter) Alert(body string) error {
	return a.Emailer.Send(a.To, "ALERT - Concert Manager", body)
}

type EmailService struct {
	smtpServer   string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

func (s *EmailService) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.smtpServer, s.smtpPort)
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s", s.fromEmail, to, subject, body))
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpServer)

	return smtp.SendMail(addr, auth, s.fromEmail, []string{to}, message)
}
