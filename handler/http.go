package handler

import (
	"encoding/json"
	"net/http"

	"github.com/capucinoxx/forlorn/manager"
	"github.com/capucinoxx/forlorn/network"
)

type HttpHandler struct {
	gm *manager.GameManager
	am *manager.AuthManager
}

// NewHttpHandler crée un nouveau gestionnaire HTTP.
func NewHttpHandler(gm *manager.GameManager, am *manager.AuthManager) *HttpHandler {
	return &HttpHandler{
		gm: gm,
		am: am,
	}
}

// Handle commence à écouter les différentes routes HTTP et les associe à des fonctions.
func (h *HttpHandler) Handle() {
	network.HandleFunc("/start", h.startGame)
	network.HandleFunc("/create", h.createAccount)
}

// createAccount crée un compte utilisateur et retourne un jeton d'authentification.
func (h *HttpHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username string `json:"username"`
		Sides    int    `json:"sides"`
	}{}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, _ := h.am.Register(payload.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// startGame démarre le serveur de jeu.
func (h *HttpHandler) startGame(w http.ResponseWriter, r *http.Request) {
	h.gm.Start()
}
