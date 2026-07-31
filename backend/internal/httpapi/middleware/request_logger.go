package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	n, err := lrw.ResponseWriter.Write(b)
	lrw.bytes += n
	return n, err
}

func WithRequestLogging(cfg *config.BackendConfig, component string, next http.Handler) http.Handler {
	if cfg == nil || cfg.Logger == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		statusCode := lrw.statusCode
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		duration := time.Since(start)
		durationMs := duration.Milliseconds()
		if durationMs == 0 && duration > 0 {
			durationMs = 1
		}

		if shouldSkipAccessLog(cfg, r.URL.Path) {
			return
		}

		role := logger.RoleAPI
		if strings.EqualFold(component, "metrics") {
			role = logger.RoleScheduler
		}
		serviceLogger := cfg.Logger.RoleService(role, logger.ServiceHTTP)

		msg := fmt.Sprintf("%s %s %d %dms", r.Method, r.URL.Path, statusCode, durationMs)
		if cfg.Log.IsHTTPLoggingEnabled {
			serviceLogger.Info(msg,
				"component", component,
				"method", r.Method,
				"path", r.URL.Path,
				"status", statusCode,
				"bytes", lrw.bytes,
				"duration_ms", durationMs,
			)
		} else {
			serviceLogger.Debug(msg,
				"component", component,
				"method", r.Method,
				"path", r.URL.Path,
				"status", statusCode,
				"bytes", lrw.bytes,
				"duration_ms", durationMs,
			)
		}
		serviceLogger.Trace("HTTP request details",
			"component", component,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", statusCode,
			"bytes", lrw.bytes,
			"duration_ms", durationMs,
			"duration_us", duration.Microseconds(),
			"client_ip", GetClientIP(r, cfg),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"x_exodus_real_ip", r.Header.Get(ExodusRealIPHeader),
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
		)
	})
}

func shouldSkipAccessLog(cfg *config.BackendConfig, path string) bool {
	if cfg == nil {
		return false
	}
	if strings.EqualFold(strings.TrimSpace(cfg.Log.NodeEnv), "development") {
		return false
	}
	path = strings.ToLower(strings.TrimSpace(path))
	if path == "" {
		return false
	}
	if path == "/favicon.ico" || path == "/site.webmanifest" || path == "/robots.txt" {
		return true
	}
	if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/locales/") {
		return true
	}
	staticExts := []string{".css", ".js", ".map", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp", ".woff", ".woff2", ".ttf", ".eot"}
	for _, ext := range staticExts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}
