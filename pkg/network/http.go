package network

import (
	"net/http"

	"github.com/capucinoxx/forlorn/pkg/utils"
)

// Middleware est une fonction qui prend un http.Handler et renvoie un http.Handler.
type Middleware func(http.Handler) http.Handler

// logMiddleware est un Middleware qui enregistre les requêtes HTTP.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Log("HTTP", r.Method, "request from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// headers est un Middleware qui définit les en-têtes de réponse HTTP
// en ajoutant certaines protections et configurations.
func headers(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prévenir le navigateur de ne pas effectuer de requêtes DNS préalables.
		w.Header().Set("X-DNS-Prefetch-Control", "off")

		// Refuse les requêtes de navigateur à travers les frames.
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")

		// Protection contre les attaques XSS.
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// TODO: Changer l'accès à l'origine pour correspondre à l'URL du jeu.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// chain renvoie une fonction http.HandlerFunc qui exécute les middlewares
// dans l'ordre donné.
func chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	handler = logMiddleware(headers(handler)).ServeHTTP
	for _, middleware := range middlewares {
		handler = middleware(handler).ServeHTTP
	}
	return handler
}

// HandleFunc enregistre un gestionnaire de requêtes HTTP avec des middlewares.
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	http.HandleFunc(pattern, chain(handler, middlewares...))
}
