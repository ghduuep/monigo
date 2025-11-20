package models

import "time"

type Website struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	URL string `json:"url"`
	Interval time.Duration `json:"interval"`
	LastStatus string `json:"last_status"`
	LastChecked *time.Time `json:"last_checked"`
}
