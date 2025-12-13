package api

import (
	"net/http"
	"os"

	"github.com/ghduuep/pingly/internal/api/handlers"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/jackc/pgx/v5/pgxpool"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupRotes(e *echo.Echo, db *pgxpool.Pool, rdb *redis.Client) {
	handler := handlers.NewHandler(db, rdb)

	v1 := e.Group("/api/v1")

	v1.POST("/register", handler.Register)
	v1.POST("/login", handler.Login)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	protected := e.Group("/api/v1")

	jwtSecret := os.Getenv("JWT_SECRET")

	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecret),
	}))

	protected.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if len(tokenString) > 7 {
				tokenString = tokenString[7:]
			}

			blacklisted, err := database.IsTokenBlackListed(c.Request().Context(), rdb, tokenString)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Auth service unavailable"})
			}

			if blacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"error": "Invalid token."})
			}

			return next(c)
		}
	})

	protected.GET("/users", handler.GetUsers)
	protected.GET("/users/:id", handler.GetUserByID)
	protected.GET("/monitors", handler.GetMonitors)
	protected.GET("/monitors/:id", handler.GetMonitorByID)
	protected.GET("/monitors/:id/stats", handler.GetMonitorStats)
	protected.GET("/monitors/:id/checks", handler.GetMonitorLastChecks)
	protected.GET("/channels", handler.GetChannels)
	protected.GET("/monitors/:id/incidents", handler.GetMonitorLastIncidents)

	protected.POST("/logout", handler.Logout)
	protected.POST("/channels", handler.CreateChannel)
	protected.POST("/monitors", handler.CreateMonitor)
	protected.DELETE("/channels/:id", handler.DeleteChannel)
	protected.DELETE("/monitors/:id", handler.DeleteMonitor)
	protected.PATCH("/users", handler.UpdateUser)
	protected.DELETE("/users", handler.DeleteUser)
}
