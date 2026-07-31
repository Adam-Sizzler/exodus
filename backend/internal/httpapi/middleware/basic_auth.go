package middleware

import (
	"crypto/subtle"
	"io"
	"net/http"
	"strings"

	"exodus/internal/config"
)

// WithMetricsBasicAuth returns an HTTP middleware enforcing Basic Auth for metrics requests.
// Health check paths (/health, /api/health) are excluded from Basic Auth.
func WithMetricsBasicAuth(cfg *config.BackendConfig, next http.Handler) http.Handler {
	if cfg == nil {
		return next
	}

	expectedUser := strings.TrimSpace(cfg.Metrics.User)
	expectedPass := strings.TrimSpace(cfg.Metrics.Pass)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/health" || path == "/api/health" || strings.HasSuffix(path, "/health") {
			next.ServeHTTP(w, r)
			return
		}

		user, pass, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "unauthorized")
			return
		}

		userMatch := subtle.ConstantTimeCompare([]byte(user), []byte(expectedUser)) == 1
		passMatch := subtle.ConstantTimeCompare([]byte(pass), []byte(expectedPass)) == 1
		if !userMatch || !passMatch {
			w.Header().Set("WWW-Authenticate", `Basic realm="metrics"`)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, "unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
