package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"log"
	"time"
)

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
