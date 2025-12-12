package monitor

import (
	"fmt"
	"github.com/ghduuep/pingly/internal/models"
	"net"
	"strings"
	"time"
)

func checkPort(m models.Monitor) models.CheckResult {
	target := m.Target
	if !strings.Contains(target, ":") {
		target = fmt.Sprintf("%s:443", target)
	}

	timeout := m.Timeout
	start := time.Now()

	conn, err := net.DialTimeout("tcp", target, timeout)

	latency := time.Since(start).Milliseconds()

	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "timeout") {
			msg = "Timeout: Port unreachable or firewall blocking"
		} else if strings.Contains(msg, "refused") {
			msg = "Connection Refused: Server is up, but service is down"
		} else if strings.Contains(msg, "no such host") {
			msg = "DNS error: Domain not found"
		}

		return models.CheckResult{
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Latency:   0,
			Message:   msg,
			CheckedAt: time.Now(),
		}
	}

	conn.Close()

	return models.CheckResult{
		MonitorID:   m.ID,
		Status:      models.StatusUp,
		Latency:     latency,
		ResultValue: fmt.Sprintf("%dms", latency),
		Message:     "Port acessible (TCP handshake)",
		CheckedAt:   time.Now(),
	}
}
