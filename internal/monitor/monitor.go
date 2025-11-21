package monitor

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/ghduuep/pingly/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool, sites []*models.Website) {
	for _, site := range sites {
		go monitor(ctx, db, site)
	}
}

func monitor(ctx context.Context, db *pgxpool.Pool, site *models.Website) {
	ticker := time.NewTicker(site.Interval)
	defer ticker.Stop()

	for {
		newStatus, err := checkSite(site.URL)
		if err != nil {
			log.Printf("[ERRO] %s: %v", site.URL, err)
			newStatus = "DOWN"
		}

		if site.LastStatus != newStatus || site.LastStatus == "UNKNOWN" {
			log.Printf("[INFO] status has changed for %s: %s", site.URL, newStatus)
			userEmail, err := database.GetUserEmail(ctx, db, site.UserID)

			if err != nil {
				log.Printf("[ERRO] failed to get user email for %s: %v", site.URL, err)
			} else {
				go func(userEmail, url, status string) {
					if err := notification.SendEmailNotification(userEmail, url, status); err != nil {
						log.Printf("[ERRO] failed to send email to %s for %s: %v", userEmail, url, err)
					}
				}(userEmail, site.URL, newStatus)
			}

			if err = database.CreateLog(ctx, db, site.ID, newStatus); err != nil {
				log.Printf("[ERRO] failed to create log for %s: %v", site.URL, err)
			}

			if err = database.UpdateWebsiteStatus(ctx, db, site.ID, newStatus); err != nil {
				log.Printf("[ERRO] failed to update status for %s: %v", site.URL, err)
			} else {
				site.LastStatus = newStatus
			}
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			continue
		}
	}
}

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
