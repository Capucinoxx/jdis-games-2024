package handler

import (
	"encoding/json"
	"net/http"

	"github.com/capucinoxx/forlorn/pkg/manager"
	"github.com/capucinoxx/forlorn/pkg/network"
)

type HttpHandler struct {
	gm *manager.GameManager
	am *manager.AuthManager
	sm *manager.ScoreManager
}

type HttpResponse struct {
	Type    string `json:"type"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func NewHttpHandler(gm *manager.GameManager, am *manager.AuthManager, sm *manager.ScoreManager) *HttpHandler {
	return &HttpHandler{
		gm: gm,
		am: am,
		sm: sm,
	}
}

func (h *HttpHandler) Handle() {
	fs := http.FileServer(http.Dir("./dist"))

	network.Handle("/", fs)
	network.HandleFunc("/start", h.startGame, h.adminOnly)
	network.HandleFunc("/create", h.register)
	network.HandleFunc("/leaderboard", h.leaderboard, h.checkLeaderboardAccess)
	network.HandleFunc("/toggle_leaderboard", h.toggleLeaderboard, h.adminOnly)
	network.HandleFunc("/kill", h.kill, h.adminOnly)

	network.HandleFunc("/freeze", h.freeze, h.adminOnly)
	network.HandleFunc("/unfrezze", h.unfreeze, h.adminOnly)
}

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

func (h *HttpHandler) freeze(w http.ResponseWriter, r *http.Request) {
	h.gm.Freeze(true)
}

func (h *HttpHandler) unfreeze(w http.ResponseWriter, r *http.Request) {
	h.gm.Freeze(false)
}
