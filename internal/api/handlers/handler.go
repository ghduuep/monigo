package handlers

import "github.com/jackc/pgx/v5/pgxpool"

type Handler struct {
	DB *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{DB: db}
}
