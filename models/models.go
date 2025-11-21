package models

import "time"

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
	LastChecked time.Time     `json:"last_checked"`
}

type CheckLog struct {
	ID        int       `json:"id"`
	WebsiteID int       `json:"website_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
