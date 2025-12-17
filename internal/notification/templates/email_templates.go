package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

const (
	colorSlate900 = "#0f172a"
	colorSlate500 = "#64748b"
	colorSlate400 = "#94a3b8"
	colorSlate100 = "#f1f5f9"
	colorWhite    = "#ffffff"
	colorGreen    = "#10b981"
	colorRed      = "#f43f5e"
	colorAmber    = "#f59e0b"
	colorBg       = "#f8fafc"

	fontFamily = `'Helvetica Neue', Helvetica, Arial, sans-serif`
)

func buildBaseEmail(title, badgeText, badgeColor, target, bodyContent string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	</head>
	<body style="margin: 0; padding: 0; background-color: %s; font-family: %s;">
		<div style="max-width: 600px; margin: 40px auto; background-color: %s; border-radius: 24px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1); border: 1px solid %s;">
			
			<div style="padding: 32px 32px 0 32px; display: flex; align-items: center; justify-content: space-between;">
				<div style="font-size: 20px; font-weight: 900; color: %s; letter-spacing: -0.05em;">Pingly</div>
			</div>

			<div style="padding: 32px;">
				<div style="margin-bottom: 24px;">
					<span style="background-color: %s15; color: %s; padding: 6px 12px; border-radius: 9999px; font-size: 11px; font-weight: 800; text-transform: uppercase; letter-spacing: 0.1em; display: inline-block; margin-bottom: 16px;">
						%s
					</span>
					<h1 style="margin: 0; font-size: 32px; font-weight: 900; color: %s; letter-spacing: -0.03em; line-height: 1.1;">
						%s
					</h1>
					<p style="margin: 8px 0 0 0; font-size: 14px; font-weight: 600; color: %s;">
						%s
					</p>
				</div>

				<div style="background-color: %s; border-radius: 16px; padding: 24px; margin-bottom: 24px;">
					%s
				</div>

				<div style="text-align: center; margin-top: 32px; padding-top: 24px; border-top: 1px solid %s;">
					<p style="margin: 0; font-size: 10px; font-weight: 800; text-transform: uppercase; letter-spacing: 0.2em; color: %s;">
						Infrastructure Monitoring
					</p>
				</div>
			</div>
		</div>
	</body>
	</html>
	`, colorBg, fontFamily, colorWhite, colorSlate100, colorSlate900, badgeColor, badgeColor, badgeText, colorSlate900, title, colorSlate500, target, colorSlate100, bodyContent, colorSlate100, colorSlate400)
}

func buildRow(label, value string, isMono bool) string {
	fontStack := fontFamily
	if isMono {
		fontStack = `'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace`
	}
	return fmt.Sprintf(`
	<div style="margin-bottom: 16px; last-child: margin-bottom: 0;">
		<p style="margin: 0 0 4px 0; font-size: 10px; font-weight: 900; text-transform: uppercase; letter-spacing: 0.1em; color: %s;">%s</p>
		<p style="margin: 0; font-size: 14px; font-weight: 600; color: %s; font-family: %s; word-break: break-all;">%s</p>
	</div>
	`, colorSlate400, label, colorSlate900, fontStack, value)
}

func BuildEmailHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var color, statusText, title string

	switch res.Status {
	case models.StatusDown:
		color = colorRed
		statusText = "CRITICAL OUTAGE"
		title = "Service Unreachable"
	case models.StatusDegraded:
		color = colorAmber
		statusText = "PERFORMANCE DEGRADED"
		title = "Latency Warning"
	default:
		color = colorGreen
		statusText = "OPERATIONAL"
		title = "Service Recovered"
	}

	content := buildRow("Response Time", fmt.Sprintf("%dms", res.Latency), true)
	if res.Message != "" {
		content += buildRow("Diagnostic Trace", res.Message, false)
	}

	if inc != nil {
		content += `<div style="margin-top: 24px; padding-top: 24px; border-top: 1px solid #cbd5e1;">`
		content += buildRow("Incident ID", fmt.Sprintf("#%d", inc.ID), true)
		content += buildRow("Started At", inc.StartedAt.Format("15:04:05 MST"), false)

		if inc.Duration != nil {
			content += buildRow("Total Downtime", inc.Duration.Round(time.Second).String(), false)
		}
		content += "</div>"
	}

	subject := fmt.Sprintf("[%s] %s: %s", statusText, title, m.Target)
	body := buildBaseEmail(title, statusText, color, m.Target, content)

	return subject, body
}

func BuildEmailDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string, inc *models.Incident) (string, string) {
	content := buildRow("Record Type", dnsType, true)
	content += buildRow("Resolved Value", res.ResultValue, true)

	if inc != nil && inc.Duration != nil {
		content += buildRow("Instability Duration", inc.Duration.Round(time.Second).String(), false)
	}

	subject := fmt.Sprintf("[RESOLVED] DNS Health: %s", m.Target)
	body := buildBaseEmail("DNS Integrity Restored", "RESOLVED", colorGreen, m.Target, content)

	return subject, body
}

func BuildEmailDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	content := buildRow("Record Type", dnsType, true)
	content += buildRow("Error Trace", res.Message, false)

	subject := fmt.Sprintf("[FAILURE] DNS Query Failed: %s", m.Target)
	body := buildBaseEmail("DNS Lookup Failed", "QUERY ERROR", colorAmber, m.Target, content)

	return subject, body
}

func BuildEmailDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	content := buildRow("Record Type", dnsType, true)
	content += buildRow("New Value Detected", res.ResultValue, true)
	content += buildRow("Alert Message", res.Message, false)

	subject := fmt.Sprintf("[ALERT] DNS Modified: %s", m.Target)
	body := buildBaseEmail("Record Mismatch Detected", "INTEGRITY ALERT", colorRed, m.Target, content)

	return subject, body
}

func BuildEmailPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var color, statusText, title string

	switch res.Status {
	case models.StatusDown:
		color = colorRed
		statusText = "CONNECTION FAILED"
		title = "Port Unreachable"
	case models.StatusDegraded:
		color = colorAmber
		statusText = "HIGH LATENCY"
		title = "Slow Connection"
	default:
		color = colorGreen
		statusText = "CONNECTED"
		title = "Port Accessible"
	}

	content := buildRow("Latency", fmt.Sprintf("%dms", res.Latency), true)

	if res.Message != "" {
		content += buildRow("Error Detail", res.Message, false)
	}

	if inc != nil && inc.Duration != nil {
		content += `<div style="margin-top: 24px; padding-top: 24px; border-top: 1px solid #cbd5e1;">`
		content += buildRow("Outage Duration", inc.Duration.Round(time.Second).String(), false)
		content += "</div>"
	}

	subject := fmt.Sprintf("[%s] TCP Alert: %s", statusText, m.Target)
	body := buildBaseEmail(title, statusText, color, m.Target, content)

	return subject, body
}
