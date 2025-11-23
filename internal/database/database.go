package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Erro ao criar pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("[ERRO] não foi possível conectar ao banco de dados: %v", err)
	}

	log.Println("[INFO] conectado com sucesso ao banco.")

	err = createTables(context.Background(), pool)

	if err != nil {
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}

	return pool
}

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS websites (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		url VARCHAR(255) NOT NULL,
		interval INTERVAL NOT NULL,
		last_checked TIMESTAMP,
		last_status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN'
	);

	CREATE TABLE IF NOT EXISTS check_logs (
		id SERIAL PRIMARY KEY,
		website_id INTEGER REFERENCES websites(id),
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := pool.Exec(ctx, query)
	return err
}
