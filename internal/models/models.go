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

type Website struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	URL         string        `json:"url"`
	Interval    time.Duration `json:"interval"`
	LastStatus  string        `json:"last_status"`
	LastChecked *time.Time    `json:"last_checked"`
}

type UptimeLog struct {
	ID        int       `json:"id"`
	WebsiteID int       `json:"website_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DNSLog struct {
	ID          int       `json:"id"`
	DNSDomainID int       `json:"dns_domain_id"`
	Diff        string    `json:"diff"`
	CreatedAt   time.Time `json:"created_at"`
}

type DNSDomains struct {
	ID          int           `json:"id"`
	UserID      int           `json:"user_id"`
	Domain      string        `json:"domain"`
	Interval    time.Duration `json:"interval"`
	LastA       []string      `json:"last_a_records"`
	LastAAAA    []string      `json:"last_aaaa_records"`
	LastMX      []string      `json:"last_mx_records"`
	LastNS      []string      `json:"last_ns_records"`
	LastChecked time.Time     `json:"checked_at"`
}
