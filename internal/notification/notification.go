package notification

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmailNotification(userEmail, domain, subject, message string) error {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_SENDER"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_SMTP_SERVER"))

	to := []string{userEmail}
	msg := []byte("To:" + userEmail + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" +
		message + ".\r\n")

	return smtp.SendMail(os.Getenv("EMAIL_SMTP_SERVER")+":"+os.Getenv("EMAIL_SMTP_PORT"), auth, os.Getenv("EMAIL_SENDER"), to, msg)
}
