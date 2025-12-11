package notification

import (
	"fmt"
	"log"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

type NotificationDispatcher struct {
	Email    *EmailService
	Telegram *TelegramService
}

func NewDispatcher(email *EmailService, telegram *TelegramService) *NotificationDispatcher {
	return &NotificationDispatcher{
		Email:    email,
		Telegram: telegram,
	}
}

func (d *NotificationDispatcher) SendAlert(channels []models.NotificationChannel, m models.Monitor, res models.CheckResult, duration time.Duration) {
	for _, ch := range channels {
		var err error

		switch ch.Type {
		case models.TypeEmail:
			go func(target string) {
				d.Email.SendStatusAlert(target, m, res, duration)
			}(ch.Target)

		case models.TypeTelegram:
			go func(target string) {
				d.Telegram.SendStatusAlert(target, m, res, duration)
			}(ch.Target)

		default:
			log.Printf("[ERROR] Unknown channel: %s", ch.Type)
		}

		if err != nil {
			log.Printf("[ERROR] Failed to send %s notification to %s", ch.Type, ch.Target)
		}
	}
}
