package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
)

func TestWithMetricsBasicAuth(t *testing.T) {
	cfg := &config.BackendConfig{
		Metrics: config.MetricsConfig{
			User: "admin",
			Pass: "secret",
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := WithMetricsBasicAuth(cfg, nextHandler)

	// 1. Health endpoint without auth should succeed
	healthReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthRec := httptest.NewRecorder()
	handler.ServeHTTP(healthRec, healthReq)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for /health without auth, got %d", healthRec.Code)
	}

	// 2. Metrics endpoint without auth should fail with 401
	metricsReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsRec := httptest.NewRecorder()
	handler.ServeHTTP(metricsRec, metricsReq)
	if metricsRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for /metrics without auth, got %d", metricsRec.Code)
	}

	// 3. Metrics endpoint with wrong auth should fail with 401
	wrongReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	wrongReq.SetBasicAuth("admin", "wrongpass")
	wrongRec := httptest.NewRecorder()
	handler.ServeHTTP(wrongRec, wrongReq)
	if wrongRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for wrong credentials, got %d", wrongRec.Code)
	}

	// 4. Metrics endpoint with correct auth should succeed
	correctReq := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	correctReq.SetBasicAuth("admin", "secret")
	correctRec := httptest.NewRecorder()
	handler.ServeHTTP(correctRec, correctReq)
	if correctRec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK for correct credentials, got %d", correctRec.Code)
	}
}
