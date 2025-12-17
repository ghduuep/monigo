package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildSMSHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	status := "UP"
	if res.Status == models.StatusDown {
		status = "CRITICAL"
	} else if res.Status == models.StatusDegraded {
		status = "DEGRADED"
	}

	msg := fmt.Sprintf("PINGLY: [%s] %s", status, m.Target)

	if res.Status != models.StatusUp {
		msg += fmt.Sprintf(" | Trace: %s", res.Message)
	} else {
		msg += fmt.Sprintf(" | %dms", res.Latency)
	}

	if inc != nil && inc.Duration != nil {
		msg += fmt.Sprintf(" | Dur: %s", inc.Duration.Round(time.Second))
	}

	return msg
}

func BuildSMSDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("PINGLY: [RESOLVED] %s (%s) matches config.", m.Target, dnsType)
}

func BuildSMSDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("PINGLY: [ALERT] %s (%s) mismatch. New: %s", m.Target, dnsType, res.ResultValue)
}

func BuildSMSDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("PINGLY: [FAIL] %s (%s). Err: %s", m.Target, dnsType, res.Message)
}

func BuildSMSPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	status := "OK"
	if res.Status == models.StatusDown {
		status = "FAIL"
	} else if res.Status == models.StatusDegraded {
		status = "SLOW"
	}

	msg := fmt.Sprintf("PINGLY: [TCP %s] %s", status, m.Target)

	if res.Status != models.StatusUp {
		msg += fmt.Sprintf(" | Err: %s", res.Message)
	} else {
		msg += fmt.Sprintf(" | %dms", res.Latency)
	}

	if inc != nil && inc.Duration != nil {
		msg += fmt.Sprintf(" | Dur: %s", inc.Duration.Round(time.Second))
	}

	return msg
}
