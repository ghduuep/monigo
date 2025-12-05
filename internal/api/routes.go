package api

import (
	"os"

	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRotes(e *echo.Echo, db *pgxpool.Pool) {
	handler := handlers.NewHandler(db)

	v1 := e.Group("/api/v1")

	v1.POST("/register", handler.Register)
	v1.POST("/login", handler.Login)

	protected := e.Group("/api/v1")

	jwtSecret := os.Getenv("JWT_SECRET")

	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
	}))

	protected.GET("/users", handler.GetUsers)
	protected.GET("/users/:id", handler.GetUserByID)
	protected.GET("/monitors", handler.GetMonitors)
	protected.GET("/monitors/:id", handler.GetMonitorByID)

	protected.POST("/monitors", handler.CreateMonitor)
	protected.DELETE("/monitors/:id", handler.DeleteMonitor)
	protected.DELETE("/users/:id", handler.DeleteUser)
}
