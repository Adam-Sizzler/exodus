package middleware

import (
	"net"
	"net/http"
	"strings"

	"exodus/internal/config"
)

// GetClientIP retrieves the client IP address from an HTTP request.
// It checks X-Forwarded-For and X-Real-IP headers before falling back to
// RemoteAddr. Note: these headers are attacker-controlled unless the request
// is verified to come from a trusted proxy first.
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
		return r.RemoteAddr
	}
	cfg.Logger.Trace("Using RemoteAddr", "ip", ip)
	return ip
}
