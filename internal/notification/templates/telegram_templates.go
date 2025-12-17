package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildTelegramHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var emoji, statusLine string

	switch res.Status {
	case models.StatusDown:
		emoji = "🔴"
		statusLine = "CRITICAL OUTAGE"
	case models.StatusDegraded:
		emoji = "🟡"
		statusLine = "PERFORMANCE DEGRADED"
	default:
		emoji = "🟢"
		statusLine = "OPERATIONAL"
	}

	subject := fmt.Sprintf("%s Pingly Alert", emoji)

	body := fmt.Sprintf("*%s*\n\n", statusLine)
	body += fmt.Sprintf("📡 *TARGET RESOURCE*\n`%s`\n\n", m.Target)
	body += fmt.Sprintf("⚡ *LATENCY*\n`%dms`\n\n", res.Latency)

	if res.Message != "" {
		body += fmt.Sprintf("📝 *DIAGNOSTIC TRACE*\n_%s_\n\n", res.Message)
	}

	if inc != nil {
		body += "➖➖➖➖➖➖➖➖➖\n"
		body += fmt.Sprintf("🆔 *INCIDENT #%d*\n", inc.ID)
		if inc.Duration != nil {
			body += fmt.Sprintf("⏱ *TOTAL DURATION*: `%s`\n", inc.Duration.Round(time.Second))
		}
		if inc.ResolvedAt != nil {
			body += fmt.Sprintf("✅ *RESOLVED*: `%s`", inc.ResolvedAt.Format("15:04:05"))
		} else {
			body += fmt.Sprintf("🕒 *STARTED*: `%s`", inc.StartedAt.Format("15:04:05"))
		}
	}

	return subject, body
}

func BuildTelegramDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string, inc *models.Incident) (string, string) {
	subject := "🟢 Pingly DNS"
	body := "*DNS INTEGRITY RESTORED*\n\n"
	body += fmt.Sprintf("🌍 *TARGET*: `%s`\n", m.Target)
	body += fmt.Sprintf("🏷 *RECORD*: `%s`\n", dnsType)
	body += fmt.Sprintf("🔢 *VALUE*: `%s`\n", res.ResultValue)

	if inc != nil && inc.Duration != nil {
		body += fmt.Sprintf("\n⏱ *INSTABILITY*: `%s`", inc.Duration.Round(time.Second))
	}
	return subject, body
}

func BuildTelegramDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := "🚨 Pingly DNS Alert"
	body := "*RECORD MISMATCH DETECTED*\n\n"
	body += fmt.Sprintf("Target: `%s` (%s)\n\n", m.Target, dnsType)
	body += "*NEW VALUE DETECTED*\n"
	body += fmt.Sprintf("`%s`\n\n", res.ResultValue)
	body += fmt.Sprintf("⚠️ *TRACE*: _%s_", res.Message)
	return subject, body
}

func BuildTelegramDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := "⚠️ Pingly DNS Warning"
	body := "*QUERY FAILURE*\n\n"
	body += fmt.Sprintf("🌍 *TARGET*: `%s` (%s)\n", m.Target, dnsType)
	body += fmt.Sprintf("❌ *TRACE*: _%s_", res.Message)
	return subject, body
}

func BuildTelegramPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var emoji, statusLine string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusLine = "CONNECTION FAILED"
	} else if res.Status == models.StatusDegraded {
		emoji = "🟡"
		statusLine = "HIGH LATENCY"
	} else {
		emoji = "🟢"
		statusLine = "CONNECTED"
	}

	subject := fmt.Sprintf("%s Pingly TCP", emoji)

	body := fmt.Sprintf("*%s*\n\n", statusLine)
	body += fmt.Sprintf("🔌 *HOST*: `%s`\n", m.Target)
	body += fmt.Sprintf("⚡ *LATENCY*: `%dms`\n", res.Latency)

	if res.Status != models.StatusUp {
		body += fmt.Sprintf("\n❌ *TRACE*: _%s_\n", res.Message)
	}

	if inc != nil && inc.Duration != nil {
		body += fmt.Sprintf("\n⏱ *DURATION*: `%s`", inc.Duration.Round(time.Second))
	}

	return subject, body
}
