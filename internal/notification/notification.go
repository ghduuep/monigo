package notification

import (
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmailNotification(userEmail string, url string, status string) error {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	auth := smtp.PlainAuth("", os.Getenv("EMAIL_SENDER"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_SMTP_HOST"))

	to := []string{userEmail}
	msg := []byte("To:" + userEmail + "\r\n" +
		"Subject: Your server status has changed\r\n" +
		"\r\n" +
		"The status of " + url + " has changed to " + status + ".\r\n")

	return smtp.SendMail(os.Getenv("EMAIL_SMTP_HOST")+":"+os.Getenv("EMAIL_SMTP_PORT"), auth, os.Getenv("EMAIL_SENDER"), to, msg)
}
