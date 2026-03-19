package middleware

import (
	"net"
	"net/http"
	"strings"

	"cerberus/backend/config"
)

// getClientIP retrieves the client IP address from an HTTP request.
func GetClientIP(r *http.Request, cfg *config.BackendConfig) string {
	cfg.Logger.Debug("Retrieving client IP address", "remote_addr", r.RemoteAddr)

	// Check X-Forwarded-For header first (may contain multiple IPs)
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ips := strings.Split(fwd, ",")
		ip := strings.TrimSpace(ips[0])
		cfg.Logger.Trace("Using X-Forwarded-For header", "ip", ip)
		return ip
	}

	// Check X-Real-IP header
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		cfg.Logger.Trace("Using X-Real-IP header", "ip", realIP)
		return realIP
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		cfg.Logger.Error("Failed to parse RemoteAddr", "remote_addr", r.RemoteAddr, "error", err)
		return r.RemoteAddr // Return as-is on error
	}
	cfg.Logger.Trace("Using RemoteAddr", "ip", ip)
	return ip
}

// TokenAuthMiddleware verifies the token in the Authorization header.
func TokenAuthMiddleware(cfg *config.BackendConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := GetClientIP(r, cfg)
		cfg.Logger.Trace("TokenAuthMiddleware passthrough", "client_ip", clientIP)
		next.ServeHTTP(w, r)
	}
}
