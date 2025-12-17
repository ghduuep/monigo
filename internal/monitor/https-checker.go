package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func checkHTTP(m models.Monitor) models.CheckResult {
	var config models.HTTPConfig
	if len(m.Config) > 0 {
		_ = json.Unmarshal(m.Config, &config)
	}

	client := http.Client{
		Timeout: m.Timeout,
	}

	start := time.Now()

	resp, err := client.Get(m.Target)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		message := err.Error()

		if strings.Contains(message, "deadline exceeded") {
			message = "Connection Timeout"
		}

		return models.CheckResult{
			MonitorID: m.ID,
			Status:    models.StatusDown,
			Latency:   latency,
			Message:   message,
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	message := resp.Status
	status := models.StatusDown
	if resp.StatusCode >= 200 && resp.StatusCode < 500 {
		status = models.StatusUp
	}

	var resultValue string

	if config.CheckSSL && resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		expiresIn := time.Until(cert.NotAfter)
		days := int(expiresIn.Hours() / 24)

		resultValue = strconv.Itoa(days)

		if expiresIn < 0 {
			status = models.StatusDown
			message = fmt.Sprintf("CRITICAL: SSL certificate expired %s (%d days ago)", cert.NotAfter.Format("02/01/2006"), -days)
		} else if days <= 30 {
			status = models.StatusDegraded
			message = fmt.Sprintf("SSl expires in %d days (%s)", days, cert.NotAfter.Format("02/01/2006"))
		}
	}

	return models.CheckResult{
		MonitorID:   m.ID,
		Status:      status,
		Latency:     latency,
		Message:     message,
		StatusCode:  resp.StatusCode,
		ResultValue: resultValue,
		CheckedAt:   time.Now(),
	}
}
