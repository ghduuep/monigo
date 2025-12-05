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
		log.Fatalf("[ERROR] Failed to create pool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("[ERROR] Failed to connect with database: %v", err)
	}

	log.Println("[INFO] Connected successfuly with database")

	err = createTables(context.Background(), pool)

	if err != nil {
		log.Fatalf("[ERROR] Failed to create database tables: %v", err)
	}

	return pool
}

func createTables(ctx context.Context, pool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);


	CREATE TABLE IF NOT EXISTS monitors (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		target TEXT NOT NULL,
		type VARCHAR(10) NOT NULL,
		config JSONB DEFAULT '{}'::jsonb,
		interval INTERVAL NOT NULL,
		last_check_status VARCHAR(10) DEFAULT 'unknown',
		last_check_at TIMESTAMPTZ,
		status_changed_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW(),

		CONSTRAINT type_check CHECK (type IN ('http', 'dns', 'ping'))
	);

	CREATE TABLE IF NOT EXISTS check_results (
		id SERIAL PRIMARY KEY,
		monitor_id INTEGER REFERENCES monitors(id) ON DELETE CASCADE,
		status VARCHAR(10) NOT NULL,
		result_value TEXT,
		message TEXT,
		status_code INTEGER,
		latency_ms INTEGER,
		checked_at TIMESTAMPTZ DEFAULT NOW(),

		CONSTRAINT status_check CHECK (status IN ('up', 'down', 'unknown'))
	)
	`
	_, err := pool.Exec(ctx, query)
	return err
}
