package middleware

import (
	"net/http"
	"time"

	"v2ray-stat/backend/config"
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

		cfg.Logger.Trace("HTTP request",
			"component", component,
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"status", statusCode,
			"bytes", lrw.bytes,
			"duration_ms", time.Since(start).Milliseconds(),
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
		)
	})
}
