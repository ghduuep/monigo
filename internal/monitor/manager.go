package monitor

import (
	"context"
	"time"

	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartMonitoring(ctx context.Context, db *pgxpool.Pool) {

}

func runMonitorRoutine(ctx context.Context, db *pgxpool.Pool, m *models.Monitor) {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			result := performCheck(m)
			//Save the result to the database
			// Compare the result from the previous check and send notifications if needed
		}
	}
}

func performCheck(m *models.Monitor) models.CheckResult {
	switch m.Type {
	case models.TypeHTTP:
		return checkHTTP(m)
	case models.TypeDNS:
		return checkDNS(m)
	default:
		return models.CheckResult{
			Status: models.StatusDown, Message: "Unknown type",
		}
	}
}
