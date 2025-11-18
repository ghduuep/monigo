package main

import (
	"net/smtp"
)

func sendEmailNotification(website Website, status string){
	auth := smtp.PlainAuth()

				to := []string{"ghduep@outlook.com"}
				msg := []byte("To: ghduep@outlook.com\r\n" +
					"Subject: Your server status has changed\r\n" +
					"\r\n" +
					"The status of " + website.URL + " has changed to " + status + ".\r\n")
			}

			err := smtp.SendMail("smtp.dominio:587", auth, from, to, msg)
			if err != nil {
				log.Fatal(err)
			}
}
