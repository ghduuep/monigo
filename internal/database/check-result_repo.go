package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateCheckResult(ctx context.Context, db *pgxpool.Pool, result *models.CheckResult) error {
	query := `INSERT INTO check_results (monitor_id, status, message, status_code, latency_ms, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	err := db.QueryRow(ctx, query, result.MonitorID, result.Status, result.Message, result.StatusCode, result.Latency).Scan(&result.ID)
	if err != nil {
		return err
	}
	return nil
}
