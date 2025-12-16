package templates

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

// Cores e Estilos comuns
const (
	colorGreen    = "#22c55e" // Success
	colorRed      = "#ef4444" // Error
	colorYellow   = "#eab308" // Warning
	bgGray        = "#f3f4f6"
	containerBody = `font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);`
)

func BuildEmailHTTPMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var color, title, emoji string
	timeLayout := "02/01/2006 às 15:04:05"

	// Definição de Estados
	switch res.Status {
	case models.StatusDown:
		color = colorRed
		emoji = "🔴"
		title = "Serviço Indisponível (DOWN)"
	case models.StatusDegraded:
		color = colorYellow
		emoji = "🟡"
		title = "Performance Degradada"
	default:
		color = colorGreen
		emoji = "🟢"
		title = "Serviço Recuperado (UP)"
	}

	// Detalhes do Incidente
	var detailsHTML string
	if inc != nil {
		durationStr := "Calculando..."
		if inc.Duration != nil {
			durationStr = inc.Duration.Round(time.Second).String()
		} else if res.Status == models.StatusUp && inc.ResolvedAt != nil {
			// Se acabou de resolver, calcula a duração final
			d := inc.ResolvedAt.Sub(inc.StartedAt)
			durationStr = d.Round(time.Second).String()
		}

		detailsHTML += fmt.Sprintf(`
			<div style="background-color: #f9fafb; padding: 15px; border-radius: 6px; margin-top: 15px; border: 1px solid #e5e7eb;">
				<p style="margin: 5px 0; color: #374151; font-size: 14px;"><strong>🆔 Incidente:</strong> #%d</p>
				<p style="margin: 5px 0; color: #374151; font-size: 14px;"><strong>🕒 Início:</strong> %s</p>
				<p style="margin: 5px 0; color: #374151; font-size: 14px;"><strong>⏱ Duração:</strong> %s</p>
			</div>
		`, inc.ID, inc.StartedAt.Format(timeLayout), durationStr)
	}

	subject := fmt.Sprintf("%s [%s] %s", emoji, title, m.Target)

	body := fmt.Sprintf(`
		<div style="background-color: %s; padding: 40px 20px;">
			<div style="%s">
				<div style="background-color: %s; padding: 20px; text-align: center;">
					<h1 style="color: white; margin: 0; font-size: 24px; font-weight: bold;">%s</h1>
				</div>
				<div style="padding: 30px;">
					<p style="font-size: 16px; color: #4b5563; margin-top: 0;">O monitor <strong>%s</strong> reportou um novo status.</p>
					
					<table style="width: 100%%; margin-top: 20px; border-collapse: collapse;">
						<tr>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6;"><strong>🎯 Alvo</strong></td>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6; text-align: right; color: #6b7280;">%s</td>
						</tr>
						<tr>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6;"><strong>📊 Status</strong></td>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6; text-align: right; color: %s; font-weight: bold;">%s</td>
						</tr>
						<tr>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6;"><strong>⚡ Latência</strong></td>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6; text-align: right; color: #6b7280;">%dms</td>
						</tr>
						<tr>
							<td style="padding: 10px;"><strong>💬 Mensagem</strong></td>
							<td style="padding: 10px; text-align: right; color: #6b7280;">%s</td>
						</tr>
					</table>

					%s

					<div style="margin-top: 30px; text-align: center; color: #9ca3af; font-size: 12px;">
						Enviado automaticamente pelo Pingly 📡
					</div>
				</div>
			</div>
		</div>
	`, bgGray, containerBody, color, title, m.Target, m.Target, color, res.Status, res.Latency, res.Message, detailsHTML)

	return subject, body
}

func BuildEmailDNSRecoveredMessage(m models.Monitor, res models.CheckResult, dnsType string, inc *models.Incident) (string, string) {
	subject := fmt.Sprintf("🟢 DNS Recuperado: %s (%s)", m.Target, dnsType)

	detailsHTML := ""
	if inc != nil && inc.Duration != nil {
		detailsHTML = fmt.Sprintf(`<p style="color: #374151;"><strong>Tempo de Instabilidade:</strong> %s</p>`, inc.Duration.Round(time.Second).String())
	}

	body := fmt.Sprintf(`
		<div style="background-color: %s; padding: 40px 20px;">
			<div style="%s">
				<div style="background-color: %s; padding: 20px; text-align: center;">
					<h1 style="color: white; margin: 0; font-size: 24px;">DNS Resolvido ✅</h1>
				</div>
				<div style="padding: 30px;">
					<p style="color: #4b5563;">O registo <strong>%s</strong> para <strong>%s</strong> está correto novamente.</p>
					
					<div style="background-color: #f0fdf4; border-left: 4px solid %s; padding: 15px; margin: 20px 0;">
						<p style="margin: 0; color: #166534; font-family: monospace; font-size: 14px;">%s</p>
					</div>
					%s
				</div>
			</div>
		</div>
	`, bgGray, containerBody, colorGreen, dnsType, m.Target, colorGreen, res.ResultValue, detailsHTML)

	return subject, body
}

func BuildEmailDNSStatusMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("⚠️ Falha na Consulta DNS: %s", m.Target)

	body := fmt.Sprintf(`
		<div style="background-color: %s; padding: 40px 20px;">
			<div style="%s">
				<div style="background-color: %s; padding: 20px; text-align: center;">
					<h1 style="color: white; margin: 0; font-size: 24px;">Falha DNS ⚠️</h1>
				</div>
				<div style="padding: 30px;">
					<p style="color: #4b5563;">Não foi possível validar o registo <strong>%s</strong> para <strong>%s</strong>.</p>
					
					<div style="background-color: #fffbeb; border: 1px solid %s; padding: 15px; border-radius: 4px; margin-top: 20px;">
						<strong style="color: #92400e;">Erro Técnico:</strong>
						<p style="margin: 5px 0; color: #b45309;">%s</p>
					</div>
				</div>
			</div>
		</div>
	`, bgGray, containerBody, colorYellow, dnsType, m.Target, colorYellow, res.Message)

	return subject, body
}

func BuildEmailDNSChangedMessage(m models.Monitor, res models.CheckResult, dnsType string) (string, string) {
	subject := fmt.Sprintf("🚨 DNS Alterado: %s (%s)", m.Target, dnsType)

	body := fmt.Sprintf(`
		<div style="background-color: %s; padding: 40px 20px;">
			<div style="%s">
				<div style="background-color: %s; padding: 20px; text-align: center;">
					<h1 style="color: white; margin: 0; font-size: 24px;">Alteração Crítica 🚨</h1>
				</div>
				<div style="padding: 30px;">
					<p style="color: #4b5563;">O registo DNS <strong>%s</strong> do alvo <strong>%s</strong> não corresponde ao esperado.</p>
					
					<div style="margin-top: 20px;">
						<p style="font-size: 12px; color: #6b7280; text-transform: uppercase; margin-bottom: 5px;">Novo Valor Detectado:</p>
						<div style="background-color: #fee2e2; border-left: 4px solid %s; padding: 15px;">
							<code style="color: #991b1b; word-break: break-all;">%s</code>
						</div>
					</div>

					<div style="margin-top: 20px; padding: 10px; border-top: 1px solid #e5e7eb;">
						<p style="color: #6b7280; font-size: 14px;"><strong>Mensagem:</strong> %s</p>
					</div>
				</div>
			</div>
		</div>
	`, bgGray, containerBody, colorRed, dnsType, m.Target, colorRed, res.ResultValue, res.Message)

	return subject, body
}

func BuildEmailPortMessage(m models.Monitor, res models.CheckResult, inc *models.Incident) (string, string) {
	var color, title, emoji string
	switch res.Status {
	case models.StatusDown:
		color = colorRed
		emoji = "🔴"
		title = "Falha de Conexão (TCP)"
	case models.StatusDegraded:
		color = colorYellow
		emoji = "🟡"
		title = "Conexão Lenta (TCP)"
	default:
		color = colorGreen
		emoji = "🟢"
		title = "Conexão Estabelecida"
	}

	subject := fmt.Sprintf("%s [%s] %s", emoji, title, m.Target)

	body := fmt.Sprintf(`
		<div style="background-color: %s; padding: 40px 20px;">
			<div style="%s">
				<div style="background-color: %s; padding: 20px; text-align: center;">
					<h1 style="color: white; margin: 0; font-size: 24px;">%s</h1>
				</div>
				<div style="padding: 30px;">
					<p style="color: #4b5563;">Status da porta TCP monitorizada.</p>
					
					<table style="width: 100%%; margin-top: 20px; border-collapse: collapse;">
						<tr>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6;"><strong>Target</strong></td>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6; text-align: right;">%s</td>
						</tr>
						<tr>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6;"><strong>Latência</strong></td>
							<td style="padding: 10px; border-bottom: 1px solid #f3f4f6; text-align: right;">%dms</td>
						</tr>
						<tr>
							<td style="padding: 10px;"><strong>Resposta</strong></td>
							<td style="padding: 10px; text-align: right;">%s</td>
						</tr>
					</table>
				</div>
			</div>
		</div>
	`, bgGray, containerBody, color, title, m.Target, res.Latency, res.Message)

	return subject, body
}
