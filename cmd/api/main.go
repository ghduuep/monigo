package main

import (
	"log"

	"github.com/ghduuep/pingly/internal/api"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	db := database.InitDB()
	defer db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(20))))

	api.SetupRotes(e, db)

	port := ":8080"
	log.Printf("API server is running on port %s", port)
	if err := e.Start(port); err != nil {
		e.Logger.Fatal(err)
	}
}
