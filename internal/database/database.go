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
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);


	CREATE TABLE IF NOT EXISTS monitors (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		target TEXT NOT NULL,
		type VARCHAR(10) NOT NULL,
		expected_value TEXT,
		interval INTERVAL NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

		last_check_status VARCHAR(10) DEFAULT 'unknown',
		last_check_at TIMESTAMPTZ

		CONSTRAINT type_check CHECK (type IN ('http', 'dns', 'ping'))
	);

	CREATE TABLE IF NOT EXISTS check_results (
		id SERIAL PRIMARY KEY,
		monitor_id INTEGER REFERENCES monitors(id) ON DELETE CASCADE,
		status VARCHAR(10) NOT NULL,
		message TEXT,
		status_code INTEGER,
		latency_ms INTERVAL,
		checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT status_check CHECK (status IN ('up', 'down', 'unknown'))
	)
	`
	_, err := pool.Exec(ctx, query)
	return err
}
