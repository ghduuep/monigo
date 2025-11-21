package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserEmail(ctx context.Context, db *pgxpool.Pool, userID int) (string, error) {
	query := `SELECT email FROM users WHERE id=$1`
	var email string
	err := db.QueryRow(ctx, query, userID).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}
