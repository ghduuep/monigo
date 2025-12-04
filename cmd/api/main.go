package main

import (
	"github.com/ghduuep/pingly/internal/database"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cannot load .env file.")
	}

	db := database.InitDB()
	defer db.Close()

}
