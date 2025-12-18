package monitor

import (
	"context"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/ghduuep/pingly/internal/notification"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"reflect"
	"time"
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
	return old.Target != new.Target || old.Interval != new.Interval || old.Timeout != new.Timeout || old.LatencyThreshold != new.LatencyThreshold || !reflect.DeepEqual(old.Config, new.Config)
}
