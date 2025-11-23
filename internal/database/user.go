package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
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

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email=$1`

	var user models.User
	err := db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(ctx context.Context, db *pgxpool.Pool, email string, passwordHash string) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2)`
	_, err := db.Exec(ctx, query, email, passwordHash)
	return err
}
