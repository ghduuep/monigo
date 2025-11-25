package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUptimeLog(ctx context.Context, db *pgxpool.Pool, websiteID int, status string) error {
	query := `INSERT INTO uptime_logs (website_id, status) VALUES ($1, $2)`

	_, err := db.Exec(ctx, query, websiteID, status)

	return err
}

func CreateDNSLog(ctx context.Context, db *pgxpool.Pool, dnsMonitorID int, details []string) error {
	query := `INSERT INTO dns_logs (dns_monitor_id, details) VALUES ($1, $2, $3)`

	_, err := db.Exec(ctx, query, dnsMonitorID, details)

	return err
}
