package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetUserEmailByID(ctx context.Context, db *pgxpool.Pool, id int) (string, error) {
	query := `SELECT email FROM users WHERE id=$1`

	var email string
	err := db.QueryRow(ctx, query, id).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}
