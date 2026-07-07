package panelsettings

import (
	"encoding/json"
	"fmt"
	"strings"

	"exodus/internal/config"
)

type apiTokenEndpointScope struct {
	Key         string `json:"key"`
	Kind        string `json:"kind"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type apiTokenResourceScopes struct {
	Resource       string                  `json:"resource"`
	ResourceScopes []string                `json:"resourceScopes"`
	Endpoints      []apiTokenEndpointScope `json:"endpoints"`
}

func buildAPITokenScopes(_ *config.BackendConfig) []apiTokenResourceScopes {
	resources := []apiTokenResourceScopes{
		scopeResource("system", []apiTokenEndpointScope{
			readScope("system:stats", "GET", "/api/system/stats", "Read system statistics"),
			readScope("system:metadata", "GET", "/api/system/metadata", "Read build metadata"),
			readScope("system:health", "GET", "/api/system/health", "Read runtime health"),
			readScope("system:recap", "GET", "/api/system/stats/recap", "Read recap information"),
		}),
		scopeResource("users", []apiTokenEndpointScope{
			readScope("users:list", "GET", "/api/users", "List users"),
			readScope("users:get", "GET", "/api/users/{uuid}", "Read user"),
			readScope("users:by-username", "GET", "/api/sub/{username}", "Read user subscription by username"),
			writeScope("users:create", "POST", "/api/users", "Create user"),
			writeScope("users:update", "PATCH", "/api/users/{uuid}", "Update user"),
			writeScope("users:delete", "DELETE", "/api/users/{uuid}", "Delete user"),
		}),
		scopeResource("nodes", []apiTokenEndpointScope{
			readScope("nodes:list", "GET", "/api/nodes", "List nodes"),
			readScope("nodes:get", "GET", "/api/nodes/{uuid}", "Read node"),
			writeScope("nodes:create", "POST", "/api/nodes", "Create node"),
			writeScope("nodes:update", "PATCH", "/api/nodes/{uuid}", "Update node"),
			writeScope("nodes:delete", "DELETE", "/api/nodes/{uuid}", "Delete node"),
		}),
		scopeResource("hosts", []apiTokenEndpointScope{
			readScope("hosts:list", "GET", "/api/hosts", "List hosts"),
			readScope("hosts:get", "GET", "/api/hosts/{uuid}", "Read host"),
			writeScope("hosts:create", "POST", "/api/hosts", "Create host"),
			writeScope("hosts:update", "PATCH", "/api/hosts/{uuid}", "Update host"),
			writeScope("hosts:delete", "DELETE", "/api/hosts/{uuid}", "Delete host"),
		}),
		scopeResource("subscription-connections", []apiTokenEndpointScope{
			readScope("subscription-connections:list", "GET", "/api/subscription-connections", "List subscription connections"),
			readScope("subscription-connections:get", "GET", "/api/subscription-connections/{uuid}", "Read subscription connection"),
			writeScope("subscription-connections:create", "POST", "/api/subscription-connections", "Create subscription connection"),
			writeScope("subscription-connections:update", "PATCH", "/api/subscription-connections/{uuid}", "Update subscription connection"),
			writeScope("subscription-connections:delete", "DELETE", "/api/subscription-connections/{uuid}", "Delete subscription connection"),
		}),
		scopeResource("config-profiles", []apiTokenEndpointScope{
			readScope("config-profiles:list", "GET", "/api/config-profiles", "List config profiles"),
			readScope("config-profiles:get", "GET", "/api/config-profiles/{uuid}", "Read config profile"),
			writeScope("config-profiles:create", "POST", "/api/config-profiles", "Create config profile"),
			writeScope("config-profiles:update", "PATCH", "/api/config-profiles/{uuid}", "Update config profile"),
			writeScope("config-profiles:delete", "DELETE", "/api/config-profiles/{uuid}", "Delete config profile"),
		}),
		scopeResource("subscription-page-configs", []apiTokenEndpointScope{
			readScope("subscription-page-configs:list", "GET", "/api/subscription-page-configs", "List subscription page configs"),
			readScope("subscription-page-configs:get", "GET", "/api/subscription-page-configs/{uuid}", "Read subscription page config"),
			writeScope("subscription-page-configs:create", "POST", "/api/subscription-page-configs", "Create subscription page config"),
			writeScope("subscription-page-configs:update", "PATCH", "/api/subscription-page-configs/{uuid}", "Update subscription page config"),
			writeScope("subscription-page-configs:delete", "DELETE", "/api/subscription-page-configs/{uuid}", "Delete subscription page config"),
		}),
		scopeResource("subscriptions", []apiTokenEndpointScope{
			readScope("subscriptions:subpage-config", "GET", "/api/sub/{shortUuid}", "Read subscription page"),
			readScope("subscriptions:get", "GET", "/api/subscriptions/{uuid}", "Read subscription"),
		}),
		scopeResource("subscription-settings", []apiTokenEndpointScope{
			readScope("subscription-settings:read", "GET", "/api/subscription-settings", "Read subscription settings"),
			writeScope("subscription-settings:write", "PATCH", "/api/subscription-settings", "Update subscription settings"),
		}),
		scopeResource("subscription-template", []apiTokenEndpointScope{
			readScope("subscription-template:list", "GET", "/api/subscription-templates", "List subscription templates"),
			readScope("subscription-template:get", "GET", "/api/subscription-templates/{uuid}", "Read subscription template"),
			writeScope("subscription-template:create", "POST", "/api/subscription-templates", "Create subscription template"),
			writeScope("subscription-template:update", "PATCH", "/api/subscription-templates/{uuid}", "Update subscription template"),
			writeScope("subscription-template:delete", "DELETE", "/api/subscription-templates/{uuid}", "Delete subscription template"),
		}),
		scopeResource("subscription-request-history", []apiTokenEndpointScope{
			readScope("subscription-request-history:read", "GET", "/api/subscription-request-history", "Read subscription request history"),
		}),
		scopeResource("metadata", []apiTokenEndpointScope{
			readScope("metadata:get-user", "GET", "/api/metadata/user/{uuid}", "Read user metadata"),
			writeScope("metadata:upsert-user", "PUT", "/api/metadata/user/{uuid}", "Update user metadata"),
			readScope("metadata:get-node", "GET", "/api/metadata/node/{uuid}", "Read node metadata"),
			writeScope("metadata:upsert-node", "PUT", "/api/metadata/node/{uuid}", "Update node metadata"),
		}),
		scopeResource("node-plugins", []apiTokenEndpointScope{
			readScope("node-plugins:list", "GET", "/api/node-plugins", "List node plugins"),
			readScope("node-plugins:get", "GET", "/api/node-plugins/{uuid}", "Read node plugin"),
			writeScope("node-plugins:create", "POST", "/api/node-plugins", "Create node plugin"),
			writeScope("node-plugins:update", "PATCH", "/api/node-plugins/{uuid}", "Update node plugin"),
			writeScope("node-plugins:delete", "DELETE", "/api/node-plugins/{uuid}", "Delete node plugin"),
		}),
		scopeResource("hwid-user-devices", []apiTokenEndpointScope{
			readScope("hwid-user-devices:list", "GET", "/api/hwid/devices", "List HWID devices"),
			readScope("hwid-user-devices:stats", "GET", "/api/hwid/devices/stats", "Read HWID stats"),
			writeScope("hwid-user-devices:create", "POST", "/api/hwid/devices", "Create HWID device"),
			writeScope("hwid-user-devices:delete", "POST", "/api/hwid/devices/delete", "Delete HWID device"),
		}),
		scopeResource("bandwidth-stats", []apiTokenEndpointScope{
			readScope("bandwidth-stats:nodes", "GET", "/api/bandwidth-stats/nodes", "Read node bandwidth stats"),
			readScope("bandwidth-stats:users", "GET", "/api/bandwidth-stats/users", "Read user bandwidth stats"),
		}),
		scopeResource("srs-lists", []apiTokenEndpointScope{
			readScope("srs-lists:list", "GET", "/api/srs-lists", "List SRS lists"),
			writeScope("srs-lists:create", "POST", "/api/srs-lists", "Create SRS list"),
			writeScope("srs-lists:update", "PATCH", "/api/srs-lists/{uuid}", "Update SRS list"),
			writeScope("srs-lists:delete", "DELETE", "/api/srs-lists/{uuid}", "Delete SRS list"),
		}),
		scopeResource("internal-squads", []apiTokenEndpointScope{
			readScope("internal-squads:read", "GET", "/api/internal-squads", "Read internal squads"),
			writeScope("internal-squads:write", "PATCH", "/api/internal-squads/{uuid}", "Update internal squads"),
		}),
		scopeResource("external-squads", []apiTokenEndpointScope{
			readScope("external-squads:read", "GET", "/api/external-squads", "Read external squads"),
			writeScope("external-squads:write", "PATCH", "/api/external-squads/{uuid}", "Update external squads"),
		}),
		scopeResource("infra-billing", []apiTokenEndpointScope{
			readScope("infra-billing:read", "GET", "/api/infra-billing/providers", "Read infra billing"),
			writeScope("infra-billing:write", "PATCH", "/api/infra-billing/providers/{uuid}", "Update infra billing"),
		}),
		scopeResource("snippets", []apiTokenEndpointScope{
			readScope("snippets:read", "GET", "/api/snippets", "Read config snippets"),
			writeScope("snippets:write", "POST", "/api/snippets", "Create config snippet"),
		}),
		scopeResource("keygen", []apiTokenEndpointScope{
			readScope("keygen:read", "GET", "/api/keygen", "Generate key material"),
		}),
	}

	return resources
}

func scopeResource(resource string, endpoints []apiTokenEndpointScope) apiTokenResourceScopes {
	return apiTokenResourceScopes{
		Resource:       resource,
		ResourceScopes: []string{resource + ":read", resource + ":write", resource + ":*"},
		Endpoints:      endpoints,
	}
}

func readScope(key, method, path, description string) apiTokenEndpointScope {
	return apiTokenEndpointScope{Key: key, Kind: "read", Method: method, Path: path, Description: description}
}

func writeScope(key, method, path, description string) apiTokenEndpointScope {
	return apiTokenEndpointScope{Key: key, Kind: "write", Method: method, Path: path, Description: description}
}

func normalizeAPITokenScopes(scopes []string) []string {
	out := make([]string, 0, len(scopes))
	seen := map[string]struct{}{}
	for _, scope := range scopes {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, ok := seen[scope]; ok {
			continue
		}
		seen[scope] = struct{}{}
		out = append(out, scope)
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}

func parseAPITokenScopes(raw string) []string {
	var scopes []string
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &scopes); err != nil {
		return []string{"*"}
	}
	return normalizeAPITokenScopes(scopes)
}

func postgresTextArrayLiteral(items []string) string {
	items = normalizeAPITokenScopes(items)
	quoted := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.ReplaceAll(item, `\`, `\\`)
		item = strings.ReplaceAll(item, `"`, `\"`)
		quoted = append(quoted, fmt.Sprintf(`"%s"`, item))
	}
	return "{" + strings.Join(quoted, ",") + "}"
}
