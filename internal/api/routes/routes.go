package routes

import (
	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(db *pgxpool.Pool, siteChan chan *models.Website) *chi.Mux {
	h := &handlers.Handler{
		DB:           db,
		NewSitesChan: siteChan,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Get("/websites", h.GetAllWebsites)
		r.Post("/websites", h.CreateWebsite)
	})

	return r
}
