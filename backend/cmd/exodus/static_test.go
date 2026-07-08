package exodus

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderPanelIndexPrefixesStaticAssetsWithBasePath(t *testing.T) {
	input := `<html><head>
<script src="/assets/index.js"></script>
<link href="/assets/index.css" rel="stylesheet">
<link rel="icon" href="/favicons/logo.svg">
<link rel="manifest" href="/site.webmanifest">
<link rel="apple-touch-startup-image" href="%BASE_URL%/splash_screens/start.png">
</head><body></body></html>`

	page := renderPanelIndex(input, "/panel/", "/panel")

	expected := []string{
		`<base href="/panel/" />`,
		`basePath:"/panel"`,
		`src="/panel/assets/index.js"`,
		`href="/panel/assets/index.css"`,
		`href="/panel/favicons/logo.svg"`,
		`href="/panel/site.webmanifest"`,
		`href="/panel/splash_screens/start.png"`,
	}
	for _, value := range expected {
		if !strings.Contains(page, value) {
			t.Fatalf("expected rendered index to contain %q, got:\n%s", value, page)
		}
	}
}

func TestRenderPanelIndexKeepsRootStaticAssetsAtRoot(t *testing.T) {
	input := `<html><head><script src="%BASE_URL%/assets/index.js"></script></head></html>`

	page := renderPanelIndex(input, "/", "/")

	if !strings.Contains(page, `src="/assets/index.js"`) {
		t.Fatalf("expected root asset path, got:\n%s", page)
	}
	if strings.Contains(page, `%BASE_URL%`) {
		t.Fatalf("expected %%BASE_URL%% placeholder to be removed, got:\n%s", page)
	}
}

func TestServePanelStaticFileServesExistingAsset(t *testing.T) {
	uiDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(uiDir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uiDir, "assets", "index.js"), []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/assets/index.js", nil)

	handled := servePanelStaticFile(recorder, request, http.FileServer(http.Dir(uiDir)), uiDir, request.URL.Path)
	if !handled {
		t.Fatal("expected static request to be handled")
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recorder.Code)
	}
	if body := recorder.Body.String(); body != "console.log('ok')" {
		t.Fatalf("unexpected static body: %q", body)
	}
}

func TestPanelHandlerDoesNotServeRootAssetsWhenBasePathIsConfigured(t *testing.T) {
	uiDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(uiDir, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uiDir, "index.html"), []byte(`<html><head></head><body></body></html>`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uiDir, "assets", "index.js"), []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}

	staticFS := http.FileServer(http.Dir(uiDir))
	handler := panelRequestHandler("/panel/", uiDir, staticFS, http.NotFoundHandler(), http.NotFoundHandler(), "/docs", "/scalar")

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/assets/index.js", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected root asset request to be rejected when base path is configured, got %d", recorder.Code)
	}
}

func TestServePanelStaticFileReturnsNotFoundForMissingAsset(t *testing.T) {
	uiDir := t.TempDir()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panel/assets/missing.js", nil)

	handled := servePanelStaticFile(recorder, request, http.FileServer(http.Dir(uiDir)), uiDir, "assets/missing.js")
	if !handled {
		t.Fatal("expected missing static request to be handled")
	}
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestPanelHandlerRoutesDocsWithTrailingSlashToAPIWithoutRedirect(t *testing.T) {
	uiDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(uiDir, "index.html"), []byte(`<html><head></head><body>spa</body></html>`), 0o644); err != nil {
		t.Fatal(err)
	}

	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/docs" {
			t.Fatalf("API path got %q, want /docs", r.URL.Path)
		}
		_, _ = w.Write([]byte("docs api"))
	})

	handler := panelRequestHandler("/panel/", uiDir, http.FileServer(http.Dir(uiDir)), apiHandler, http.NotFoundHandler(), "/docs", "/scalar")
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panel/docs/", nil)

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status got %d, want %d", recorder.Code, http.StatusOK)
	}
	if location := recorder.Header().Get("Location"); location != "" {
		t.Fatalf("expected no redirect Location header, got %q", location)
	}
	if recorder.Body.String() != "docs api" {
		t.Fatalf("expected API docs response, got %q", recorder.Body.String())
	}
}

func TestDocsAPIRequestPathMatchesWithAndWithoutTrailingSlash(t *testing.T) {
	tests := map[string]string{
		"docs":              "/docs",
		"docs/":             "/docs",
		"docs/openapi.json": "/docs/openapi.json",
		"scalar":            "/scalar",
		"scalar/":           "/scalar",
	}

	for input, expected := range tests {
		actual, ok := docsAPIRequestPath(input, "/docs", "/scalar")
		if !ok {
			t.Fatalf("expected %q to be recognized as docs path", input)
		}
		if actual != expected {
			t.Fatalf("for %q got %q, want %q", input, actual, expected)
		}
	}

	if actual, ok := docsAPIRequestPath("dashboard", "/docs", "/scalar"); ok {
		t.Fatalf("expected dashboard to stay in SPA, got %q", actual)
	}
}
