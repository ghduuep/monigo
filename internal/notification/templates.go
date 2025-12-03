package notification

import (
	"fmt"

	"github.com/ghduuep/pingly/internal/models"
)

func (s *EmailService) SendStatusAlert(userEmail string, m models.Monitor, result models.CheckResult) error {

	if m.Type == models.TypeHTTP {
		return s.sendHTTPAlert(userEmail, m, result)
	}

	if m.Type == models.TypeDNS {
		if result.Status == models.StatusDown && result.ResultValue != "" {
			return s.sendDNSChangedAlert(userEmail, m, result) // Caso 3
		}

		return s.sendDNSStatusAlert(userEmail, m, result) // Caso 2
	}

	return nil
}

func (s *EmailService) sendHTTPAlert(to string, m models.Monitor, res models.CheckResult) error {
	subject := fmt.Sprintf("📡 Monitor HTTP: %s está %s", m.Target, res.Status)
	color := "#e53e3e" // Vermelho
	if res.Status == models.StatusUp {
		color = "#38a169"
	} // Verde

	body := fmt.Sprintf(`
		<h2>Atualização de Status HTTP</h2>
		<p>O monitor <strong>%s</strong> mudou para <span style="color:%s"><strong>%s</strong></span>.</p>
		<p><strong>Motivo:</strong> %sms</p>
		<p><strong>Latência:</strong> %v</p>
	`, m.Target, color, res.Status, res.Message, res.Latency)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendDNSStatusAlert(to string, m models.Monitor, res models.CheckResult) error {
	subject := fmt.Sprintf("⚠️ Falha de DNS: %s", m.Target)

	body := fmt.Sprintf(`
		<h2>Problema de Resolução DNS</h2>
		<p>Não foi possível verificar os registros DNS para <strong>%s</strong>.</p>
		<p><strong>Status:</strong> %s</p>
		<p><strong>Erro Técnico:</strong> %s</p>
		<p><em>Verifique se o domínio expirou ou se os servidores de nome estão respondendo.</em></p>
	`, m.Target, res.Status, res.Message)

	return s.SendEmail(to, subject, body)
}

func (s *EmailService) sendDNSChangedAlert(to string, m models.Monitor, res models.CheckResult) error {
	subject := fmt.Sprintf("🚨 ALERTA CRÍTICO: DNS de %s foi Alterado!", m.Target)

	body := fmt.Sprintf(`
		<div style="border: 2px solid red; padding: 15px; background-color: #fff5f5;">
			<h2 style="color: red;">Alteração de Registro Detectada</h2>
			<p>O registro DNS monitorado não corresponde à configuração esperada.</p>
			
			<ul>
				<li><strong>Alvo:</strong> %s</li>
				<li><strong>Valor Encontrado (Atual):</strong> <code>%s</code></li>
				<li><strong>Mensagem do Sistema:</strong> %s</li>
			</ul>

			<p><strong>Ação Recomendada:</strong> Verifique imediatamente se o seu domínio foi comprometido ou se houve uma atualização não planejada.</p>
		</div>
	`, m.Target, res.ResultValue, res.Message)

	return s.SendEmail(to, subject, body)
}
