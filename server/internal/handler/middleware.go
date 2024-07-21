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

func (h *HttpHandler) adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.URL.Query().Get("tkn")

		if tkn == "" {
			cookie, err := r.Cookie("tkn")
			if err == nil {
				tkn = cookie.Value
			}
		}

		if tkn == "" {
			http.Error(w, "admin restricted", http.StatusForbidden)
			return
		}

		_, _, isAdmin, _ := h.am.Authenticate(tkn)
		if !isAdmin {
			http.Error(w, "admin restricted", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
