package handlers

import (
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	DB           *pgxpool.Pool
	TokenAuth *jwtauth.JWTAuth
}
