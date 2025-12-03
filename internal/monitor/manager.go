package monitor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool, emailService *notification.EmailService) {
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
				go runMonitorRoutine(monitorCtx, db, *m, emailService)
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

func runMonitorRoutine(ctx context.Context, db *pgxpool.Pool, m models.Monitor, emailService *notification.EmailService) {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			processCheck(ctx, db, &m, emailService)
		}
	}
}

func handleAutoDiscovery(ctx context.Context, db *pgxpool.Pool, m *models.Monitor, res models.CheckResult) {
	if m.Type != models.TypeDNS || res.Status != models.StatusUp {
		return
	}

	var dnsConfig models.DNSConfig
	if err := json.Unmarshal(m.Config, &dnsConfig); err != nil {
		return
	}

	if dnsConfig.ExpectedValue != "" {
		return
	}

	log.Printf("[LOG] Auto Discovery: learning value '%s' for monitor %d", res.ResultValue, m.ID)

	dnsConfig.ExpectedValue = res.ResultValue
	newConfigJson, _ := json.Marshal(dnsConfig)

	m.Config = newConfigJson

	if err := database.UpdateMonitorConfig(ctx, db, m.ID, newConfigJson); err != nil {
		log.Printf("[ERROR] Failed to save auto config: %v", err)
	}
}

func processCheck(ctx context.Context, db *pgxpool.Pool, m *models.Monitor, emailService *notification.EmailService) {
	result := performCheck(*m)

	handleAutoDiscovery(ctx, db, m, result)

	if err := database.CreateCheckResult(ctx, db, &result); err != nil {
		log.Printf("[ERROR] failed to save check result for monitor %d: %v", m.ID, err)
	}

	if result.Status != m.LastCheckStatus && result.Status != models.StatusUnknown {
		log.Printf("[LOG] Monitor %d has changed from %s to %s", m.ID, m.LastCheckStatus, result.Status)

		m.LastCheckStatus = result.Status
		if err := database.UpdateMonitorStatus(ctx, db, m.ID, string(result.Status)); err != nil {
			log.Printf("[ERROR] Failed to update monitor %d: %v", m.ID, err)
		}

		userEmail, _ := database.GetUserEmailByID(ctx, db, m.UserID)

		if userEmail != "" {
			go func(mon models.Monitor, res models.CheckResult) {
				if err := emailService.SendStatusAlert(userEmail, mon, res); err != nil {
					log.Printf("[ERROR] Failed to send e-mail: %v", err)
				}
			}(*m, result)
		}
	}
}

func performCheck(m models.Monitor) models.CheckResult {
	switch m.Type {
	case models.TypeHTTP:
		return checkHTTP(m)
	case models.TypeDNS:
		return checkDNS(m)
	default:
		return models.CheckResult{
			Status: models.StatusDown, Message: "Unknown type",
		}
	}
}
