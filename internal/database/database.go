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
	queryStandard := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS user_channels (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		type VARCHAR(20) NOT NULL,
		target VARCHAR(255) NOT NULL,
		enabled BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		CONSTRAINT unique_channel_target UNIQUE (user_id, type, target)
	);

	CREATE TABLE IF NOT EXISTS monitors (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		target TEXT NOT NULL,
		type VARCHAR(10) NOT NULL,
		config JSONB DEFAULT '{}'::jsonb,
		interval INTERVAL NOT NULL,
		timeout INTERVAL NOT NULL DEFAULT INTERVAL '30 seconds',
		latency_threshold_ms INTEGER DEFAULT 0,
		last_check_status VARCHAR(10) DEFAULT 'unknown',
		last_check_at TIMESTAMPTZ,
		status_changed_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT NOW()
	);
	CREATE INDEX IF NOT EXISTS idx_monitors_user_id ON monitors(user_id);

	CREATE TABLE IF NOT EXISTS incidents (
		id SERIAL PRIMARY KEY,
		monitor_id INTEGER REFERENCES monitors(id) ON DELETE CASCADE,
		started_at TIMESTAMPTZ DEFAULT NOW(),
		resolved_at TIMESTAMPTZ,
		duration INTERVAL,
		error_cause TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_incidents_monitor_id ON incidents(monitor_id);
	`
	if _, err := pool.Exec(ctx, queryStandard); err != nil {
		return err
	}

	queryTimescaleBase := `
	CREATE EXTENSION IF NOT EXISTS timescaledb;

	CREATE TABLE IF NOT EXISTS check_results (
		id BIGSERIAL, 
		monitor_id INTEGER REFERENCES monitors(id) ON DELETE CASCADE,
		status VARCHAR(10) NOT NULL,
		result_value TEXT,
		message TEXT,
		status_code INTEGER,
		latency_ms INTEGER,
		checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`
	if _, err := pool.Exec(ctx, queryTimescaleBase); err != nil {
		return err
	}

	queryHypertable := `
	SELECT create_hypertable('check_results', 'checked_at', chunk_time_interval => INTERVAL '1 day', if_not_exists => TRUE);
	`
	if _, err := pool.Exec(ctx, queryHypertable); err != nil {
		return err
	}

	queryCompression := `
	CREATE INDEX IF NOT EXISTS idx_check_results_monitor_date ON check_results(monitor_id, checked_at DESC);

	ALTER TABLE check_results SET (
		timescaledb.compress,
		timescaledb.compress_segmentby = 'monitor_id',
		timescaledb.compress_orderby = 'checked_at DESC'
	);

	SELECT add_compression_policy('check_results', INTERVAL '3 days', if_not_exists => TRUE);
	SELECT add_retention_policy('check_results', INTERVAL '1 year', if_not_exists => TRUE);
	`
	if _, err := pool.Exec(ctx, queryCompression); err != nil {
		return err
	}

	queryView := `
	CREATE MATERIALIZED VIEW IF NOT EXISTS monitor_stats_hourly
	WITH (timescaledb.continuous) AS
	SELECT
		time_bucket('1 hour', checked_at) as bucket,
		monitor_id,
		SUM(latency_ms) as sum_latency,
		MIN(latency_ms) as min_latency,
		MAX(latency_ms) as max_latency,
		COUNT(*) FILTER (WHERE status = 'up') as up_count,
		COUNT(*) as total_checks
	FROM check_results
	GROUP BY bucket, monitor_id;
	`
	if _, err := pool.Exec(ctx, queryView); err != nil {
		return err
	}

	queryViewPolicy := `
	SELECT add_continuous_aggregate_policy('monitor_stats_hourly',
		start_offset => NULL,
		end_offset => INTERVAL '5 minutes',
		schedule_interval => INTERVAL '5 minutes',
		if_not_exists => TRUE);
	`
	if _, err := pool.Exec(ctx, queryViewPolicy); err != nil {
		return err
	}

	return nil
}
