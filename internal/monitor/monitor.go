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

func checkSite(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "DOWN", err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "UP", nil
	}

	return "DOWN", nil
}

func startMonitoring(ctx context.Context, db *pgxpool.Pool, sites []*models.Website) {
	for _, site := range sites {
		newStatus, err := checkSite(site.URL)

		if err != nil {
			newStatus = "DOWN"
		}

		if site.LastStatus != newStatus && site.LastStatus != "UNKNOWN" {
			log.Printf("Status of %s has changed to: %s", site.URL, newStatus)
			if err := database.CreateLog(ctx, db, site.ID, newStatus); err != nil {
				log.Printf("Error creating log for %s: %v", site.URL, err)
			}

			if err = database.UpdateWebsiteStatus(ctx, db, site.ID, newStatus); err != nil {
				log.Printf("Error updating status for %s: %v", site.URL, err)
			}
		}
	}
}
