package main

import (
	"log"
	"net/smtp"
)

func main() {
	sendEmailNotification("https://example.com", "UP")
}

func sendEmailNotification(url string, status string) {
	auth := smtp.PlainAuth("", "ghduuep@gmail.com", "krve whaq yzpi jwcq", "smtp.gmail.com")

	to := []string{"ghduep@outlook.com"}
	msg := []byte("To: ghduep@outlook.com\r\n" +
		"Subject: Your server status has changed\r\n" +
		"\r\n" +
		"The status of " + url + " has changed to " + status + ".\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, "ghduep@outlook.com", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
