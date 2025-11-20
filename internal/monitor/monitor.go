package monitor

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckSite(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return resp.Status, err
	}
	defer resp.Body.Close()

	return resp.Status, nil
}

func startMonitoring(ctx context.Context, db *pgxpool.Pool, sites []*models.Website) {
	for _, site := range sites {
		status, err := CheckSite(site.URL)
		if err != nil {
			log.Printf("Cannot connect to website: %v", err)
		}
		if status != site.LastStatus && status != "UNKNOWN" {
			log.Printf("Status changed for %s: %s -> %s", site.URL, site.LastStatus, status)
			if err := database.UpdateWebsiteStatus(ctx, db, site.ID, status); err != nil {
				log.Printf("Error updating website status: %v", err)
			}
		}
	}
}
