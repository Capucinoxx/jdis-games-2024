package handler

// Package handler provides an abstraction for managing HTTP requests and
// responses in application. This package includes functionalities
// for handling user registration, game start, leaderboard management, and
// user management.

import "net/http"

// checkLeaderboardAccess is a middleware function that checks if the leaderboard
// is accessible. If the leaderboard is not visible, it returns a 403 Forbidden
// status. Otherwise, it passes the request to the next handler.
func (h *HttpHandler) checkLeaderboardAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/leaderboard" && !h.sm.IsVisible() {
			http.Error(w, "Leaderboard access is disabled.", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// adminOnly is a middleware function that restricts access to certain routes
// to only admin users. It checks for an admin token in the URL query or cookies.
// If the token is not present or the user is not an admin, it returns a 403
// Forbidden status. Otherwise, it passes the request to the next handler.
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
