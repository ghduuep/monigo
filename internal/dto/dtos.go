package dto

import (
	"encoding/json"
	"github.com/ghduuep/pingly/internal/models"
	"time"
)

type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,alpha"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type MonitorRequest struct {
	Target   string             `json:"target" db:"target" validate:"required"`
	Type     models.MonitorType `json:"type" db:"type" validate:"required,oneof=http dns"`
	Config   json.RawMessage    `json:"config" db:"config" swaggertype:"string"`
	Interval time.Duration      `json:"interval" db:"interval" validate:"required,min=30000000000"`
}

type MonitorStatsResponse struct {
	MonitorID        int                  `json:"monitor_id"`
	LastStatus       models.MonitorStatus `json:"last_status"`
	LastCheckAt      *time.Time           `json:"last_check_at"`
	UptimePercentage float64              `json:"uptime_percentage"`
	AvgLatency       float64              `json:"avg_latency"`
	Last24HChecks    int                  `json:"last_24h_checks"`
}

type CreateChannelRequest struct {
	Type   string `json:"type" validate:"required,oneof="email telegram"`
	Target string `json:"target" validate:"required"`
}
