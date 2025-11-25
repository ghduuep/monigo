package notification

import "github.com/ghduuep/pingly/internal/models"

func (s *EmailService) SendStatusAlert(userEmail string, m models.Monitor, newStatus models.MonitorStatus, rootCause string) error {

}
