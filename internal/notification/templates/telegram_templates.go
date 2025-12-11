package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildTelegramHTTPMessage(m models.Monitor, res models.CheckResult, d time.Duration) (string, string) {
	var emoji, statusText string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		statusText = "DOWN"
	} else {
		emoji = "🟢"
		statusText = "UP"
	}

	subject := fmt.Sprintf("%s Monitor HTTP: %s", emoji, m.Target)

	body := fmt.Sprintf("\n\n📊 *Status:* %s", statusText)
	body += fmt.Sprintf("\n🔍 *Motivo:* %s", res.Message)
	body += fmt.Sprintf("\n⚡ *Latência:* %dms", res.Latency)

	timeLayout := "02/01 15:04:05"
	if res.Status == models.StatusDown {
		body += fmt.Sprintf("\n🕒 *Começou em:* %s", res.CheckedAt.Format(timeLayout))
	} else if res.Status == models.StatusUp && m.LastCheckStatus == models.StatusDown {
		body += fmt.Sprintf("\n🕒 *Resolvido em:* %s", res.CheckedAt.Format(timeLayout))
		if d > 0 {
			body += fmt.Sprintf("\n⏱ *Duração da Queda:* %s", d.Round(time.Second).String())
		}
	}

	return subject, body
}

func BuildTelegramDNSDetectedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("🟢 DNS %s Detectado: %s", dnsType, m.Target)

	body := fmt.Sprintf("\n\nA monitorização foi configurada com sucesso.")
	body += fmt.Sprintf("\n\n📄 *Valor Atual:* `%s`", res.ResultValue)
	body += "\n\n_Avisaremos se houver alterações._"

	return subject, body
}

func BuildTelegramDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("🚨 DNS %s Alterado: %s", dnsType, m.Target)

	body := "\n\n⚠️ *Atenção! O registo DNS mudou inesperadamente.*"
	body += fmt.Sprintf("\n\n🔻 *Valor Encontrado:* `%s`", res.ResultValue)
	body += fmt.Sprintf("\n💬 *Mensagem:* %s", res.Message)
	body += fmt.Sprintf("\n🕒 *Detectado em:* %s", res.CheckedAt.Format("15:04:05"))
	body += "\n\n_Verifique o seu domínio imediatamente._"

	return subject, body
}

func BuildTelegramDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("⚠️ Falha DNS %s: %s", dnsType, m.Target)

	body := "\n\nNão foi possível resolver o registo DNS."
	body += fmt.Sprintf("\n\n❌ *Erro:* %s", res.Message)
	body += fmt.Sprintf("\n📊 *Status:* %s", res.Status)

	return subject, body
}
