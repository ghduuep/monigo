package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateLog(ctx context.Context, db *pgxpool.Pool, websiteID int, status string) error {
	query := `INSERT INTO check_logs (website_id, status) VALUES ($1, $2)`

	_, err := db.Exec(ctx, query, websiteID, status)

	return err
}
