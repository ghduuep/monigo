package routes

import (
	"net/http"

	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, siteChan chan *models.Website) *http.ServeMux {
	h := &handlers.Handler{
		DB:           db,
		NewSitesChan: siteChan,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /websites", h.CreateWebsite)
	mux.HandleFunc("GET /websites", h.GetAllWebsites)

	return mux
}
