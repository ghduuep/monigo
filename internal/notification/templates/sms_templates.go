package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildSMSHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	status := "UP"
	if res.Status == models.StatusDown {
		status = "DOWN"
	} else if res.Status == models.StatusDegraded {
		status = "SLOW"
	}

	msg := fmt.Sprintf("Pingly: [%s] %s", status, m.Target)

	if res.Status != models.StatusUp {
		msg += fmt.Sprintf(" - %s", res.Message)
	} else {
		msg += fmt.Sprintf(" - %dms", res.Latency)
	}

	if inc != nil && inc.Duration != nil {
		msg += fmt.Sprintf(" (Dur: %s)", inc.Duration.Round(time.Second))
	}

	return msg
}

func BuildSMSDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("Pingly: [DNS OK] %s (%s) matches config. Val: %s", m.Target, dnsType, res.ResultValue)
}

func BuildSMSDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("Pingly: [DNS CHANGE] %s (%s). New: %s. Msg: %s", m.Target, dnsType, res.ResultValue, res.Message)
}

func BuildSMSDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("Pingly: [DNS FAIL] %s (%s). Err: %s", m.Target, dnsType, res.Message)
}

func BuildSMSPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	status := "OK"
	if res.Status == models.StatusDown {
		status = "FAIL"
	} else if res.Status == models.StatusDegraded {
		status = "SLOW"
	}

	msg := fmt.Sprintf("Pingly: [TCP %s] %s", status, m.Target)

	if res.Status != models.StatusUp {
		msg += fmt.Sprintf(" - %s", res.Message)
	} else {
		msg += fmt.Sprintf(" - %dms", res.Latency)
	}

	if inc != nil && inc.Duration != nil {
		msg += fmt.Sprintf(" (Dur: %s)", inc.Duration.Round(time.Second))
	}

	return msg
}
