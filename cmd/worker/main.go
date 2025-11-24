package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/monitor"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	db := database.InitDB()
	defer db.Close()

	log.Println("Worker is running...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go monitor.StartHttpMonitoring(ctx, db)
	go monitor.StartDNSMonitoring(ctx, db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Worker is shutting down...")
}
