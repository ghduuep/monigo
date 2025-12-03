package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCheckResult(ctx context.Context, db *pgxpool.Pool, result *models.CheckResult) error {
	query := `INSERT INTO check_results (monitor_id, status, latency_ms, status_code, result_value, message, checked_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id`
	err := db.QueryRow(ctx, query, result.MonitorID, result.Status).Scan(&result.ID)
	if err != nil {
		return err
	}
	return nil
}
