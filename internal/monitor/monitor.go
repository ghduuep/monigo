package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckSite(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "DOWN", err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func startMonitoring(ctx context.Context, db *pgxpool.Pool, sites []*models.Website) {
}
