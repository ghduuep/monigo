package models

import (
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
	TypeDNS_A    MonitorType = "dns_a"
	TypeDNS_AAAA MonitorType = "dns_aaaa"
	TypeDNS_MX   MonitorType = "dns_mx"
	TypeDNS_NS   MonitorType = "dns_ns"
)

type MonitorStatus string

const (
	StatusUp      MonitorStatus = "up"
	StatusDown    MonitorStatus = "down"
	StatusChanged MonitorStatus = "changed"
	StatusUnknown MonitorStatus = "unknown"
)

type Monitor struct {
	ID              int           `json:"id"`
	UserID          int           `json:"user_id"`
	Target          string        `json:"target"`
	Type            MonitorType   `json:"type"`
	ExpectedValue   string        `json:"expected_value,omitempty"`
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
