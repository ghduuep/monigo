package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func BuildEmailHTTPMessage(m models.Monitor, res models.CheckResult, d time.Duration) (string, string) {
	emoji := "🟢"
	color := "#38a169" // Verde

	timeLayout := "02/01/2006 15:04:05"
	var timeDetails string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		color = "#e53e5e" // Vermelho
		timeDetails += fmt.Sprintf("<p><strong>Começou em:</strong> %s</p>", res.CheckedAt.Format(timeLayout))
	} else if res.Status == models.StatusUp && m.LastCheckStatus == models.StatusDown {
		timeDetails += fmt.Sprintf("<p><strong>Começou em:</strong> %s</p>", m.StatusChangedAt.Format(timeLayout))
		timeDetails += fmt.Sprintf("<p><strong>Resolvido em:</strong> %s</p>", res.CheckedAt.Format(timeLayout))

		if d > 0 {
			timeDetails += fmt.Sprintf("<p><strong>Duração:</strong> %s</p>", d.Round(time.Second).String())
		}
	}

	subject := fmt.Sprintf("%s %s está %s", emoji, m.Target, res.Status)

	body := fmt.Sprintf(`
		<h2>Atualização de Status HTTP</h2>
		<p>O monitor <strong>%s</strong> mudou para <span style="color:%s"><strong>%s</strong></span>.</p>
		<p><strong>Motivo:</strong> %s</p>
		<p><strong>Latência:</strong> %vms</p>
		%s
	`, m.Target, color, res.Status, res.Message, res.Latency, timeDetails)

	return subject, body
}

func BuildEmailDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("⚠️ Falha de DNS Tipo %s: %s", dnsType, m.Target)

	body := fmt.Sprintf(`
		<h2>Problema de Resolução DNS Tipo %s</h2>
		<p>Não foi possível verificar os registos DNS para <strong>%s</strong>.</p>
		<p><strong>Status:</strong> %s</p>
		<p><strong>Erro Técnico:</strong> %s</p>
		<p><em>Verifique se o domínio expirou ou se os servidores de nome estão a responder.</em></p>
	`, dnsType, m.Target, res.Status, res.Message)

	return subject, body
}

func BuildEmailDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("🚨 DNS tipo %s de %s foi Alterado!", dnsType, m.Target)

	body := fmt.Sprintf(`
		<div style="border: 2px solid red; padding: 15px; background-color: #fff5f5;">
			<h2 style="color: red;">Alteração de Registo Detectada</h2>
			<p>O registo DNS Tipo %s monitorizado não corresponde à configuração esperada.</p>
			
			<ul>
				<li><strong>Alvo:</strong> %s</li>
				<li><strong>Valor Encontrado (Atual):</strong> <code>%s</code></li>
				<li><strong>Mensagem do Sistema:</strong> %s</li>
				<li><strong>Detectado em: %s</strong></li>
			</ul>

			<p><strong>Ação Recomendada:</strong> Verifique imediatamente se o seu domínio foi comprometido ou se houve uma atualização não planeada.</p>
		</div>
	`, dnsType, m.Target, res.ResultValue, res.Message, res.CheckedAt.Format("02/01/2006 15:04:05"))

	return subject, body
}

func BuildEmailDNSDetectedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("🟢 DNS Tipo %s Detectado: %s", dnsType, m.Target)

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; border: 1px solid #38a169; border-radius: 5px; background-color: #f0fff4;">
			<h2 style="color: #38a169;">Monitorização DNS Tipo %s Ativa</h2>
			<p>A monitorização para <strong>%s</strong> foi atualizada com sucesso.</p>
			
			<ul>
				<li><strong>Status:</strong> <span style="color: #38a169; font-weight: bold;">UP (Ativo)</span></li>
				<li><strong>Valor Detectado:</strong> <code>%s</code></li>
			</ul>

			<p style="font-size: 12px; color: #666;">A partir de agora, avisaremos se esse valor mudar.</p>
		</div>
	`, dnsType, m.Target, res.ResultValue)

	return subject, body
}

func BuildEmailPingMessage(m models.Monitor, res models.CheckResult, d time.Duration) (string, string) {
	emoji := "🟢"
	color := "#38a169" // Verde
	statusText := "CONECTADO"

	timeLayout := "02/01/2006 15:04:05"
	var timeDetails string

	if res.Status == models.StatusDown {
		emoji = "🔴"
		color = "#e53e5e" // Vermelho
		statusText = "FALHA"
		timeDetails += fmt.Sprintf("<p><strong>Começou em:</strong> %s</p>", res.CheckedAt.Format(timeLayout))
	} else if res.Status == models.StatusUp && m.LastCheckStatus == models.StatusDown {
		timeDetails += fmt.Sprintf("<p><strong>Começou em:</strong> %s</p>", m.StatusChangedAt.Format(timeLayout))
		timeDetails += fmt.Sprintf("<p><strong>Resolvido em:</strong> %s</p>", res.CheckedAt.Format(timeLayout))

		if d > 0 {
			timeDetails += fmt.Sprintf("<p><strong>Tempo de inatividade:</strong> %s</p>", d.Round(time.Second).String())
		}
	}

	subject := fmt.Sprintf("%s Ping/TCP %s: %s", emoji, m.Target, statusText)

	body := fmt.Sprintf(`
		<h2>Atualização de Conectividade (Ping/TCP)</h2>
		<p>O monitor <strong>%s</strong> reportou o estado: <span style="color:%s"><strong>%s</strong></span>.</p>
		
		<div style="border-left: 4px solid %s; padding-left: 10px; margin: 15px 0;">
			<p><strong>Target:</strong> %s</p>
			<p><strong>Mensagem:</strong> %s</p>
			<p><strong>Latência:</strong> %vms</p>
		</div>

		%s
		<p style="font-size: 12px; color: #666;">Verificação realizada via TCP Handshake.</p>
	`, m.Target, color, statusText, color, m.Target, res.Message, res.Latency, timeDetails)

	return subject, body
}
