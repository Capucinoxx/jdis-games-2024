package handler

import "net/http"

func (h *HttpHandler) checkLeaderboardAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/leaderboard" && !h.sm.IsVisible() {
			http.Error(w, "Leaderboard access is disabled.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
