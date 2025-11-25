package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllDNSMonitors(ctx context.Context, db *pgxpool.Pool) ([]*models.DNSDomains, error) {
	query := `SELECT id, domain, interval, last_a_records, last_aaaa_records, last_mx_records, last_ns_records, last_checked FROM dns_monitors`
	rows, _ := db.Query(ctx, query)

	dnsMonitors, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.DNSDomains])

	return dnsMonitors, err
}

func CreateDNSMonitor(ctx context.Context, db *pgxpool.Pool, dnsMonitor *models.DNSDomains) error {
	query := `INSERT INTO dns_monitors (domain, interval, last_a_records, last_aaaa_records, last_mx_records, last_ns_records, last_checked) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(ctx, query, dnsMonitor.Domain, dnsMonitor.Interval, dnsMonitor.LastA, dnsMonitor.LastAAAA, dnsMonitor.LastMX, dnsMonitor.LastNS, dnsMonitor.LastChecked).Scan(&dnsMonitor.ID)

	return err
}

func UpdateDNSMonitorRecords(ctx context.Context, db *pgxpool.Pool, dnsMonitorID int, lastA, lastAAAA, lastMX, lastNS []string) error {
	query := `UPDATE dns_monitors SET last_a_records = $1, last_aaaa_records, last_mx_records = $2, last_ns_records = $3, last_checked = NOW() WHERE id = $4`

	_, err := db.Exec(ctx, query, lastA, lastAAAA, lastMX, lastNS, dnsMonitorID)

	return err
}

func DeleteDNSMonitor(ctx context.Context, db *pgxpool.Pool, dnsMonitorID int) error {
	query := `DELETE FROM dns_monitors WHERE id = $1`

	_, err := db.Exec(ctx, query, dnsMonitorID)

	return err
}
