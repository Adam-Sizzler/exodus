package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"exodus/internal/config"
	"exodus/internal/security"
)

func TestAPITokenFullScopeEnforcement(t *testing.T) {
	appSecret := "test-secret-32-chars-long-abcdefgh"
	cfg := &config.BackendConfig{}
	cfg.JWT.AuthSecret = appSecret
	cfg.Backend.BasePath = "/exodus_path"

	// Create test API tokens
	// 1. users:read token
	usersReadToken, _, err := security.SignAuthJWTWithLifetime(appSecret, "read-only-token", "token-uuid-1", "API", 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}

	// 2. nodes:read token
	nodesReadToken, _, err := security.SignAuthJWTWithLifetime(appSecret, "nodes-token", "token-uuid-2", "API", 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}

	// 3. wildcard * token
	wildcardToken, _, err := security.SignAuthJWTWithLifetime(appSecret, "admin-api-token", "token-uuid-3", "API", 24*time.Hour)
	if err != nil {
		t.Fatalf("failed to sign jwt: %v", err)
	}

	// Mock resolver / middleware testing
	mockHandler := func(scopes []string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			principal := &AuthPrincipal{
				AdminUUID: "token-uuid",
				Username:  "test-token",
				Role:      "API",
				TokenType: "jwt_api_token",
				Scopes:    scopes,
			}

			if !requireAPITokenScope(principal, r, cfg) {
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden: insufficient API token scope"})
				return
			}

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		})
	}

	// Test 1: users:read token accessing GET /exodus_path/api/users -> 200 OK
	h1 := mockHandler([]string{"users:read"})
	rec1 := httptest.NewRecorder()
	req1 := httptest.NewRequest(http.MethodGet, "/exodus_path/api/users", nil)
	req1.Header.Set("Authorization", "Bearer "+usersReadToken)
	h1.ServeHTTP(rec1, req1)
	if rec1.Code != http.StatusOK {
		t.Errorf("expected 200 for users:read on GET /exodus_path/api/users, got %d", rec1.Code)
	}

	// Test 2: users:read token accessing POST /exodus_path/api/users -> 403 Forbidden
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/exodus_path/api/users", nil)
	req2.Header.Set("Authorization", "Bearer "+usersReadToken)
	h1.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusForbidden {
		t.Errorf("expected 403 for users:read on POST /exodus_path/api/users, got %d", rec2.Code)
	}

	// Test 3: users:read token accessing GET /exodus_path/api/nodes -> 403 Forbidden
	rec3 := httptest.NewRecorder()
	req3 := httptest.NewRequest(http.MethodGet, "/exodus_path/api/nodes", nil)
	req3.Header.Set("Authorization", "Bearer "+usersReadToken)
	h1.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusForbidden {
		t.Errorf("expected 403 for users:read on GET /exodus_path/api/nodes, got %d", rec3.Code)
	}

	// Test 4: nodes:read token accessing GET /exodus_path/api/nodes -> 200 OK
	h2 := mockHandler([]string{"nodes:read"})
	rec4 := httptest.NewRecorder()
	req4 := httptest.NewRequest(http.MethodGet, "/exodus_path/api/nodes", nil)
	req4.Header.Set("Authorization", "Bearer "+nodesReadToken)
	h2.ServeHTTP(rec4, req4)
	if rec4.Code != http.StatusOK {
		t.Errorf("expected 200 for nodes:read on GET /exodus_path/api/nodes, got %d", rec4.Code)
	}

	// Test 5: wildcard * token accessing GET /exodus_path/api/nodes, POST /exodus_path/api/users, DELETE /exodus_path/api/hosts/1 -> 200 OK
	h3 := mockHandler([]string{"*"})
	rec5 := httptest.NewRecorder()
	req5 := httptest.NewRequest(http.MethodDelete, "/exodus_path/api/hosts/1", nil)
	req5.Header.Set("Authorization", "Bearer "+wildcardToken)
	h3.ServeHTTP(rec5, req5)
	if rec5.Code != http.StatusOK {
		t.Errorf("expected 200 for * on DELETE /exodus_path/api/hosts/1, got %d", rec5.Code)
	}
}
