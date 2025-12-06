package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/dto"
	"github.com/ghduuep/pingly/internal/models"
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

func GetMonitorStats(ctx context.Context, db *pgxpool.Pool, monitorID int) (dto.MonitorStatsResponse, error) {
	query := `SELECT
				COUNT(*) FILTER (WHERE checked_at > NOW() - INTERVAL '24 hours') as last_24_checks,
				COALESCE(AVG(latency_ms), 0) as avg_latency,
				COALESCE(
					(COUNT(*) FILTER (WHERE status = 'up') * 100.0 / NULLIF(COUNT(*), 0)),
					0
				) as uptime_percentage
			FROM check_results
			WHERE monitor_id = $1
			AND checked_at > NOW() - INTERVAL '30 days'
			`

	var stats dto.MonitorStatsResponse
	stats.MonitorID = monitorID

	err := db.QueryRow(ctx, query, monitorID).Scan(
		&stats.Last24HChecks,
		&stats.AvgLatency,
		&stats.UptimePercentage,
	)

	if err != nil {
		return dto.MonitorStatsResponse{}, err
	}

	return stats, nil
}
