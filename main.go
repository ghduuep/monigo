package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Website struct {
	URL        string
	Interval   time.Duration
	LastStatus string
	CheckedAt time.Time
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Erro ao criar pool: %v", err)
	}

	defer pool.Close()


	err = createTable(context.Background(), pool)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

}

func createTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
		CREATE TABLE IF NOT EXISTS checks (
			id SERIAL PRIMARY KEY,
			url TEXT NOT NULL,
			status TEXT NOT NULL,
			checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	_, err := pool.Exec(ctx, query)
	return err
}

// func monitorLoop(website *Website) {
// 	for {
// 		newStatus, err := Check(website.URL)

// 		if err != nil {
// 			log.Printf("[%s] Erro: %v", website.URL, err)
// 		}

// 		if website.LastStatus != "UNKNOWN" && website.LastStatus != newStatus {
// 			log.Printf("MUDANÇA DE STATUS: %s - %s", website.URL, newStatus)

// 			go func() {
// 				err := sendEmailNotification(website.URL, newStatus)
// 				if err != nil {
// 					log.Printf("Erro ao enviar e-mail de aviso: %v", err)
// 				}
// 			}()
// 		} else {
// 			log.Printf("[%s] Status: %s", website.URL, newStatus)
// 		}

// 		website.LastStatus = newStatus
// 		time.Sleep(website.Interval)
// 	}
// }
