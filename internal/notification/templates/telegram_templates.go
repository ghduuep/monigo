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
		statusLine = "*SERVIÇO FORA DO AR*"
	case models.StatusDegraded:
		emoji = "🟡"
		statusLine = "*PERFORMANCE DEGRADADA*"
	default:
		emoji = "🟢"
		statusLine = "*SERVIÇO OPERACIONAL*"
	}

	subject := fmt.Sprintf("%s Pingly Alert", emoji)

	body := fmt.Sprintf("%s\n\n", statusLine)
	body += fmt.Sprintf("🔗 *Alvo:* `%s`\n", m.Target)
	body += fmt.Sprintf("📡 *Status:* %s\n", res.Status)
	body += fmt.Sprintf("⚡ *Latência:* `%dms`\n", res.Latency)

	if res.Message != "" {
		body += fmt.Sprintf("📝 *Info:* _%s_\n", res.Message)
	}

	if inc != nil {
		body += "\n➖➖➖➖➖➖➖\n"
		body += fmt.Sprintf("🆔 *Incidente #%d*\n", inc.ID)
		if inc.Duration != nil {
			body += fmt.Sprintf("⏱ *Duração:* %s\n", inc.Duration.Round(time.Second))
		}
		if inc.ResolvedAt != nil {
			body += fmt.Sprintf("✅ *Resolvido em:* %s\n", inc.ResolvedAt.Format("15:04:05"))
		} else {
			body += fmt.Sprintf("🕒 *Início:* %s\n", inc.StartedAt.Format("15:04:05"))
		}
	}

	return subject, body
}

func BuildTelegramDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string, inc *models.Incident) (string, string) {
	subject := "🟢 Pingly DNS"
	body := fmt.Sprintf("✅ *DNS Resolvido*\n\n")
	body += fmt.Sprintf("🌍 *Alvo:* `%s`\n", m.Target)
	body += fmt.Sprintf("🏷 *Tipo:* %s\n", dnsType)
	body += fmt.Sprintf("🔢 *Valor:* `%s`\n", res.ResultValue)

	if inc != nil && inc.Duration != nil {
		body += fmt.Sprintf("\n⏱ *Instabilidade:* %s", inc.Duration.Round(time.Second))
	}
	return subject, body
}

func BuildTelegramDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := "🚨 Pingly DNS Alert"
	body := fmt.Sprintf("🚨 *ALTERAÇÃO DE DNS DETECTADA*\n\n")
	body += fmt.Sprintf("O registo %s para `%s` foi modificado.\n\n", dnsType, m.Target)
	body += fmt.Sprintf("🔻 *Novo Valor:*\n`%s`\n\n", res.ResultValue)
	body += fmt.Sprintf("⚠️ *Mensagem:* %s", res.Message)
	return subject, body
}

func BuildTelegramDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := "⚠️ Pingly DNS Warning"
	body := fmt.Sprintf("⚠️ *Falha na Consulta DNS*\n\n")
	body += fmt.Sprintf("🌍 *Alvo:* `%s` (%s)\n", m.Target, dnsType)
	body += fmt.Sprintf("❌ *Erro:* _%s_", res.Message)
	return subject, body
}

func BuildTelegramPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var emoji, statusLine string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusLine = "*FALHA DE CONEXÃO TCP*"
	} else if res.Status == models.StatusDegraded {
		emoji = "🟡"
		statusLine = "*CONEXÃO LENTA*"
	} else {
		emoji = "🟢"
		statusLine = "*CONEXÃO ESTABELECIDA*"
	}

	subject := fmt.Sprintf("%s Pingly TCP", emoji)

	body := fmt.Sprintf("%s\n\n", statusLine)
	body += fmt.Sprintf("🔌 *Host:* `%s`\n", m.Target)
	body += fmt.Sprintf("⚡ *Latência:* `%dms`\n", res.Latency)

	if res.Status != models.StatusUp {
		body += fmt.Sprintf("❌ *Erro:* _%s_\n", res.Message)
	}

	if inc != nil && inc.Duration != nil {
		body += fmt.Sprintf("\n⏱ *Duração:* %s", inc.Duration.Round(time.Second))
	}

	return subject, body
}
