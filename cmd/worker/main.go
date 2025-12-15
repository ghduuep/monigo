package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/monitor"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	emailService := notification.NewEmailService(
		os.Getenv("RESEND_API_KEY"),
		os.Getenv("EMAIL_SENDER"),
	)

	telegramService := notification.NewTelegramService(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
	)

	smsService := notification.NewSMSService(
		os.Getenv("TWILIO_ACCOUNT_SID"),
		os.Getenv("TWILIO_AUTH_TOKEN"),
		os.Getenv("TWILIO_NUMBER"),
	)

	dispatcher := notification.NewDispatcher(emailService, telegramService, smsService)

	db := database.InitDB()
	defer db.Close()

	log.Println("Worker is running...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitor.StartMonitoring(ctx, db, *dispatcher)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Worker is shutting down...")
}
