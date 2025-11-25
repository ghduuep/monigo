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
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		url VARCHAR(255) NOT NULL,
		interval INTERVAL NOT NULL DEFAULT '5 minutes',
		last_checked TIMESTAMP,
		last_status VARCHAR(50) NOT NULL DEFAULT 'UNKNOWN',

		UNIQUE(user_id, url)
	);

	CREATE TABLE IF NOT EXISTS dns_domains (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		domain VARCHAR(255) NOT NULL,
		interval INTERVAL NOT NULL DEFAULT '1 hour',
		last_a_records JSONB,
		last_aaaa_records JSONB,
		last_mx_records JSONB,
		last_ns_records JSONB,
		last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		UNIQUE(user_id, domain)
	);

	CREATE TABLE IF NOT EXISTS uptime_logs (
		id SERIAL PRIMARY KEY,
		website_id INTEGER REFERENCES websites(id) ON DELETE CASCADE,
		status VARCHAR(50) NOT NULL,
		root_cause VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS dns_logs (
		id SERIAL PRIMARY KEY,
		dns_monitor_id INTEGER REFERENCES dns_domains(id) ON DELETE CASCADE,
		diff JSONB,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := pool.Exec(ctx, query)
	return err
}
