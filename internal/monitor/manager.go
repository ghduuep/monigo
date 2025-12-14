package monitor

import (
	"context"
	"encoding/json"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool, dispatcher notification.NotificationDispatcher) {
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
				go runMonitorRoutine(monitorCtx, db, *m, dispatcher)
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

func runMonitorRoutine(ctx context.Context, db *pgxpool.Pool, m models.Monitor, dispatcher notification.NotificationDispatcher) {
	const downInterval = 30 * time.Second

	timer := time.NewTimer(m.Interval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			currentStatus := processCheck(ctx, db, &m, dispatcher)

			nextCheckDuration := m.Interval

			if currentStatus == models.StatusDown {
				nextCheckDuration = downInterval
			}

			timer.Reset(nextCheckDuration)
		}
	}
}

func processCheck(ctx context.Context, db *pgxpool.Pool, m *models.Monitor, dispatcher notification.NotificationDispatcher) models.MonitorStatus {
	result := performCheck(*m)

	var config models.DNSConfig

	if m.Type == models.TypeDNS && result.Status == models.StatusUp && result.ResultValue != "" {
		if err := database.SetInitialDNSConfig(ctx, db, m.ID, result.ResultValue); err != nil {
			log.Printf("[ERROR] Failed to set initial DNS config: %v", err)
		}

		if err := json.Unmarshal(m.Config, &config); err != nil {
			log.Printf("[ERROR] Failed to unmarsh json config")
		} else {
			if config.ExpectedValue == "" {
				log.Printf("[INFO] Learning DNS value for monitor %d: %s", m.ID, result.ResultValue)
				config.ExpectedValue = result.ResultValue
				newJSON, _ := json.Marshal(config)
				m.Config = newJSON
			}
		}
	}

	if err := database.CreateCheckResult(ctx, db, &result); err != nil {
		log.Printf("[ERROR] failed to save check result for monitor %d: %v", m.ID, err)
	}

	if result.Status != m.LastCheckStatus && result.Status != models.StatusUnknown {
		log.Printf("[LOG] Monitor %d has changed from %s to %s", m.ID, m.LastCheckStatus, result.Status)
		var incident *models.Incident
		var dbErr error

		if result.Status == models.StatusDown {
			incident, dbErr = database.CreateIncident(ctx, db, m.ID, result.Message)
			if dbErr != nil {
				log.Printf("[ERROR] Failed to create incident: %v", dbErr)
			}
		}

		if m.StatusChangedAt != nil && m.LastCheckStatus == models.StatusDown && result.Status == models.StatusUp {
			incident, dbErr = database.ResolveIncident(ctx, db, m.ID)
			if dbErr != nil {
				log.Printf("[ERROR] Failed to resolve incident: %v", dbErr)
			}
		}

		if err := database.UpdateMonitorStatus(ctx, db, m.ID, string(result.Status)); err != nil {
			log.Printf("[ERROR] Failed to update monitor %d: %v", m.ID, err)
		}

		if incident != nil {
			channels, _ := database.GetUserChannels(ctx, db, m.UserID)
			go dispatcher.SendAlert(channels, *m, result, incident)
		}

		m.LastCheckStatus = result.Status
		m.StatusChangedAt = &result.CheckedAt
	} else {
		if err := database.UpdateLastCheck(ctx, db, m.ID); err != nil {
			log.Printf("[ERROR] Failed to update last check for monitor %d.: %v", m.ID, err)
		}
	}

	return result.Status
}

func performCheck(m models.Monitor) models.CheckResult {
	switch m.Type {
	case models.TypeHTTP:
		return checkHTTP(m)
	case models.TypeDNS:
		return checkDNS(m)
	case models.TypePort:
		return checkPort(m)
	default:
		return models.CheckResult{
			Status: models.StatusDown, Message: "Unknown type",
		}
	}
}
