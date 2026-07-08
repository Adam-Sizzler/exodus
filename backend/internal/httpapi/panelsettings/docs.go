package panelsettings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"exodus/internal/config"
	"exodus/internal/constant"
	"exodus/internal/httpapi/shared"
)

//go:embed openapi.json
var openapiJSON []byte

var (
	exodusOpenAPIOnce  sync.Once
	exodusOpenAPIBytes []byte
	exodusOpenAPIErr   error
)

// buildDocsResponse builds the "docs" field for GET /api/tokens response.
// When IS_DOCS_ENABLED=false all paths are null — the frontend hides the buttons.
// When enabled, paths are prefixed with APP_PATH so the links work under any base path.
func buildDocsResponse(cfg *config.BackendConfig) map[string]any {
	if !cfg.Docs.IsEnabled {
		return map[string]any{
			"enabled":     false,
			"scalarPath":  nil,
			"swaggerPath": nil,
		}
	}

	basePath := strings.TrimRight(cfg.Panel.BasePath, "/")

	return map[string]any{
		"enabled":     true,
		"scalarPath":  fmt.Sprintf("%s%s", basePath, cfg.Docs.ScalarPath),
		"swaggerPath": fmt.Sprintf("%s%s", basePath, cfg.Docs.SwaggerPath),
	}
}

// DocsScalarHandler serves the Scalar API reference UI.
// Route: GET /scalar (or whatever SCALAR_PATH is set to, registered in router.go).
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

// DocsSwaggerHandler serves a Swagger UI page compatible with SWAGGER_PATH.
// The raw OpenAPI spec remains available at SWAGGER_PATH/openapi.json.
func DocsSwaggerHandler(cfg *config.BackendConfig) http.HandlerFunc {
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
		specURL := fmt.Sprintf("%s%s/openapi.json", basePath, cfg.Docs.SwaggerPath)

		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Exodus API Schema</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; background: #0f1117; }
    .swagger-ui .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: %q,
      dom_id: '#swagger-ui',
      deepLinking: true,
      persistAuthorization: true,
      layout: 'BaseLayout'
    })
  </script>
</body>
</html>`, specURL)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(html))
	}
}

// DocsOpenAPIHandler serves the raw OpenAPI JSON spec.
// Route: GET /scalar/openapi.json or /docs/openapi.json.
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
		spec, err := exodusOpenAPISpec()
		if err != nil {
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to prepare openapi spec")
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(spec)
	}
}

func exodusOpenAPISpec() ([]byte, error) {
	exodusOpenAPIOnce.Do(func() {
		exodusOpenAPIBytes, exodusOpenAPIErr = buildExodusOpenAPISpec()
	})
	return exodusOpenAPIBytes, exodusOpenAPIErr
}

func buildExodusOpenAPISpec() ([]byte, error) {
	var doc map[string]any
	if err := json.Unmarshal(openapiJSON, &doc); err != nil {
		return nil, err
	}

	if paths, ok := doc["paths"].(map[string]any); ok {
		for path := range paths {
			if strings.HasPrefix(path, "/api/ip-control/") || strings.HasPrefix(path, "/api/node-plugins/torrent-blocker") {
				delete(paths, path)
			}
		}
	}

	if info, ok := doc["info"].(map[string]any); ok {
		info["title"] = "Exodus API"
		info["description"] = "Exodus dashboard and node-management API."
		info["version"] = constant.Version
		info["license"] = map[string]any{"name": "AGPL-3.0"}
	}

	if tags, ok := doc["tags"].([]any); ok {
		cleanTags := make([]any, 0, len(tags))
		for _, rawTag := range tags {
			tag, ok := rawTag.(map[string]any)
			if !ok {
				cleanTags = append(cleanTags, rawTag)
				continue
			}
			name := strings.ToLower(fmt.Sprint(tag["name"]))
			if strings.Contains(name, "ip control") || strings.Contains(name, "torrent blocker") {
				continue
			}
			cleanTags = append(cleanTags, rawTag)
		}
		doc["tags"] = cleanTags
	}

	if components, ok := doc["components"].(map[string]any); ok {
		if schemas, ok := components["schemas"].(map[string]any); ok {
			for name := range schemas {
				lower := strings.ToLower(name)
				if strings.Contains(lower, "torrentblocker") || strings.Contains(lower, "ipcontrol") {
					delete(schemas, name)
				}
			}
		}
	}

	spec, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	replacer := strings.NewReplacer(
		"exodusSettings", "ExodusSettings",
		"Getexodus", "GetExodus",
		"Updateexodus", "UpdateExodus",
		"exodus Settings Controller", "Exodus Settings Controller",
		"exodus settings", "Exodus settings",
		"exodus Node", "Exodus Node",
		"exodus Information", "Exodus Information",
		"exodus Health", "Exodus Health",
		"exodus health", "Exodus health",
		"exodus API", "Exodus API",
	)

	return []byte(replacer.Replace(string(spec))), nil
}
