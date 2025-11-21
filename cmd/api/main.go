package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ghduuep/pingly/internal/api/routes"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/monitor"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	siteChannel := make(chan *models.Website)

	ctx := context.Background()

	router := routes.NewRouter(db, siteChannel)

	go monitor.StartMonitoring(ctx, db, siteChannel)

	log.Println("API server is running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
