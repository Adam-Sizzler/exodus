package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"exodus/internal/config"
)

const ExodusRealIPHeader = "X-Exodus-Real-IP"

type clientIPContextKey struct{}

type clientIPCandidate struct {
	value  string
	ip     net.IP
	source string
}

// WithClientIP resolves the client IP once per request and stores it in the
// request context. Handlers should use GetClientIP instead of reading
// RemoteAddr directly when they need the end-user address.
func WithClientIP(cfg *config.BackendConfig, next http.Handler) http.Handler {
	if next == nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := ResolveClientIP(r, cfg)
		ctx := context.WithValue(r.Context(), clientIPContextKey{}, clientIP)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClientIP returns the request client IP resolved by WithClientIP, or
// resolves it on demand when the middleware has not run for this request.
func GetClientIP(r *http.Request, cfg *config.BackendConfig) string {
	if r == nil {
		return ""
	}
	if value, ok := r.Context().Value(clientIPContextKey{}).(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	return ResolveClientIP(r, cfg)
}

// ResolveClientIP retrieves the client IP address from an HTTP request. It
// prefers the Exodus explicit real-IP header, then common proxy/CDN headers,
// then RemoteAddr. If a public address is present in proxy headers, it is
// preferred over private Docker/proxy addresses from the same chain.
func ResolveClientIP(r *http.Request, cfg *config.BackendConfig) string {
	if r == nil {
		return ""
	}

	candidates := make([]clientIPCandidate, 0, 8)
	for _, header := range []string{
		ExodusRealIPHeader,
		"CF-Connecting-IP",
		"True-Client-IP",
		"X-Forwarded-For",
		"X-Real-IP",
	} {
		candidates = appendHeaderIPCandidates(candidates, header, r.Header.Values(header))
	}
	if candidate, ok := parseIPCandidate(r.RemoteAddr, "RemoteAddr"); ok {
		candidates = append(candidates, candidate)
	}

	selected := selectClientIPCandidate(candidates)
	if selected.value == "" {
		selected.value = "0.0.0.0"
	}
	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Debug("Resolved client IP address", "client_ip", selected.value, "source", selected.source, "remote_addr", r.RemoteAddr)
	}
	return selected.value
}

func appendHeaderIPCandidates(candidates []clientIPCandidate, header string, values []string) []clientIPCandidate {
	for _, value := range values {
		for _, item := range strings.Split(value, ",") {
			if candidate, ok := parseIPCandidate(item, header); ok {
				candidates = append(candidates, candidate)
			}
		}
	}
	return candidates
}

func parseIPCandidate(value, source string) (clientIPCandidate, bool) {
	trimmed := strings.Trim(strings.TrimSpace(value), "\"'")
	if trimmed == "" {
		return clientIPCandidate{}, false
	}

	if host, _, err := net.SplitHostPort(trimmed); err == nil {
		trimmed = host
	}
	trimmed = strings.Trim(trimmed, "[]")
	if strings.HasPrefix(strings.ToLower(trimmed), "::ffff:") {
		trimmed = trimmed[7:]
	}
	if zone := strings.LastIndex(trimmed, "%"); zone > -1 {
		trimmed = trimmed[:zone]
	}
	if ip := net.ParseIP(trimmed); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			ip = ip4
		}
		return clientIPCandidate{value: ip.String(), ip: ip, source: source}, true
	}
	return clientIPCandidate{}, false
}

func selectClientIPCandidate(candidates []clientIPCandidate) clientIPCandidate {
	if len(candidates) == 0 {
		return clientIPCandidate{}
	}
	for _, candidate := range candidates {
		if !strings.EqualFold(candidate.source, "RemoteAddr") && isPublicClientIP(candidate.ip) {
			return candidate
		}
	}
	return candidates[0]
}

func isPublicClientIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsLoopback() &&
		!ip.IsLinkLocalUnicast() &&
		!ip.IsLinkLocalMulticast() &&
		!ip.IsUnspecified()
}
