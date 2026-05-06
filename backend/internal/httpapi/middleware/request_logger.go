package middleware

import (
	"net/http"
	"time"

	"exodus/internal/config"
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
		durationUs := duration.Microseconds()

		// Keep request timing visible even when debug is disabled.
		cfg.Logger.Info("HTTP request",
			"component", component,
			"method", r.Method,
			"path", r.URL.Path,
			"status", statusCode,
			"duration_ms", durationMs,
			"duration_us", durationUs,
		)
		cfg.Logger.Debug("HTTP request debug",
			"component", component,
			"method", r.Method,
			"path", r.URL.Path,
			"status", statusCode,
			"bytes", lrw.bytes,
			"duration_ms", durationMs,
			"duration_us", durationUs,
		)
		cfg.Logger.Trace("HTTP request details",
			"component", component,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", statusCode,
			"bytes", lrw.bytes,
			"duration_ms", durationMs,
			"duration_us", durationUs,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
		)
	})
}
