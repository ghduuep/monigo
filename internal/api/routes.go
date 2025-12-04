package api

import (
	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func SetupRotes(e *echo.Echo, db *pgxpool.Pool) {
	handler := handlers.NewHandler(db)

	apiGroup := e.Group("/api/v1")

	apiGroup.GET("/users", handler.GetUsers)
	apiGroup.GET("/users/:id", handler.GetUserByID)
}
