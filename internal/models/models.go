package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type MonitorType string

const (
	TypeHTTP MonitorType = "http"
	TypePing MonitorType = "ping"
	TypeDNS  MonitorType = "dns"
)

type MonitorStatus string

const (
	StatusUp      MonitorStatus = "up"
	StatusDown    MonitorStatus = "down"
	StatusUnknown MonitorStatus = "unknown"
)

type Monitor struct {
	ID            int           `json:"id"`
	UserID        int           `json:"user_id"`
	Target        string        `json:"target"`
	Type          MonitorType   `json:"type"`
	ExpectedValue string        `json:"expected_value"`
	Interval      time.Duration `json:"interval"`
	CreatedAt     time.Time     `json:"created_at"`
}

type CheckResult struct {
	ID         int           `json:"id"`
	MonitorID  int           `json:"monitor_id"`
	Status     MonitorStatus `json:"status"`
	Latency    time.Duration `json:"latency"`
	StatusCode int           `json:"status_code,omitempty"`
	Message    string        `json:"message,omitempty"`
	CheckedAt  time.Time     `json:"checked_at"`
}
