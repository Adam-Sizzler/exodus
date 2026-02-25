package api

import (
	"fmt"
	"net/http"
	"strings"
	"v2ray-stat/backend/config"
	"v2ray-stat/constant"
)

// Answer handles basic server information requests.
func Answer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverHeader := fmt.Sprintf("MuxCloud/%s (WebServer)", constant.Version)
		w.Header().Set("Server", serverHeader)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Powered-By", "MuxCloud")
		fmt.Fprintf(w, "MuxCloud / %s\n", constant.Version)
	}
}

// WithCORS adds CORS headers to all responses to allow cross-origin requests from the frontend.
// If allowedOrigins contains "*", all origins are allowed.
// Otherwise, only origins in the list are allowed.
func WithCORS(cfg *config.BackendConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use CORS settings from config
		allowedOrigins := cfg.CORS.AllowedOrigins
		allowedMethods := cfg.CORS.AllowedMethods
		allowedHeaders := cfg.CORS.AllowedHeaders
		
		// Set defaults if not configured
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}
		if len(allowedMethods) == 0 {
			allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
		}
		if len(allowedHeaders) == 0 {
			allowedHeaders = []string{"Content-Type", "Authorization", "X-API-Token"}
		}
		
		origin := r.Header.Get("Origin")
		allowOrigin := ""
		
		// Check if "*" is in the list (allow all)
		for _, o := range allowedOrigins {
			if o == "*" {
				allowOrigin = "*"
				break
			}
			if o == origin {
				allowOrigin = origin
				break
			}
		}
		
		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}
		
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
		w.Header().Set("Access-Control-Max-Age", "86400")
		
		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// WithServerHeader adds a Server header to all responses.
func WithServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverHeader := fmt.Sprintf("MuxCloud/%s (WebServer)", constant.Version)
		w.Header().Set("Server", serverHeader)
		w.Header().Set("X-Powered-By", "MuxCloud")
		next.ServeHTTP(w, r)
	})
}
