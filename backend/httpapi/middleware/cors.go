package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"exodus/backend/config"
	"exodus/constant"
)

var (
	defaultCORSMethods = []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	defaultCORSHeaders = []string{"Content-Type", "Authorization", "X-API-Token", "Cookie"}
)

// WithCORS adds CORS headers and Server header to all responses.
// If allowedOrigins contains "*", all origins are allowed.
// Otherwise, only origins in the list are allowed.
func WithCORS(cfg *config.BackendConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg != nil && !cfg.Panel.AllowInsecureHTTP && !IsSecureRequest(r, cfg) && !isHealthPath(r.URL.Path) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUpgradeRequired)
			_, _ = w.Write([]byte(`{"error":"https_required","message":"panel requires HTTPS"}`))
			return
		}

		// Add Server header
		serverHeader := fmt.Sprintf("MuxCloud/%s (WebServer)", constant.Version)
		w.Header().Set("Server", serverHeader)
		w.Header().Set("X-Powered-By", "MuxCloud")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")

		// Use CORS settings from config
		allowedOrigins := cfg.CORS.AllowedOrigins

		// If no allowed origins configured, do not emit CORS headers.

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
			w.Header().Add("Vary", "Origin")
		}
		if allowOrigin != "" {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(defaultCORSMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(defaultCORSHeaders, ", "))
			w.Header().Set("Access-Control-Max-Age", "86400")
		}

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			if allowOrigin != "" {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsSecureRequest(r *http.Request, cfg *config.BackendConfig) bool {
	if r.TLS != nil {
		return true
	}
	if cfg == nil {
		return false
	}

	remoteIP := parseRemoteIP(r.RemoteAddr)
	if remoteIP != nil && remoteIP.IsLoopback() {
		return true
	}
	if remoteIP == nil || !cfg.Panel.IsTrustedProxy(remoteIP) {
		return false
	}

	if proto := firstHeaderValue(r.Header.Get("X-Forwarded-Proto")); proto != "" {
		return strings.EqualFold(proto, "https")
	}
	if proto := firstHeaderValue(r.Header.Get("X-Forwarded-Scheme")); proto != "" {
		return strings.EqualFold(proto, "https")
	}
	if strings.EqualFold(r.Header.Get("X-Forwarded-Ssl"), "on") {
		return true
	}
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		if proto := parseForwardedProto(forwarded); proto != "" {
			return strings.EqualFold(proto, "https")
		}
	}
	return false
}

func parseRemoteIP(remoteAddr string) net.IP {
	if remoteAddr == "" {
		return nil
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err == nil {
		return net.ParseIP(host)
	}
	return net.ParseIP(remoteAddr)
}

func firstHeaderValue(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}

func parseForwardedProto(forwarded string) string {
	parts := strings.Split(forwarded, ",")
	for _, part := range parts {
		segments := strings.Split(part, ";")
		for _, segment := range segments {
			segment = strings.TrimSpace(segment)
			if strings.HasPrefix(strings.ToLower(segment), "proto=") {
				return strings.Trim(strings.TrimPrefix(segment, "proto="), "\"")
			}
		}
	}
	return ""
}

func isHealthPath(path string) bool {
	return path == "/api/health"
}
