package panelsettings

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"exodus/internal/config"
)

func TestDocsSwaggerHandlerServesSwaggerUI(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.Docs.IsEnabled = true
	cfg.Docs.SwaggerPath = "/docs"
	cfg.Panel.BasePath = "/exodus_path"

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
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
	if !strings.Contains(body, `/exodus_path/docs/openapi.json`) {
		t.Fatalf("expected base-path-aware openapi URL, got %s", body)
	}
}

func TestDocsOpenAPIHandlerDisabled(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.Docs.IsEnabled = false

	req := httptest.NewRequest(http.MethodGet, "/docs/openapi.json", nil)
	rec := httptest.NewRecorder()

	DocsOpenAPIHandler(cfg)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status got %d, want %d", rec.Code, http.StatusNotFound)
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
	if _, exists := spec.Paths["/api/exodus-settings"]; !exists {
		t.Fatal("expected /api/exodus-settings path")
	}
	for path := range spec.Paths {
		if strings.Contains(path, "remnawave-settings") || strings.Contains(path, "ip-control") || strings.Contains(path, "torrent-blocker") {
			t.Fatalf("unexpected unsupported path in openapi spec: %s", path)
		}
	}
}
