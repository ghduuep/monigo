package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/ghduuep/pingly/internal/models"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) CreateWebsite(w http.ResponseWriter, req *http.Request) {

	var newWebsite models.Website

	if err := json.NewDecoder(req.Body).Decode(&newWebsite); err != nil {
		http.Error(w, "Erro ao decodificar o corpo da requisição", http.StatusBadRequest)
		return
	}

	newWebsite.LastStatus = "UNKNOWN"

	if err := database.CreateWebsite(req.Context(), h.DB, &newWebsite); err != nil {
		http.Error(w, "Erro ao criar website", http.StatusInternalServerError)
		return
	}

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

func (h *Handler) DeleteWebsite(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	if err := database.DeleteWebsite(req.Context(), h.DB, id); err != nil {
		http.Error(w, "Erro ao deletar website", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
