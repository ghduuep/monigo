package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCheckResult(ctx context.Context, db *pgxpool.Pool, result *models.CheckResult) error {
	query := `INSERT INTO check_results (monitor_id, status, latency_ms, status_code, result_value, message, checked_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id`
	err := db.QueryRow(ctx, query, result.MonitorID, result.Status, result.Latency, result.StatusCode, result.ResultValue, result.Message).Scan(&result.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetLatestCheckResults(ctx context.Context, db *pgxpool.Pool, monitorID int) ([]*models.CheckResult, error) {
	query := `SELECT id, monitor_id, status, latency_ms, status_code, result_value, message, checked_at 
		FROM check_results 
		WHERE monitor_id = $1 
		ORDER BY checked_at DESC 
		LIMIT $2`

	rows, err := db.Query(ctx, query, monitorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.CheckResult])
	if err != nil {
		return nil, err
	}
	return checks, nil
}
