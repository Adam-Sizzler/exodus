package panelsettings

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"exodus/internal/config"
	"exodus/internal/security"
)

func TestDocsSwaggerHandlerServesSwaggerUI(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.Panel.BasePath = "/exodus_path"

	req := httptest.NewRequest(http.MethodGet, "/api/backend-tools/swagger", nil)
	rec := httptest.NewRecorder()

	DocsSwaggerHandler(cfg)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status got %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("content type got %q, want html", rec.Header().Get("Content-Type"))
	}
	if !strings.Contains(body, "SwaggerUIBundle") {
		t.Fatalf("expected Swagger UI bundle in response, got %s", body)
	}
	if !strings.Contains(body, `/exodus_path/api/backend-tools/swagger/openapi.json`) {
		t.Fatalf("expected base-path-aware openapi URL, got %s", body)
	}
}

func TestDocsOpenAPIHandlerServesExodusSpec(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.Docs.IsEnabled = true

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	rec := httptest.NewRecorder()

	DocsOpenAPIHandler(cfg)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status got %d, want %d", rec.Code, http.StatusOK)
	}

	var spec struct {
		Info struct {
			Title string `json:"title"`
		} `json:"info"`
		Paths map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("openapi response is invalid JSON: %v", err)
	}
	if spec.Info.Title != "Exodus API" {
		t.Fatalf("title got %q, want Exodus API", spec.Info.Title)
	}
	if _, exists := spec.Paths["/exodus-settings"]; !exists {
		if _, exists2 := spec.Paths["/api/exodus-settings"]; !exists2 {
			t.Fatal("expected /exodus-settings or /api/exodus-settings path")
		}
	}
	for path := range spec.Paths {
		if strings.Contains(path, "ip-control") || strings.Contains(path, "torrent-blocker") {
			t.Fatalf("unexpected unsupported path in openapi spec: %s", path)
		}
	}
}

func TestToolsAuthMiddlewareWithOTT(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.JWT.AuthSecret = "test_auth_secret_1234567890123456"
	cfg.Panel.BasePath = "/exodus_path"

	ott, err := security.SignOttJWT(cfg.JWT.AuthSecret)
	if err != nil {
		t.Fatalf("failed to sign OTT: %v", err)
	}

	target := "/exodus_path/api/backend-tools/queues?ott=" + ott
	req := httptest.NewRequest(http.MethodGet, target, nil)
	rec := httptest.NewRecorder()

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	ToolsAuthMiddleware(cfg)(dummyHandler).ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected redirect 302 StatusFound, got %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "/exodus_path/api/backend-tools/queues" {
		t.Fatalf("expected redirect Location '/exodus_path/api/backend-tools/queues', got %q", location)
	}

	cookies := rec.Result().Cookies()
	var toolsCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == BackendToolsAuthCookieName {
			toolsCookie = c
			break
		}
	}
	if toolsCookie == nil {
		t.Fatal("expected tools session cookie to be set")
	}

	// Now test subsequent request with cookie
	req2 := httptest.NewRequest(http.MethodGet, "/exodus_path/api/backend-tools/queues", nil)
	req2.AddCookie(toolsCookie)
	rec2 := httptest.NewRecorder()

	ToolsAuthMiddleware(cfg)(dummyHandler).ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200 OK with session cookie, got %d", rec2.Code)
	}
}
