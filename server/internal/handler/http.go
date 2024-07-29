package handler

// Package handler provides an abstraction for managing HTTP requests and
// responses in application. This package includes functionalities
// for handling user registration, game start, leaderboard management, and
// user management.

import (
	"encoding/json"
	"net/http"

	"github.com/capucinoxx/jdis-games-2024/pkg/manager"
	"github.com/capucinoxx/jdis-games-2024/pkg/network"
)

// HttpHandler is a structure that holds references to the game manager,
// authentication manager, and score manager.
type HttpHandler struct {
	gm *manager.GameManager
	am *manager.AuthManager
	sm *manager.ScoreManager
}

// HttpResponse is a structure used for formatting JSON responses.
type HttpResponse struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// NewHttpHandler creates a new instance of HttpHandler.
func NewHttpHandler(gm *manager.GameManager, am *manager.AuthManager, sm *manager.ScoreManager) *HttpHandler {
	return &HttpHandler{
		gm: gm,
		am: am,
		sm: sm,
	}
}

// Handle sets up the HTTP routes and handlers for the server.
func (h *HttpHandler) Handle() {
	fs := http.FileServer(http.Dir("./dist"))

	network.Handle("/", fs)
	network.HandleFunc("/start", h.startGame, h.adminOnly)
	network.HandleFunc("/create", h.register)
	network.HandleFunc("/leaderboard", h.leaderboard, h.checkLeaderboardAccess)
	network.HandleFunc("/toggle_leaderboard", h.toggleLeaderboard, h.adminOnly)
	network.HandleFunc("/kill", h.kill, h.adminOnly)
	network.HandleFunc("/users", h.users, h.adminOnly)

	network.HandleFunc("/freeze", h.freeze, h.adminOnly)
	network.HandleFunc("/unfreeze", h.unfreeze, h.adminOnly)
}

// register handles user registration requests.
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

	token, err := h.am.Register(payload.Username, false)
	var resp HttpResponse
	resp.Subject = "Token generation"
	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		resp.Type = "success"
		resp.Message = token
	} else {
		resp.Type = "error"
		resp.Message = err.Error()
	}
	json.NewEncoder(w).Encode(resp)
}

// startGame handles requests to start the game.
// restrictions: admins only.
func (h *HttpHandler) startGame(w http.ResponseWriter, r *http.Request) {
	h.gm.Start()
}

// users handles requests to list all users.
// restrictions: admins only.
func (h *HttpHandler) users(w http.ResponseWriter, r *http.Request) {
	users, err := h.am.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
	}

	payload := map[string]interface{}{
		"users": users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payload)
}

// leaderboard handles requests to retrieve the leaderboard.
// restrictions: only if leaderboard isn't freeze.
func (h *HttpHandler) leaderboard(w http.ResponseWriter, r *http.Request) {
	leaderboard, histories, err := h.sm.Rank()
	if err != nil {
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Leaderboard []manager.PlayerScore            `json:"leaderboard"`
		Histories   map[string][]manager.PlayerEntry `json:"histories"`
	}{
		Leaderboard: leaderboard,
		Histories:   histories,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// toggleLeaderboard handles requests to toggle the visibility of the leaderboard.
// restrictions: admins only.
func (h *HttpHandler) toggleLeaderboard(w http.ResponseWriter, r *http.Request) {
	visible := h.sm.ToggleVisibility()

	status := "disabled"
	if visible {
		status = "enabled"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "leaderboard access has been " + status})
}

// kill handles requests to kill a player.
// restrictions: admins only.
func (h *HttpHandler) kill(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		h.gm.Kill(name)
	}
}

// freeze handles requests to freeze the game.
// restrictions: admins only.
func (h *HttpHandler) freeze(w http.ResponseWriter, r *http.Request) {
	h.gm.Freeze(true)
}

// // unfreeze handles requests to unfreeze and restart the game.
// restrictions: admins only.
func (h *HttpHandler) unfreeze(w http.ResponseWriter, r *http.Request) {
	h.gm.Freeze(false)
	h.gm.Start()
}
