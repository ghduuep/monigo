package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildSMSHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	var emoji, statusText string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusText = "DOWN"
	} else if res.Status == models.StatusDegraded {
		emoji = "🟡"
		statusText = "SLOW"
	} else {
		emoji = "🟢"
		statusText = "UP"
	}

	msg := fmt.Sprintf("%s [HTTP] %s is %s.", emoji, m.Target, statusText)

	if res.Status == models.StatusDown {
		msg += fmt.Sprintf(" %s", res.Message)
	} else if res.Status == models.StatusDegraded {
		msg += fmt.Sprintf(" %s", res.Message)
	} else {
		msg += fmt.Sprintf(" %dms.", res.Latency)

		if inc != nil {
			if m.LastCheckStatus == models.StatusDown && inc.Duration != nil {
				msg += fmt.Sprintf(" Down for: %s", inc.Duration.Round(time.Second).String())
			}
			if m.LastCheckStatus == models.StatusDegraded && inc.Duration != nil {
				msg += fmt.Sprintf(" Slow for: %s", inc.Duration.Round(time.Second).String())
			}
		}
	}

	return msg
}

func BuildSMSDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("Pingly DNS OK: %s record for %s matches config. Val: %s", dnsType, m.Target, res.ResultValue)
}

func BuildSMSDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("🚨 [DNS %s] %s CHANGED! New: %s. Msg: %s", dnsType, m.Target, res.ResultValue, res.Message)
}

func BuildSMSDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("⚠️ [DNS %s] %s Failed. Err: %s", dnsType, m.Target, res.Message)
}

func BuildSMSPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) string {
	var emoji, statusText string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusText = "FAIL"
	} else if res.Status == models.StatusDegraded {
		emoji = "🟡"
		statusText = "SLOW"
	} else {
		emoji = "🟢"
		statusText = "OK"
	}

	msg := fmt.Sprintf("%s [TCP] %s is %s.", emoji, m.Target, statusText)

	if res.Status == models.StatusDown {
		msg += fmt.Sprintf(" %s", res.Message)
	} else if res.Status == models.StatusDegraded {
		msg += fmt.Sprintf(" %s", res.Message)
	} else {
		msg += fmt.Sprintf(" %dms.", res.Latency)
		if inc != nil {
			if m.LastCheckStatus == models.StatusDown && inc.Duration != nil {
				msg += fmt.Sprintf(" Down for: %s", inc.Duration.Round(time.Second).String())
			}
			if m.LastCheckStatus == models.StatusDegraded && inc.Duration != nil {
				msg += fmt.Sprintf(" Slow for: %s", inc.Duration.Round(time.Second).String())
			}
		}
	}

	return msg
}
