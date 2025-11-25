package monitor

import (
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func checkHTTP(m *models.Monitor) models.CheckResult {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	start := time.Now()

	resp, err := client.Get(m.Target)
	latency := time.Since(start)

	if err != nil {
		return models.CheckResult{
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Latency:   latency,
			Message:   err.Error(),
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	status := models.StatusDown
	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		status = models.StatusUp
	}

	return models.CheckResult{
		MonitorID:  m.ID,
		Status:     status,
		Latency:    latency,
		StatusCode: resp.StatusCode,
		CheckedAt:  time.Now(),
	}
}
