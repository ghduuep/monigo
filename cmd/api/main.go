package main

import (
	"log"

	_ "github.com/ghduuep/pingly/docs"
	"github.com/ghduuep/pingly/internal/api"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// @title Pingly API
// @version 1.0
// @description API for Pingly monitoring service.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@pingly.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
