package routes

import (
	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)


func NewRouter(db *pgxpool.Pool, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	h := &handlers.Handler{
		DB:           db,
		TokenAuth:    tokenAuth,
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", h.Login)
	r.Post("/register", h.Register)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/websites", h.GetAllWebsites)
		r.Post("/websites", h.CreateWebsite)
	})

	return r
}
