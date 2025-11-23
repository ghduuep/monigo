package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/ghduuep/pingly/internal/api/routes"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/monitor"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	db := database.InitDB()
	defer db.Close()

	siteChannel := make(chan *models.Website)

	ctx := context.Background()

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	router := routes.NewRouter(db, siteChannel, tokenAuth)

	go monitor.StartMonitoring(ctx, db, siteChannel)

	log.Println("API server is running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
