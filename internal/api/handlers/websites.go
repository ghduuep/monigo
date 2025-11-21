package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
)

func (h *Handler) CreateWebsite(w http.ResponseWriter, req *http.Request) {

	var newWebsite models.Website

	if err := json.NewDecoder(req.Body).Decode(&newWebsite); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	if err := database.CreateWebsite(req.Context(), h.DB, &newWebsite); err != nil {
		http.Error(w, "Erro ao criar website", http.StatusInternalServerError)
		return
	}

	go func() {
		h.NewSitesChan <- &newWebsite
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newWebsite)
}

func (h *Handler) GetAllWebsites(w http.ResponseWriter, req *http.Request) {
	websites, err := database.GetAllWebsites(req.Context(), h.DB)
	if err != nil {
		http.Error(w, "Erro ao buscar websites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)
}
