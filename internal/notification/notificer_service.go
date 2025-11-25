package notification

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	Host string
	Port string
	Sender string
	Password string
}

func NewEmailService(host, port, sender, password string) *EmailService {
	return &EmailService{
		Host: host,
		Port: port,
		Sender: sender,
		Password: password,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.Sender, s.Password, s.Host)
	address := s.Host + ":" + s.Port

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	var msg []byte
	msg = fmt.Appendf(msg, "To: %s\r\nSubject: %s\r\n%s\r\n%s", to, subject, mime, body)

	return smtp.SendMail(address, auth, s.Sender, []string{to}, msg)
}
