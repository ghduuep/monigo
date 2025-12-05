package notification

import (
	"fmt"
	"time"

	"github.com/ghduuep/pingly/internal/models"
)

func (s *EmailService) SendStatusAlert(userEmail string, m models.Monitor, result models.CheckResult, dnsType string, duration time.Duration) error {

	if m.Type == models.TypeHTTP {
		return s.sendHTTPAlert(userEmail, m, result, duration)
	}

	if m.Type == models.TypeDNS {

		if result.Status == models.StatusUp {
			return s.sendDNSDetectedAlert(userEmail, m, result, dnsType)
		}

		if result.Status == models.StatusDown && result.ResultValue != "" {
			return s.sendDNSChangedAlert(userEmail, m, result, dnsType)
		}

		return s.sendDNSStatusAlert(userEmail, m, result, dnsType)
	}

	return nil
}

func (s *EmailService) sendHTTPAlert(to string, m models.Monitor, res models.CheckResult, d time.Duration) error {
	emoji := "🟢"
	color := "#38a169"

	timeLayout := "02/01/2006 15:04:05"
	var timeDetails string
	if res.Status == models.StatusDown {
		emoji = "🔴"
		color = "#e53e5e"
		timeDetails += fmt.Sprintf("<p><strong>Começou em:</strong> %s</p>", res.CheckedAt.Format(timeLayout))
	} else if res.Status == models.StatusUp && m.LastCheckStatus == models.StatusDown {
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

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendDNSStatusAlert(to string, m models.Monitor, res models.CheckResult, dnsConfig models.DNSConfig) error {
	subject := fmt.Sprintf("⚠️ Falha de DNS tipo %s: %s", dnsConfig.RecordType, m.Target)

	body := fmt.Sprintf(`
		<h2>Problema de Resolução DNS</h2>
		<p>Não foi possível verificar os registros DNS para <strong>%s</strong>.</p>
		<p><strong>Status:</strong> %s</p>
		<p><strong>Erro Técnico:</strong> %s</p>
		<p><em>Verifique se o domínio expirou ou se os servidores de nome estão respondendo.</em></p>
	`, m.Target, res.Status, res.Message)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendDNSChangedAlert(to string, m models.Monitor, res models.CheckResult, dnsConfig models.DNSConfig) error {
	subject := fmt.Sprintf("🚨 DNS tipo %s de %s foi Alterado!", dnsConfig.RecordType, m.Target)

	body := fmt.Sprintf(`
		<div style="border: 2px solid red; padding: 15px; background-color: #fff5f5;">
			<h2 style="color: red;">Alteração de Registro Detectada</h2>
			<p>O registro DNS monitorado não corresponde à configuração esperada.</p>
			
			<ul>
				<li><strong>Alvo:</strong> %s</li>
				<li><strong>Valor Encontrado (Atual):</strong> <code>%s</code></li>
				<li><strong>Mensagem do Sistema:</strong> %s</li>
				<li><strong>Detectado em: %s</strong></li>
			</ul>

			<p><strong>Ação Recomendada:</strong> Verifique imediatamente se o seu domínio foi comprometido ou se houve uma atualização não planejada.</p>
		</div>
	`, m.Target, res.ResultValue, res.Message, res.CheckedAt)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendDNSDetectedAlert(to string, m models.Monitor, res models.CheckResult, dnsConfig models.DNSConfig) error {
	subject := fmt.Sprintf("🟢 DNS tipo %s Detectado: %s", dnsConfig.RecordType, m.Target)

	body := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px; border: 1px solid #38a169; border-radius: 5px; background-color: #f0fff4;">
			<h2 style="color: #38a169;">Monitoramento DNS Ativo</h2>
			<p>O monitoramento para <strong>%s</strong> foi atualizado com sucesso.</p>
			
			<ul>
				<li><strong>Status:</strong> <span style="color: #38a169; font-weight: bold;">UP (Ativo)</span></li>
				<li><strong>Valor Detectado:</strong> <code>%s</code></li>
			</ul>

			<p style="font-size: 12px; color: #666;">A partir de agora, avisaremos se esse valor mudar.</p>
		</div>
	`, m.Target, res.ResultValue)

	return s.SendEmail(to, subject, body)
}
