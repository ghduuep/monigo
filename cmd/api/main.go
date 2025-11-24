package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ghduuep/pingly/internal/api/routes"
	"github.com/ghduuep/pingly/internal/database"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	db := database.InitDB()
	defer db.Close()

	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("JWT_SECRET")), nil)

	router := routes.NewRouter(db, tokenAuth)

	log.Printf("API server is running on %s", os.Getenv("SERVER_PORT"))

	if err := http.ListenAndServeTLS(":" + os.Getenv("SERVER_PORT"), router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
