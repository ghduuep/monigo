package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
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

func GetUserByEmail(ctx context.Context, db *pgxpool.Pool, email string) (models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user models.User
	err := db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, err
}
