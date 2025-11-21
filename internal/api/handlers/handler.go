package handlers

import (
	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB           *pgxpool.Pool
	NewSitesChan chan *models.Website
}
