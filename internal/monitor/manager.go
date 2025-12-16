package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	SyncInterval      = 10 * time.Second
	DownCheckInterval = 30 * time.Second
	FailureThreshold  = 2
	FlappingTTLMulti  = 3
)

type activeMonitor struct {
	cancel context.CancelFunc
	config models.Monitor
}

type MonitorManager struct {
	db             *pgxpool.Pool
	redis          *redis.Client
	dispatcher     notification.NotificationDispatcher
	activeMonitors map[int]*activeMonitor
}

func NewMonitorManager(db *pgxpool.Pool, rdb *redis.Client, dispatcher notification.NotificationDispatcher) *MonitorManager {
	return &MonitorManager{
		db:             db,
		redis:          rdb,
		dispatcher:     dispatcher,
		activeMonitors: make(map[int]*activeMonitor),
	}
}

func (m *MonitorManager) Start(ctx context.Context) {
	log.Println("[INFO] Monitor Manager stated...")

	ticker := time.NewTicker(SyncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			m.stopAll()
			return
		case <-ticker.C:
			m.syncMonitors(ctx)
		}
	}
}

func (m *MonitorManager) syncMonitors(ctx context.Context) {
	monitors, err := database.GetAllMonitors(ctx, m.db)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch monitors: %v", err)
		return
	}

	currentIDs := make(map[int]bool)

	for _, mon := range monitors {
		currentIDs[mon.ID] = true
		active, exists := m.activeMonitors[mon.ID]

		if !exists {
			m.startMonitor(ctx, *mon)
		} else if m.hasChanged(active.config, *mon) {
			log.Printf("[INFO] Configuration changed for monitor %d", mon.ID)
			active.cancel()
			m.startMonitor(ctx, *mon)
		}
	}

	for id, active := range m.activeMonitors {
		if !currentIDs[id] {
			log.Printf("[INFO] Stopping monitor %d (removed)", id)
			active.cancel()
			delete(m.activeMonitors, id)
		}
	}
}

func (m *MonitorManager) startMonitor(ctx context.Context, mon models.Monitor) {
	monCtx, cancel := context.WithCancel(ctx)

	m.activeMonitors[mon.ID] = &activeMonitor{
		cancel: cancel,
		config: mon,
	}

	go m.runWorker(monCtx, mon)
	log.Printf("[INFO] Started monitoring for %s (%s)", mon.Target, mon.Type)
}

func (m *MonitorManager) runWorker(ctx context.Context, mon models.Monitor) {
	log.Printf("[INFO] Perfoming initial check for monitor %d", mon.ID)
	initialStatus := m.processCheck(ctx, &mon)

	initialDelay := mon.Interval
	if initialStatus == models.StatusDown || initialStatus == models.StatusDegraded {
		initialDelay = DownCheckInterval
	}

	timer := time.NewTimer(initialDelay)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			status := m.processCheck(ctx, &mon)

			nextInterval := mon.Interval
			if status == models.StatusDown || status == models.StatusDegraded {
				nextInterval = DownCheckInterval
			}

			timer.Reset(nextInterval)
		}
	}
}

func (m *MonitorManager) processCheck(ctx context.Context, mon *models.Monitor) models.MonitorStatus {
	result := performCheck(*mon)

	m.handleDNSLearning(ctx, mon, &result)

	shouldProceed := m.isConfirmedFailure(ctx, mon, result.Status)
	if !shouldProceed {
		_ = database.UpdateLastCheck(ctx, m.db, mon.ID)
		return mon.LastCheckStatus
	}

	m.analyzePerformance(ctx, mon, &result)

	if err := database.CreateCheckResult(ctx, m.db, &result); err != nil {
		log.Printf("[ERROR] Failed to save check result for monitor %d", mon.ID)
	}

	if result.Status != mon.LastCheckStatus && result.Status != models.StatusUnknown {
		m.handleStateChange(ctx, mon, result)
	} else {
		_ = database.UpdateLastCheck(ctx, m.db, mon.ID)
	}

	return result.Status
}

