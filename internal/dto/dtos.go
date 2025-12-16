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
	Target           string             `json:"target" db:"target" validate:"required"`
	Type             models.MonitorType `json:"type" db:"type" validate:"required,oneof=http dns port"`
	Config           json.RawMessage    `json:"config" db:"config" swaggertype:"string"`
	Interval         string             `json:"interval" validate:"required,oneof=30s 1m 5m 30m 1h 12h 24h"`
	Timeout          string             `json:"timeout" validate:"required,oneof=1s 15s 30s 45s 60s"`
	LatencyThreshold int64              `json:"latency_threshold_ms" db:"latency_threshold_ms" validate:"min=0"`
}

type MonitorStatsResponse struct {
	MonitorID        int     `json:"monitor_id"`
	UptimePercentage float64 `json:"uptime_percentage"`
	AvgLatency       float64 `json:"avg_latency"`
	MinLatency       float64 `json:"min_latency"`
	MaxLatency       float64 `json:"max_latency"`
	TotalChecks      int     `json:"total_checks"`
	ApdexScore       float64 `json:"apdex_score"`
}

type CreateChannelRequest struct {
	Type   string `json:"type" validate:"required,oneof=email telegram sms"`
	Target string `json:"target" validate:"required"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=6"`
}

type UpdateMonitorRequest struct {
	Target           *string         `json:"target" validate:"omitempty"`
	Interval         *string         `json:"interval" validate:"omitempty,oneof=30s 1m 5m 30m 1h 12h 24h"`
	Timeout          *string         `json:"timeout" validate:"omitempty,oneof=1s 30s 45s 60s"`
	Config           json.RawMessage `json:"config" validate:"omitempty"`
	LatencyThreshold *int64          `json:"latency_threshold_ms" db:"latency_threshold_ms"`
}

type UpdateChannelRequest struct {
	Type    *string `json:"type" validate:"omitempty,oneof=email telegram sms"`
	Target  *string `json:"target" validate:"omitempty"`
	Enabled *bool   `json:"enabled" validate:"omitempty"`
}

type MonitorSummaryResponse struct {
	Total    int `json:"total"`
	Up       int `json:"up"`
	Down     int `json:"down"`
	Degraded int `json:"degraded"`
}

type PaginationQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	Perpage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

type PaginatedResponse struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta"`
}
