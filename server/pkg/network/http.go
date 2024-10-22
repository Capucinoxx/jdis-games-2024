package network

// Package network provides an abstraction for managing HTTP requests and
// responses, including middleware for logging and setting response headers.
// This package is used to enhance the security and monitoring of HTTP
// communication

import (
	"net/http"

	"github.com/capucinoxx/jdis-games-2024/pkg/utils"
)

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// logMiddleware is a middleware that logs HTTP requests.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.Log("HTTP", r.Method, "request from %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// headers is a middleware that sets HTTP response headers
// by adding certain protections and configurations.
func headers(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// chain returns an http.HandlerFunc that executes the middlewares in the given order.
func chain(handler http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	handler = headers(handler).ServeHTTP
	handler = logMiddleware(handler).ServeHTTP
	for _, middleware := range middlewares {
		handler = middleware(handler).ServeHTTP
	}

	return handler
}

// HandleFunc registers an HTTP request handler with middlewares.
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request), middlewares ...Middleware) {
	http.HandleFunc(pattern, chain(handler, middlewares...))
}

// Handle registers an HTTP request handler with middlewares.
func Handle(pattern string, handle http.Handler, middlewares ...Middleware) {
	http.Handle(pattern, chain(handle.ServeHTTP, middlewares...))
}
