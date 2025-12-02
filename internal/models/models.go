package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type MonitorType string

const (
	TypeHTTP     MonitorType = "http"
	TypePing     MonitorType = "ping"
	TypeDNS MonitorType = "dns"
}

type MonitorStatus string

const (
	StatusUp      MonitorStatus = "up"
	StatusDown    MonitorStatus = "down"
	StatusUnknown MonitorStatus = "unknown"
)

type Monitor struct {
	ID              int           `json:"id"`
	UserID          int           `json:"user_id"`
	Target          string        `json:"target"`
	Type            MonitorType   `json:"type"`
	Config json.RawMessage `json:"config"`
	Interval        time.Duration `json:"interval"`
	LastCheckStatus MonitorStatus `json:"last_check_status"`
	LastCheckAt     time.Time     `json:"last_check_at"`
	CreatedAt       time.Time     `json:"created_at"`
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

type DNSConfig struct {
	RecordType string `json:"record_type"`
	ExpectedValue string `json:"expected_value"`
}
