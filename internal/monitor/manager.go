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
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			processCheck(ctx, db, &m, dispatcher)
		}
	}
}

func processCheck(ctx context.Context, db *pgxpool.Pool, m *models.Monitor, dispatcher notification.NotificationDispatcher) {
	result := performCheck(*m)

	var config models.DNSConfig

	if m.Type == models.TypeDNS && result.Status == models.StatusUp && result.ResultValue != "" {
		if err := database.SetInitialDNSConfig(ctx, db, m.ID, result.ResultValue); err != nil {
			log.Printf("[ERROR] Failed to set initial DNS config: %v", err)
		}

		if err := json.Unmarshal(m.Config, &config); err != nil {
			log.Printf("[ERROR] Failed to unmarsh json config")
			return
		}

		if config.ExpectedValue != "" {
			return
		}

		log.Printf("[INFO] Learning DNS value for monitor %d: %s", m.ID, result.ResultValue)

		config.ExpectedValue = result.ResultValue
		newJSON, _ := json.Marshal(config)
		m.Config = newJSON
	}

	if err := database.CreateCheckResult(ctx, db, &result); err != nil {
		log.Printf("[ERROR] failed to save check result for monitor %d: %v", m.ID, err)
	}

	if result.Status != m.LastCheckStatus && result.Status != models.StatusUnknown {
		log.Printf("[LOG] Monitor %d has changed from %s to %s", m.ID, m.LastCheckStatus, result.Status)

		var downtimeDuration time.Duration
		if m.StatusChangedAt != nil && m.LastCheckStatus == models.StatusDown && result.Status == models.StatusUp {
			downtimeDuration = result.CheckedAt.Sub(*m.StatusChangedAt)
		}

		if err := database.UpdateMonitorStatus(ctx, db, m.ID, string(result.Status)); err != nil {
			log.Printf("[ERROR] Failed to update monitor %d: %v", m.ID, err)
		}

		channels, err := database.GetUserChannels(ctx, db, m.UserID)
		if err != nil {
			log.Printf("[ERROR] Failed to fetch user channels: %v", err)
		}

		go dispatcher.SendAlert(channels, *m, result, downtimeDuration)

		m.LastCheckStatus = result.Status
		m.StatusChangedAt = &result.CheckedAt
	} else {
		if err := database.UpdateLastCheck(ctx, db, m.ID); err != nil {
			log.Printf("[ERROR] Failed to update last check for monitor %d.: %v", m.ID, err)
		}
	}
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
