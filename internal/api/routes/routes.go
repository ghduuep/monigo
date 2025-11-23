package routes

import (
	"log"

	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)


func NewRouter(db *pgxpool.Pool, siteChan chan *models.Website, tokenAuth *jwtauth.JWTAuth) *chi.Mux {
	h := &handlers.Handler{
		DB:           db,
		NewSitesChan: siteChan,
		TokenAuth:    tokenAuth,
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
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
