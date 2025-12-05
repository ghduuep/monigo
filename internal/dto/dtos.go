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
	Username string `json:"username" validate:"required, alpha"`
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"password" validate:"required, min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type MonitorRequest struct {
	Target   string             `json:"target" db:"target" validate:"required"`
	Type     models.MonitorType `json:"type" db:"type" validate:"required,oneof=http dns"`
	Config   json.RawMessage    `json:"config" db:"config" swaggertype:"string"`
	Interval time.Duration      `json:"interval" db:"interval" validate:"required, min=30000000000"`
}