func (m *MonitorManager) isConfirmedFailure(ctx context.Context, mon *models.Monitor, currentStatus models.MonitorStatus) bool {
	if currentStatus == models.StatusUp {
		key := fmt.Sprintf("monitor:%d:fails", mon.ID)
		m.redis.Del(ctx, key)
		return true
	}

	if mon.LastCheckStatus == models.StatusDown || mon.LastCheckStatus == models.StatusDegraded {
		return true
	}

	key := fmt.Sprintf("monitor:%d:fails", mon.ID)
	count, err := m.redis.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("[ERROR] Redis error on flapping check: %v", err)
		return true
	}

	m.redis.Expire(ctx, key, mon.Interval*time.Duration(FlappingTTLMulti))

	if count < FailureThreshold {
		log.Printf("[INFO] Monitor %d flapping detected (%d/%d). Suppressing alert", mon.ID, count, FailureThreshold)
		return false
	}

	return true
}

func (m *MonitorManager) handleDNSLearning(ctx context.Context, mon *models.Monitor, res *models.CheckResult) {
	if mon.Type != models.TypeDNS || res.Status != models.StatusUp || res.ResultValue != "" {
		return
	}

	var config models.DNSConfig
	if err := json.Unmarshal(mon.Config, &config); err == nil {
		log.Printf("[INFO] Learning DNS value for monitor %d: %s", mon.ID, res.ResultValue)

		if err := database.SetInitialDNSConfig(ctx, m.db, mon.ID, res.ResultValue); err != nil {
			log.Printf("[ERROR] Failed to persist learned DNS config: %v", err)
		}

		config.ExpectedValue = res.ResultValue
		newJSON, _ := json.Marshal(config)
		mon.Config = newJSON
	}
}

func (m *MonitorManager) analyzePerformance(ctx context.Context, mon *models.Monitor, res *models.CheckResult) {
	if res.Status != models.StatusUp {
		return
	}

	if mon.LatencyThreshold > 0 && res.Latency > mon.LatencyThreshold {
		log.Printf("[WARN] Latency threshold exceeded for monitor %d (%dms > %dms)", mon.ID, res.Latency, mon.LatencyThreshold)
		res.Status = models.StatusDegraded
		res.Message = fmt.Sprintf("Low performance: %dms (Limit: %dms)", res.Latency, mon.LatencyThreshold)
		return
	}

	history, err := database.GetRecentLatencies(ctx, m.db, mon.ID, 30)
	if err == nil {
		isAnom, msg := isAnomaly(res.Latency, history)
		if isAnom {
			log.Printf("[INFO] Anomaly detected for monitor %d", mon.ID)
			res.Status = models.StatusDegraded
			res.Message = msg
		}
	}
}

func (m *MonitorManager) handleStateChange(ctx context.Context, mon *models.Monitor, res models.CheckResult) {
	log.Printf("[INFO] Monitor %d state change: %s -> %s", mon.ID, mon.LastCheckStatus, res.Status)

	var incident *models.Incident
	var err error

	if res.Status == models.StatusDown || res.Status == models.StatusDegraded {
		incident, err = database.CreateIncident(ctx, m.db, mon.ID, res.Message)
		if err != nil {
			log.Printf("[ERROR] Failed to create incident: %v", err)
		}
	}

	wasBad := mon.LastCheckStatus == models.StatusDown || mon.LastCheckStatus == models.StatusDegraded
	isRecovered := res.Status == models.StatusUp

	if wasBad && isRecovered {
		incident, err = database.ResolveIncident(ctx, m.db, mon.ID)
		if err != nil {
			log.Printf("[ERROR] Failed to resolve incident: %v", err)
		}
	}

	if err := database.UpdateMonitorStatus(ctx, m.db, mon.ID, string(res.Status)); err != nil {
		log.Printf("[ERROR] Failed to update monitor status: %v", err)
	}

	mon.LastCheckStatus = res.Status
	mon.StatusChangedAt = &res.CheckedAt

	if incident != nil {
		channels, _ := database.GetEnabledUserChannels(ctx, m.db, mon.UserID)
		go m.dispatcher.SendAlert(channels, *mon, res, incident)
	}
}

func (m *MonitorManager) stopAll() {
	for _, active := range m.activeMonitors {
		active.cancel()
	}
	log.Printf("All monitors stopped.")
}

func (m *MonitorManager) hasChanged(old, new models.Monitor) bool {
	return old.Target != new.Target || old.Interval != new.Interval || old.Timeout != new.Timeout || !reflect.DeepEqual(old.Config, new.Config)
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
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Message:   "Unknown monitor type",
			CheckedAt: time.Now(),
		}
	}
}
