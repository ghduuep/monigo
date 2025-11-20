package notification

import (
	"net/smtp"
)

func SendEmailNotification(url string, status string) error {
	auth := smtp.PlainAuth("", "ghduuep@gmail.com", "krve whaq yzpi jwcq", "smtp.gmail.com")

	to := []string{"ghduep@outlook.com"}
	msg := []byte("To: ghduep@outlook.com\r\n" +
		"Subject: Your server status has changed\r\n" +
		"\r\n" +
		"The status of " + url + " has changed to " + status + ".\r\n")

	return smtp.SendMail("smtp.gmail.com:587", auth, "ghduuep@gmail.com", to, msg)
}
