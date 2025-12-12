package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateIncident(ctx context.Context, db *pgxpool.Pool, monitorID int, cause string) error {
	queryCheck := `SELECT id FROM incidents WHERE monitor_id = $1 AND resolved_at IS NULL`
	var existingID int
	err := db.QueryRow(ctx, queryCheck, monitorID).Scan(&existingID)
	if err == nil {
		return nil
	}

	query := `INSERT INTO incidents (monitor_id, started_at, error_cause) VALUES ($1, NOW(), $2)`
	_, err = db.Exec(ctx, query, monitorID, cause)
	return err
}

func ResolveIncident(ctx context.Context, db *pgxpool.Pool, monitorID int) error {
	query := `
        UPDATE incidents 
        SET resolved_at = NOW(),
            duration = NOW() - started_at
        WHERE monitor_id = $1 AND resolved_at IS NULL
    `
	_, err := db.Exec(ctx, query, monitorID)
	return err
}
