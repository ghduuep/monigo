package database

import (
	"context"
	"time"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateIncident(ctx context.Context, db *pgxpool.Pool, monitorID int, cause string) (*models.Incident, error) {
	var exists int
	_ = db.QueryRow(ctx, "SELECT 1 FROM incidents WHERE monitor_id = $1 AND resolved_at IS NULL", monitorID).Scan(&exists)
	if exists == 1 {
		return nil, nil // Já existe, ignora ou retorna erro
	}

	query := `INSERT INTO incidents (monitor_id, started_at, error_cause) 
	          VALUES ($1, NOW(), $2) 
	          RETURNING id, monitor_id, started_at, error_cause`

	var inc models.Incident
	err := db.QueryRow(ctx, query, monitorID, cause).Scan(&inc.ID, &inc.MonitorID, &inc.StartedAt, &inc.ErrorCause)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func ResolveIncident(ctx context.Context, db *pgxpool.Pool, monitorID int) (*models.Incident, error) {
	query := `
		UPDATE incidents 
		SET resolved_at = NOW(),
		    duration = NOW() - started_at
		WHERE monitor_id = $1 AND resolved_at IS NULL
		RETURNING id, monitor_id, started_at, resolved_at, duration, error_cause
	`
	var inc models.Incident
	err := db.QueryRow(ctx, query, monitorID).Scan(
		&inc.ID, &inc.MonitorID, &inc.StartedAt, &inc.ResolvedAt, &inc.Duration, &inc.ErrorCause,
	)
	if err != nil {
		return nil, err
	}
	return &inc, nil
}

func GetIncidentsByID(ctx context.Context, db *pgxpool.Pool, monitorID int, from, to time.Time) ([]*models.Incident, error) {
	query := `SELECT * from incidents WHERE monitor_id = $1 AND started_at >= $2 AND started_at <= $3 ORDER BY started_at DESC`

	rows, err := db.Query(ctx, query, monitorID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Incident])
	if err != nil {
		return nil, err
	}

	return results, nil
}
