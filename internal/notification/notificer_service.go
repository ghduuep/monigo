package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification/templates"
	"github.com/resend/resend-go/v3"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
	"net/http"
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
		To:      []string{to},
		From:    s.Sender,
		Subject: subject,
		Html:    body,
	}

	_, err := s.Client.Emails.Send(params)
	return err
}

func (s *EmailService) SendStatusAlert(to string, m models.Monitor, result models.CheckResult, inc *models.Incident) error {
	var subject, body string

	if m.Type == models.TypeHTTP {
		subject, body = templates.BuildEmailHTTPMessage(m, result, inc)
	} else if m.Type == models.TypeDNS {
		var config models.DNSConfig
		var dnsType string
		if err := json.Unmarshal(m.Config, &config); err == nil {
			dnsType = config.RecordType
		} else {
			dnsType = "N/A"
		}

		if result.Status == models.StatusUp {
			subject, body = templates.BuildEmailDNSRecoveredMessage(m, result, dnsType, inc)
		} else if result.Status == models.StatusDown && result.ResultValue != "" {
			subject, body = templates.BuildEmailDNSChangedMessage(m, result, dnsType)
		} else {
			subject, body = templates.BuildEmailDNSStatusMessage(m, result, dnsType)
		}
	} else if m.Type == models.TypePort {
		subject, body = templates.BuildEmailPortMessage(m, result, inc)
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
		"chat_id":    to,
		"text":       msg,
		"parse_mode": "Markdown",
	}

	data, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

func (t *TelegramService) SendStatusAlert(chatID string, m models.Monitor, result models.CheckResult, inc *models.Incident) error {
	var subject, body string

	if m.Type == models.TypeHTTP {
		subject, body = templates.BuildTelegramHTTPMessage(m, result, inc)
	} else if m.Type == models.TypeDNS {
		var config models.DNSConfig
		json.Unmarshal(m.Config, &config)

		if result.Status == models.StatusUp {
			subject, body = templates.BuildTelegramDNSRecoveredMessage(m, result, config.RecordType, inc)
		} else if result.Status == models.StatusDown && result.ResultValue != "" {
			subject, body = templates.BuildTelegramDNSChangedMessage(m, result, config.RecordType)
		} else {
			subject, body = templates.BuildTelegramDNSStatusMessage(m, result, config.RecordType)
		}

	} else if m.Type == models.TypePort {
		subject, body = templates.BuildTelegramPortMessage(m, result, inc)
	}

	return t.Send(chatID, subject, body)
}

type SMSService struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

func NewSMSService(accountSid, authToken, number string) *SMSService {
	return &SMSService{
		AccountSID: accountSid,
		AuthToken:  authToken,
		FromNumber: number,
	}
}

func (s *SMSService) Send(to, body string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: s.AccountSID,
		Password: s.AuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(s.FromNumber)
	params.SetBody(body)

	_, err := client.Api.CreateMessage(params)
	return err
}

func (s *SMSService) SendStatusAlert(to string, m models.Monitor, result models.CheckResult, inc *models.Incident) error {
	var body string

	if m.Type == models.TypeHTTP {
		body = templates.BuildSMSHTTPMessage(m, result, inc)
	} else if m.Type == models.TypeDNS {
		var config models.DNSConfig
		json.Unmarshal(m.Config, &config)

		if result.Status == models.StatusUp {
			body = templates.BuildSMSDNSRecoveredMessage(m, result, config.RecordType)
		} else if result.Status == models.StatusDown && result.ResultValue != "" {
			body = templates.BuildSMSDNSChangedMessage(m, result, config.RecordType)
		} else {
			body = templates.BuildSMSDNSStatusMessage(m, result, config.RecordType)
		}
	} else if m.Type == models.TypePort {
		body = templates.BuildSMSPortMessage(m, result, inc)
	}

	return s.Send(to, body)
}
