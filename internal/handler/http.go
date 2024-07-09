package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/capucinoxx/forlorn/pkg/manager"
	"github.com/capucinoxx/forlorn/pkg/network"
)

type HttpHandler struct {
	gm *manager.GameManager
	am *manager.AuthManager
	sm *manager.ScoreManager
}

// NewHttpHandler crée un nouveau gestionnaire HTTP.
func NewHttpHandler(gm *manager.GameManager, am *manager.AuthManager, sm *manager.ScoreManager) *HttpHandler {
	return &HttpHandler{
		gm: gm,
		am: am,
		sm: sm,
	}
}

// Handle commence à écouter les différentes routes HTTP et les associe à des fonctions.
func (h *HttpHandler) Handle() {
	network.HandleFunc("/start", h.startGame)
	network.HandleFunc("/create", h.register)
	network.HandleFunc("/leaderboard", h.leaderboard, h.checkLeaderboardAccess)
	network.HandleFunc("/toggle_leaderboard", h.toggleLeaderboard)
	network.HandleFunc("/kill", h.kill)
}

// register crée un compte utilisateur et retourne un jeton d'authentification.
func (h *HttpHandler) register(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Username string `json:"username"`
	}{}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, _ := h.am.Register(payload.Username, false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"type": "success", "message": fmt.Sprintf("Token: %s", token)})
}

// startGame démarre le serveur de jeu.
func (h *HttpHandler) startGame(w http.ResponseWriter, r *http.Request) {
	h.gm.Start()
}

func (h *HttpHandler) leaderboard(w http.ResponseWriter, r *http.Request) {
	l, err := h.sm.Rank()
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(l)
}

func (h *HttpHandler) toggleLeaderboard(w http.ResponseWriter, r *http.Request) {
	visible := h.sm.ToggleVisibility()

	status := "disabled"
	if visible {
		status = "enabled"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "leaderboard access has been " + status})
}

func (h *HttpHandler) kill(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		h.gm.Kill(name)
	}
}
