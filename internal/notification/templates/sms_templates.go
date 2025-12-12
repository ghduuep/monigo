package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildSMSHTTPMessage(m models.Monitor, res models.CheckResult, d time.Duration) string {
	var emoji, statusText string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusText = "DOWN"
	} else {
		emoji = "🟢"
		statusText = "UP"
	}

	msg := fmt.Sprintf("%s [HTTP] %s is %s.", emoji, m.Target, statusText)

	if res.Status == models.StatusDown {
		msg += fmt.Sprintf(" %s", res.Message)
	} else {
		msg += fmt.Sprintf(" %dms.", res.Latency)

		if m.LastCheckStatus == models.StatusDown && d > 0 {
			msg += fmt.Sprintf(" Down for: %s", d.Round(time.Second).String())
		}
	}

	return msg
}

func BuildSMSDNSDetectedMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("🟢 [DNS %s] %s Configured. Val: %s", dnsType, m.Target, res.ResultValue)
}

func BuildSMSDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("🚨 [DNS %s] %s CHANGED! New: %s. Msg: %s", dnsType, m.Target, res.ResultValue, res.Message)
}

func BuildSMSDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) string {
	return fmt.Sprintf("⚠️ [DNS %s] %s Failed. Err: %s", dnsType, m.Target, res.Message)
}

func BuildSMSPortMessage(m models.Monitor, res models.CheckResult, d time.Duration) string {
	var emoji, statusText string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusText = "FAIL"
	} else {
		emoji = "🟢"
		statusText = "OK"
	}

	msg := fmt.Sprintf("%s [TCP] %s is %s.", emoji, m.Target, statusText)

	if res.Status == models.StatusDown {
		msg += fmt.Sprintf(" %s", res.Message)
	} else {
		msg += fmt.Sprintf(" %dms.", res.Latency)
		if m.LastCheckStatus == models.StatusDown && d > 0 {
			msg += fmt.Sprintf(" Down for: %s", d.Round(time.Second).String())
		}
	}

	return msg
}
