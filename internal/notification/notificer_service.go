package notification

import (
	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	Client *resend.Client
	Sender string
}

func NewEmailService(apiKey, sender string) *EmailService {
	client := resend.NewClient(apiKey)

	return &EmailService{
		Client: client,
		Sender: sender,
	}
}

func (s *EmailService) SendEmail(to, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    s.Sender,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := s.Client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
