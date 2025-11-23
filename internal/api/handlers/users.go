package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ghduuep/pingly/internal/database"
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByEmail(req.Context(), h.DB, creds.Email)
	if err != nil {
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials.", http.StatusUnauthorized)
		return
	}

	claims := map[string]any{
		"user_id": user.ID,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(24 * time.Hour))

	_, tokenString, err := h.TokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}


func(h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(req.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	err = database.CreateUser(req.Context(), h.DB, creds.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
