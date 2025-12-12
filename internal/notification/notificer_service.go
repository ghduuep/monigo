package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification/templates"
	"net/http"
	"net/smtp"
	"time"
)

type Notifier interface {
	Send(to, subject, body string) error
}

type EmailService struct {
	Host     string
	Password string
	Port     string
	Username string
	Sender   string
}

func NewEmailService(host, password, port, username, sender string) *EmailService {
	return &EmailService{
		Host:     host,
		Password: password,
		Port:     port,
		Username: username,
		Sender:   sender,
	}
}

func (s *EmailService) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	headers := make(map[string]string)
	headers["From"] = s.Sender
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += "\r\n" + body

	return smtp.SendMail(addr, auth, s.Sender, []string{to}, []byte(message))
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
