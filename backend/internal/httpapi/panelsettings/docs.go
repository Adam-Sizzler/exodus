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

//go:embed docs/swagger.json
var openapiJSON []byte

var (
	exodusOpenAPIOnce  sync.Once
	exodusOpenAPIBytes []byte
	exodusOpenAPIErr   error

	defaultControllerTags = []map[string]string{
		{"name": "Users Controller", "description": "Manage users, change their status, reset traffic, etc."},
		{"name": "Users Bulk Actions Controller", "description": "Manage users in bulk, perform mass operations, etc."},
		{"name": "HWID User Devices Controller", "description": "Manage HWID user devices, reset HWID, etc."},
		{"name": "[Public] Subscription Controller", "description": "Public subscription endpoints for clients to fetch configs and user details."},
		{"name": "[Protected] Subscriptions Controller", "description": "Protected endpoints for subscription management."},
		{"name": "Nodes Controller", "description": "Manage nodes, change their status, restart, etc."},
		{"name": "Node Plugins Controller", "description": "Manage node plugins, torrent blocker, etc."},
		{"name": "Bandwidth Stats Controller", "description": "View bandwidth statistics for users and nodes."},
		{"name": "Connections Controller", "description": "Manage connections, connection profiles, etc."},
		{"name": "Config Profiles Controller", "description": "Management of Config Profiles."},
		{"name": "Snippets Controller", "description": "Manage reusable configuration snippets for config profiles."},
		{"name": "Internal Squads Controller", "description": "Manage Internal Squads."},
		{"name": "External Squads Controller", "description": "Manage External Squads."},
		{"name": "Hosts Controller", "description": "Manage hosts, change their status, etc."},
		{"name": "Hosts Bulk Actions Controller", "description": "Manage hosts in bulk, perform mass operations, etc."},
		{"name": "Subscription Template Controller", "description": "Manage subscription templates, modify them, etc."},
		{"name": "Subscription Settings Controller", "description": "Manage subscription settings, modify them, etc."},
		{"name": "Subscription Page Configs Controller", "description": "Manage subscription page configs, modify them, etc."},
		{"name": "Subscription Request History Controller", "description": "Manage subscription request history, view logs, etc."},
		{"name": "Infra Billing Controller", "description": "Manage infrastructure billing, view invoice history, etc."},
		{"name": "System Controller", "description": "System information, stats, tools, etc."},
		{"name": "Keygen Controller", "description": "Key generation for nodes and clients."},
		{"name": "Metadata Controller", "description": "Custom metadata for users and nodes."},
		{"name": "Auth Controller", "description": "Admin authentication, login, session, and OAuth2."},
		{"name": "Passkeys Controller", "description": "WebAuthn / Passkeys management and authentication."},
		{"name": "API Tokens Controller", "description": "Manage scoped API tokens and OTT tokens for panel automation."},
		{"name": "Exodus Settings Controller", "description": "Exodus dashboard preferences, passkey, oauth2, and security policies."},
		{"name": "SRS Lists Controller", "description": "Sing-box binary rule-set lists and routing assets."},
		{"name": "Health", "description": "Backend health and liveness probe."},
	}
)

// DocsScalarHandler serves the Scalar API reference UI.
// Route: GET /scalar (or whatever SCALAR_PATH is set to, registered in router.go).
func DocsScalarHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		basePath := strings.TrimRight(cfg.Panel.BasePath, "/")
		specURL := fmt.Sprintf("%s/api/backend-tools/scalar/openapi.json", basePath)

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

func DocsSwaggerHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		basePath := strings.TrimRight(cfg.Panel.BasePath, "/")
		specURL := fmt.Sprintf("%s/api/backend-tools/swagger/openapi.json", basePath)

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

func DocsOpenAPIHandler(_ *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

	doc["tags"] = defaultControllerTags

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
