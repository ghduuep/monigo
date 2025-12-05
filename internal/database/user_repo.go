package database

import (
	"context"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, db *pgxpool.Pool, user *models.User) error {
	query := `INSERT INTO users(username, email, password_hash, created_at) VALUES($1, $2, $3, $4) RETURNING id`

	err := db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.CreatedAt).Scan(&user.ID)
	return err
}

func GetAllUsers(ctx context.Context, db *pgxpool.Pool) ([]models.User, error) {
	query := `SELECT * FROM users`

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, err
	}
	return users, nil
}

func DeleteUser(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func GetUserEmailByID(ctx context.Context, db *pgxpool.Pool, id int) (string, error) {
	query := `SELECT email FROM users WHERE id=$1`

	var email string
	err := db.QueryRow(ctx, query, id).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

func GetUserByID(ctx context.Context, db *pgxpool.Pool, id int) (models.User, error) {
	query := `SELECT * FROM users WHERE id=$1`
	var user models.User
	err := db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, nil
	}
	return user, nil
}

func GetUserByUsername(ctx context.Context, db *pgxpool.Pool, username string) (models.User, error) {
	query := `SELECT * FROM users WHERE username = $1`

	var user models.User
	err := db.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return models.User{}, err
	}

	return user, err
}
