package monitor

import (
	"context"
	"log"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool) {
	activeMonitors := make(map[int]context.CancelFunc)

	for {
		monitors, err := database.GetAllMonitors(ctx, db)
		if err != nil {
			log.Printf("Failed to fetch monitors: %v", err)
			time.Sleep(10 * time.Second)
			continue
		}

		currentMonitorIDs := make(map[int]bool)
		for _, m := range monitors {
			currentMonitorIDs[m.ID] = true
			if _, exists := activeMonitors[m.ID]; !exists {
				monitorCtx, cancel := context.WithCancel(ctx)
				activeMonitors[m.ID] = cancel
				go runMonitorRoutine(monitorCtx, db, *m)
				log.Printf("Started monitoring for monitor ID %d", m.ID)
			}
		}

		for id, cancel := range activeMonitors {
			if _, exists := currentMonitorIDs[id]; !exists {
				cancel()
				delete(activeMonitors, id)
				log.Printf("Stopped monitoring for monitor ID %d", id)
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func runMonitorRoutine(ctx context.Context, db *pgxpool.Pool, m models.Monitor) {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result := performCheck(m)
			if err := database.CreateCheckResult(ctx, db, &result); err != nil {
				log.Printf("Failed to save check result for monitor %d: %v", m.ID, err)
			}
			if result.Status != m.LastCheckStatus && result.Status != models.StatusUnknown {
				log.Printf("Monitor %d status changed from %s to %s", m.ID, m.LastCheckStatus, result.Status)
				m.LastCheckStatus = result.Status
				if err := database.UpdateMonitorStatus(ctx, db, m.ID, string(result.Status)); err != nil {
					log.Printf("Failed to update monitor %d: %v", m.ID, err)
				}
			}
		}
	}
}

func performCheck(m models.Monitor) models.CheckResult {
	switch m.Type {
	case models.TypeHTTP:
		return checkHTTP(m)
	case models.TypeDNS_A, models.TypeDNS_AAAA, models.TypeDNS_MX, models.TypeDNS_NS:
		return checkDNS(m)
	default:
		return models.CheckResult{
			Status: models.StatusDown, Message: "Unknown type",
		}
	}
}
