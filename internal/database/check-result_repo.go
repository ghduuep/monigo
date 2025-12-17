package database

import (
	"context"
	"time"

	"github.com/ghduuep/pingly/internal/dto"
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

func GetMonitorStats(ctx context.Context, db *pgxpool.Pool, monitorID int, threshold int64, from, to time.Time) (dto.MonitorStatsResponse, error) {
	query := `SELECT
				COALESCE(SUM(total_checks), 0) as total_checks,
				COALESCE(SUM(sum_latency) / NULLIF(SUM(total_checks), 0), 0) as avg_latency,
				COALESCE(MIN(min_latency), 0) as min_latency,
				COALESCE(MAX(max_latency), 0) as max_latency,
				COALESCE(
					(SUM(up_count) * 100.0 / NULLIF(SUM(total_checks), 0)),
					0
				) as uptime_percentage
			FROM monitor_stats_hourly
			WHERE monitor_id = $1 
			AND bucket >= $2 AND bucket <= $3`

	var stats dto.MonitorStatsResponse
	stats.MonitorID = monitorID

	err := db.QueryRow(ctx, query, monitorID, from, to).Scan(
		&stats.TotalChecks,
		&stats.AvgLatency,
		&stats.MinLatency,
		&stats.MaxLatency,
		&stats.UptimePercentage,
	)

	if err != nil {
		return dto.MonitorStatsResponse{}, err
	}

	if threshold > 0 {
		var satisfactory, tolerating, totalApdexChecks int64

		queryApdex := `
				SELECT
                COUNT(*) FILTER (WHERE status IN ('up', 'degraded') AND latency_ms <= $1),
                COUNT(*) FILTER (WHERE status IN ('up', 'degraded') AND latency_ms > $1 AND latency_ms <= ($1 * 4)),
                COUNT(*)
            FROM check_results
            WHERE monitor_id = $2 AND checked_at >= $3 AND checked_at <= $4
            `

		err := db.QueryRow(ctx, queryApdex, threshold, monitorID, from, to).Scan(&satisfactory, &tolerating, &totalApdexChecks)

		if err == nil && totalApdexChecks > 0 {
			score := (float64(satisfactory) + (float64(tolerating) / 2.0)) / float64(totalApdexChecks)
			stats.ApdexScore = score
		}
	}
	return stats, nil
}

func GetLastChecks(ctx context.Context, db *pgxpool.Pool, monitorID int, from, to time.Time) ([]*models.CheckResult, error) {
	query := `
	SELECT id, monitor_id, status, result_value, message, status_code, latency_ms, checked_at
	FROM check_results
	WHERE monitor_id = $1
	AND checked_at >= $2 AND checked_at <= $3
	ORDER BY checked_at DESC
	`

	rows, err := db.Query(ctx, query, monitorID, from, to)
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

func ExportCheckResults(ctx context.Context, db *pgxpool.Pool, monitorID int, from, to time.Time) (pgx.Rows, error) {
	query := `SELECT checked_at, status, latency_ms, status_code, result_value, message
	FROM check_results
	WHERE monitor_id = $1 AND checked_at >= $2 AND checked_at <= $3
	ORDER BY checked_at DESC`

	return db.Query(ctx, query, monitorID, from, to)
}
