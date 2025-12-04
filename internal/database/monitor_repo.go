package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllMonitors(ctx context.Context, db *pgxpool.Pool) ([]*models.Monitor, error) {
	query := `SELECT id, user_id, target, type, config, interval, last_check_status, last_check_at, created_at FROM monitors`
	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monitors, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Monitor])
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

func CreateMonitor(ctx context.Context, db *pgxpool.Pool, monitor *models.Monitor) error {
	query := `INSERT INTO monitors (user_id, target, type, config, interval) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := db.QueryRow(ctx, query, monitor.UserID, monitor.Target, monitor.Type, monitor.Config, monitor.Interval).Scan(&monitor.ID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMonitorStatus(ctx context.Context, db *pgxpool.Pool, monitorID int, status string) error {
	query := `UPDATE monitors SET last_check_status = $1, last_check_at = NOW() WHERE id = $2`
	_, err := db.Exec(ctx, query, status, monitorID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateMonitorConfig(ctx context.Context, db *pgxpool.Pool, monitorID int, newConfig []byte) error {
	query := `UPDATE monitors SET config = $1 WHERE id = $2`
	_, err := db.Exec(ctx, query, newConfig, monitorID)
	return err
}

func DeleteMonitor(ctx context.Context, db *pgxpool.Pool, monitorID int64) error {
	query := `DELETE FROM monitors WHERE id = $1`
	_, err := db.Exec(ctx, query, monitorID)
	if err != nil {
		return err
	}
	return nil
}

func SetInitialDNSConfig(ctx context.Context, db *pgxpool.Pool, monitorID int, detectedValue string) error {
	query := `UPDATE monitors SET config = config || jsonb_build_object('expected_value', $1::text) WHERE id = $2
	AND (config->>'expected_value' IS NULL or config->>'expected_value' = '')`

	_, err := db.Exec(ctx, query, detectedValue, monitorID)
	return err
}
