package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
	"exodus/internal/httpapi/health"
)

func TestMetricsMuxEndpoints(t *testing.T) {
	t.Run("root base path", func(t *testing.T) {
		cfg := &config.BackendConfig{}
		cfg.Backend.BasePath = "/"
		prefix := cfg.Backend.Trimmed()

		mux := http.NewServeMux()
		mux.HandleFunc(prefix+"/health", health.HealthHandler())

		// Root endpoint should succeed
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /health status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("custom base path", func(t *testing.T) {
		cfg := &config.BackendConfig{}
		cfg.Backend.BasePath = "/exodus_path/"
		prefix := cfg.Backend.Trimmed()

		mux := http.NewServeMux()
		mux.HandleFunc(prefix+"/health", health.HealthHandler())

		// Custom path should succeed
		req := httptest.NewRequest(http.MethodGet, "/exodus_path/health", nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("GET /exodus_path/health status = %d, want %d", rec.Code, http.StatusOK)
		}

		// Root path must return 404 (strictly no dual path!)
		reqRoot := httptest.NewRequest(http.MethodGet, "/health", nil)
		recRoot := httptest.NewRecorder()
		mux.ServeHTTP(recRoot, reqRoot)
		if recRoot.Code != http.StatusNotFound {
			t.Fatalf("GET /health with custom BasePath got status = %d, want 404 Not Found", recRoot.Code)
		}
	})
}
