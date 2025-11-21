package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ghduuep/pingly/internal/api/routes"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/monitor"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	db = database.InitDB()

	defer db.Close()

	router := routes.NewRouter(db)

	ctx := context.Background()

	go monitor.StartMonitoring(ctx, db)

	log.Println("API server is running on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	select {}

}
