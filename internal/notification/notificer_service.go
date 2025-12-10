package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/resend/resend-go/v3"
)

type Notifier interface {
	Send(to, subject, body string) error
}

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

func (s *EmailService) Send(to, subject, body string) error {
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

type TelegramService struct {
	BotToken string
}

func NewTelegramService(botToken string) *TelegramService {
	return &TelegramService{BotToken: botToken}
}

func (t *TelegramService) Send(to, subject, body string) error {
	msg := subject + "\n" + body
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := map[string]string{
		"chat_id": to,
		"text":    msg,
	}

	data, _ := json.Marshal(payload)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	return err
}
