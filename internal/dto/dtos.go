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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MonitorRequest struct {
	Target   string             `json:"target" db:"target"`
	Type     models.MonitorType `json:"type" db:"type"`
	Config   json.RawMessage    `json:"config" db:"config" swaggertype:"string"`
	Interval time.Duration      `json:"interval" db:"interval" swaggertype:"integer"`
}
