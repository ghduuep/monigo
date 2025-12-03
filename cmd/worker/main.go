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
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_SENDER"),
		os.Getenv("SMTP_PASSWORD"),
	)

	db := database.InitDB()
	defer db.Close()

	log.Println("Worker is running...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitor.StartMonitoring(ctx, db, emailService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Worker is shutting down...")
}
