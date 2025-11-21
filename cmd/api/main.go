package main

import (
	"context"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/monitor"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	db = database.InitDB()

	defer db.Close()

	ctx := context.Background()

	monitor.StartMonitoring(ctx, db)

	select {}

	// healthHandler := func(w http.ResponseWriter, req *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`{"status":"ok"}`))
	// }

	// http.HandleFunc("/health", healthHandler)
	// http.HandleFunc("/websites", websitesHandler)
	// http.ListenAndServe(":8080", nil)
	// log.Println("Server started on :8080")
}

// func websitesHandler(w http.ResponseWriter, req *http.Request) {
// 	switch req.Method{
// 		case http.MethodGet:
// 			getAllWebsites(w, req)
// 		case http.MethodPost:
// 			createWebsite(w, req)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func getAllWebsites(w http.ResponseWriter, req *http.Request) {
// 	websites, err := database.GetWebsites(req.Context(), db)

// 	if err != nil {
// 		http.Error(w, "Erro ao buscar websites", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(websites)
// }

// func createWebsite(w http.ResponseWriter, req *http.Request) {
// 	var newWebsite models.Website

// 	if err := json.NewDecoder(req.Body).Decode(&newWebsite); err != nil {
// 		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
// 		return
// 	}

// 	if err := database.CreateWebsite(req.Context(), db, &newWebsite); err != nil {
// 		log.Printf("Erro na criação do website: %v", err)
// 		http.Error(w, "Erro ao criar website", http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(newWebsite)
// }
