package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
)

func TestRequireAPITokenScope(t *testing.T) {
	adminPrincipal := &AuthPrincipal{
		Role:      "ADMIN",
		TokenType: "jwt_auth",
	}

	allAccessAPIToken := &AuthPrincipal{
		Role:      "API",
		TokenType: "jwt_api_token",
		Scopes:    []string{"*"},
	}

	usersReadToken := &AuthPrincipal{
		Role:      "API",
		TokenType: "jwt_api_token",
		Scopes:    []string{"users:read"},
	}

	usersWildcardToken := &AuthPrincipal{
		Role:      "API",
		TokenType: "jwt_api_token",
		Scopes:    []string{"users:*"},
	}

	usersListToken := &AuthPrincipal{
		Role:      "API",
		TokenType: "jwt_api_token",
		Scopes:    []string{"users:list"},
	}

	reqGetUsers := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	reqGetUserByID := httptest.NewRequest(http.MethodGet, "/api/users/123", nil)
	reqPostUser := httptest.NewRequest(http.MethodPost, "/api/users", nil)
	reqDeleteUser := httptest.NewRequest(http.MethodDelete, "/api/users/123", nil)
	reqGetNodes := httptest.NewRequest(http.MethodGet, "/api/nodes", nil)

	// Admin principal can access anything
	if !requireAPITokenScope(adminPrincipal, reqGetUsers, nil) {
		t.Errorf("admin should have access to GET /api/users")
	}
	if !requireAPITokenScope(adminPrincipal, reqPostUser, nil) {
		t.Errorf("admin should have access to POST /api/users")
	}

	// * token can access anything
	if !requireAPITokenScope(allAccessAPIToken, reqGetUsers, nil) {
		t.Errorf("allAccessAPIToken should have access to GET /api/users")
	}
	if !requireAPITokenScope(allAccessAPIToken, reqPostUser, nil) {
		t.Errorf("allAccessAPIToken should have access to POST /api/users")
	}
	if !requireAPITokenScope(allAccessAPIToken, reqGetNodes, nil) {
		t.Errorf("allAccessAPIToken should have access to GET /api/nodes")
	}

	// users:read can GET /api/users and /api/users/123, but not POST or GET /api/nodes
	if !requireAPITokenScope(usersReadToken, reqGetUsers, nil) {
		t.Errorf("usersReadToken should have access to GET /api/users")
	}
	if !requireAPITokenScope(usersReadToken, reqGetUserByID, nil) {
		t.Errorf("usersReadToken should have access to GET /api/users/123")
	}
	if requireAPITokenScope(usersReadToken, reqPostUser, nil) {
		t.Errorf("usersReadToken should NOT have access to POST /api/users")
	}
	if requireAPITokenScope(usersReadToken, reqGetNodes, nil) {
		t.Errorf("usersReadToken should NOT have access to GET /api/nodes")
	}

	// users:* can GET, POST, DELETE users, but not nodes
	if !requireAPITokenScope(usersWildcardToken, reqGetUsers, nil) {
		t.Errorf("usersWildcardToken should have access to GET /api/users")
	}
	if !requireAPITokenScope(usersWildcardToken, reqPostUser, nil) {
		t.Errorf("usersWildcardToken should have access to POST /api/users")
	}
	if !requireAPITokenScope(usersWildcardToken, reqDeleteUser, nil) {
		t.Errorf("usersWildcardToken should have access to DELETE /api/users/123")
	}
	if requireAPITokenScope(usersWildcardToken, reqGetNodes, nil) {
		t.Errorf("usersWildcardToken should NOT have access to GET /api/nodes")
	}

	// users:list can GET /api/users, but not DELETE or POST
	if !requireAPITokenScope(usersListToken, reqGetUsers, nil) {
		t.Errorf("usersListToken should have access to GET /api/users")
	}
	if requireAPITokenScope(usersListToken, reqDeleteUser, nil) {
		t.Errorf("usersListToken should NOT have access to DELETE /api/users/123")
	}
}

func TestRequireAPITokenScopeWithBasePath(t *testing.T) {
	cfg := &config.BackendConfig{}
	cfg.Backend.BasePath = "/exodus_path"

	nodesReadToken := &AuthPrincipal{
		Role:      "API",
		TokenType: "jwt_api_token",
		Scopes:    []string{"nodes:read"},
	}

	reqGetNodesWithPrefix := httptest.NewRequest(http.MethodGet, "/exodus_path/api/nodes", nil)
	if !requireAPITokenScope(nodesReadToken, reqGetNodesWithPrefix, cfg) {
		t.Errorf("expected nodes:read token to access /exodus_path/api/nodes")
	}

	reqPostUserWithPrefix := httptest.NewRequest(http.MethodPost, "/exodus_path/api/users", nil)
	if requireAPITokenScope(nodesReadToken, reqPostUserWithPrefix, cfg) {
		t.Errorf("expected nodes:read token to be denied on /exodus_path/api/users")
	}
}
