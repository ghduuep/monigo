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

type MonitorControl struct {
	Cancel context.CancelFunc
	Data   models.Website
}

func hasConfigChanged(current, new models.Website) bool {
	if current.URL != new.URL {
		return true
	}

	if current.Interval != new.Interval {
		return true
	}

	return false
}

func StartHttpMonitoring(ctx context.Context, db *pgxpool.Pool) {
	monitoringMap := make(map[int]MonitorControl)

	for {
		websites, err := database.GetAllWebsites(ctx, db)
		if err != nil {
			log.Printf("[ERRO] failed to fetch websites: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		validIds := make(map[int]bool)

		for _, site := range websites {
			validIds[site.ID] = true

			existingMonitor, exists := monitoringMap[site.ID]

			if exists && hasConfigChanged(existingMonitor.Data, *site) {
				existingMonitor.Cancel()
				delete(monitoringMap, site.ID)
				log.Printf("[INFO] configuration changed for %s", site.URL)
				exists = false
			}

			if !exists {
				siteCtx, cancel := context.WithCancel(ctx)

				monitoringMap[site.ID] = MonitorControl{
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
				go func(userEmail, url, subject, message string) {
					if err := notification.SendEmailNotification(userEmail, url, subject, message); err != nil {
						log.Printf("[ERRO] failed to send email to %s for %s: %v", userEmail, url, err)
					} else {
						log.Printf("[INFO] sent email to %s for %s: %s", userEmail, url, message)
					}
				}(userEmail, site.URL, site.URL+"is "+newStatus, "The status of "+site.URL+" has changed to "+newStatus+".")
			}

			if err = database.CreateUptimeLog(ctx, db, site.ID, newStatus); err != nil {
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
