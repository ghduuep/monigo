package routes

import (
	"net/http"

	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool) *http.ServeMux {
	h := &handlers.Handler{
		DB: db,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /websites", h.CreateWebsite)
	mux.HandleFunc("GET /websites", h.GetAllWebsites)

	return mux
}
