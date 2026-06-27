package panelsettings

import (
	_ "embed"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
)

//go:embed openapi.json
var openapiJSON []byte

// buildDocsResponse builds the "docs" field for GET /api/tokens response.
// When IS_DOCS_ENABLED=false all paths are null — the frontend hides the buttons.
// When enabled, paths are prefixed with APP_PATH so the links work under any base path.
func buildDocsResponse(cfg *config.BackendConfig) map[string]any {
	if !cfg.Docs.IsEnabled {
		return map[string]any{
			"isDocsEnabled": false,
			"scalarPath":    nil,
			"swaggerPath":   nil,
		}
	}

	basePath := strings.TrimRight(cfg.Panel.BasePath, "/")

	return map[string]any{
		"isDocsEnabled": true,
		"scalarPath":    fmt.Sprintf("%s%s", basePath, cfg.Docs.ScalarPath),
		"swaggerPath":   fmt.Sprintf("%s%s", basePath, cfg.Docs.SwaggerPath),
	}
}

// DocsScalarHandler serves the Scalar API reference UI.
// Route: GET /api/docs  (or whatever SCALAR_PATH is set to, registered in router.go)
// Protected by WithPanelAuth — only authenticated users can access.
func DocsScalarHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Docs.IsEnabled {
			shared.WriteJSONError(w, http.StatusNotFound, "docs are not enabled")
			return
		}
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		basePath := strings.TrimRight(cfg.Panel.BasePath, "/")
		specURL := fmt.Sprintf("%s%s/openapi.json", basePath, cfg.Docs.ScalarPath)

		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Exodus API</title>
</head>
<body>
  <script
    id="api-reference"
    data-url="%s"
    data-configuration='{
      "theme": "purple",
      "darkMode": true,
      "layout": "modern",
      "showSidebar": true,
      "hideModels": false,
      "hideDownloadButton": false,
      "hideTestRequestButton": false,
      "persistAuth": true,
      "telemetry": false,
      "defaultHttpClient": {
        "targetKey": "js",
        "clientKey": "axios"
      }
    }'
  ></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`, specURL)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	}
}

// DocsOpenAPIHandler serves the raw OpenAPI JSON spec.
// Route: GET /api/docs/openapi.json
func DocsOpenAPIHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.Docs.IsEnabled {
			shared.WriteJSONError(w, http.StatusNotFound, "docs are not enabled")
			return
		}
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openapiJSON)
	}
}
