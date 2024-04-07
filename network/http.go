package network

import (
	"net/http"

	"github.com/capucinoxx/forlorn/utils"
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

// chain renvoie une fonction http.HandlerFunc qui exécute les middlewares
// dans l'ordre donné.
func chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	handler = logMiddleware(handler).ServeHTTP
	for _, middleware := range middlewares {
		handler = middleware(handler).ServeHTTP
	}
	return handler
}

// HandleFunc enregistre un gestionnaire de requêtes HTTP avec des middlewares.
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	http.HandleFunc(pattern, chain(handler, middlewares...))
}
