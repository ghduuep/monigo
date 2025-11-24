package monitor

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool) {
	monitoringMap := make(map[int]models.MonitorControl)

	for {
		websites, err := database.GetAllWebsites(ctx, db)
		if err != nil {
			log.Printf("[ERRO] failed to fetch websites: %v", err)
			continue
		}

		validIds := make(map[int]bool)

		for _, site := range websites {
			validIds[site.ID] = true

			existingMonitor, exists := monitoringMap[site.ID]

			if exists {
				if existingMonitor.Data.URL != site.URL || existingMonitor.Data.Interval != site.Interval {
					log.Printf("[INFO] Change detected for %s", site.URL)
					existingMonitor.Cancel()
					delete(monitoringMap, site.ID)
					exists = false
				}
			}

			if !exists {
				siteCtx, cancel := context.WithCancel(ctx)

				monitoringMap[site.ID] = models.MonitorControl{
					Cancel: cancel,
					Data:   *site,
				}

				go monitor(siteCtx, db, site)
				log.Printf("[INFO] started monitoring %s", site.URL)
			}
		}

		for id, control := range monitoringMap {
			if _, ok := validIds[id]; !ok {
				control.Cancel()
				delete(monitoringMap, id)
				log.Printf("[INFO] stopped monitoring website ID %d", id)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func monitor(ctx context.Context, db *pgxpool.Pool, site *models.Website) {
	ticker := time.NewTicker(site.Interval)
	defer ticker.Stop()

	for {
		newStatus, err := checkSite(site.URL)
		log.Printf("[INFO] checked %s: %s", site.URL, newStatus)
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
					} else {
						log.Printf("[INFO] sent email to %s for %s: %s", userEmail, url, status)
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
