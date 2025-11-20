package database

import (
	"context"

	"github.com/ghduuep/pingly/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetWebsites(ctx context.Context, db *pgxpool.Pool) ([]models.Website, error){
	rows, _ := db.Query(ctx, "SELECT id, user_id, url, interval, last_checked, last_status FROM websites")

	websites, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Website])

	return websites, err
}

func CreateWebsite(ctx context.Context, db *pgxpool.Pool, website *models.Website) error {
	query := `INSERT INTO websites (user_id, url, interval, last_checked, last_status) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := db.QueryRow(ctx, query, website.UserID, website.URL, website.Interval, website.LastChecked, website.LastStatus).Scan(&website.ID)

	return err
}
