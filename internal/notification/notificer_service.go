package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification/templates"
	"github.com/resend/resend-go/v3"
	"net/http"
	"time"
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

func (s *EmailService) SendStatusAlert(to string, m models.Monitor, result models.CheckResult, duration time.Duration) error {
	var subject, body string

	if m.Type == models.TypeHTTP {
		subject, body = templates.BuildEmailHTTPMessage(m, result, duration)
	} else if m.Type == models.TypeDNS {
		var config models.DNSConfig
		var dnsType string
		if err := json.Unmarshal(m.Config, &config); err == nil {
			dnsType = config.RecordType
		} else {
			dnsType = "N/A"
		}

		if result.Status == models.StatusUp {
			subject, body = templates.BuildEmailDNSDetectedMessage(m, result, dnsType)
		} else if result.Status == models.StatusDown && result.ResultValue != "" {
			subject, body = templates.BuildEmailDNSChangedMessage(m, result, dnsType)
		} else {
			subject, body = templates.BuildEmailDNSStatusMessage(m, result, dnsType)
		}
	}

	return s.Send(to, subject, body)
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

func (t *TelegramService) SendStatusAlert(chatID string, m models.Monitor, result models.CheckResult, duration time.Duration) error {
	var subject, body string

	if m.Type == models.TypeHTTP {
		subject, body = templates.BuildTelegramHTTPMessage(m, result, duration)
	} else if m.Type == models.TypeDNS {
		var config models.DNSConfig
		json.Unmarshal(m.Config, &config)

		if result.Status == models.StatusUp {
			subject, body = templates.BuildTelegramDNSDetectedMessage(m, result, config.RecordType)
		} else if result.Status == models.StatusDown && result.ResultValue != "" {
			subject, body = templates.BuildTelegramDNSChangedMessage(m, result, config.RecordType)
		} else {
			subject, body = templates.BuildTelegramDNSStatusMessage(m, result, config.RecordType)
		}
	}

	return t.Send(chatID, subject, body)
}
