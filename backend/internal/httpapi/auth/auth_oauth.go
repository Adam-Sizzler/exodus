package auth

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
	"exodus/internal/notifications"
	"exodus/internal/security"
)

const (
	oauthStateTTL        = 10 * time.Minute
	oauthScope           = "openid email profile"
	customOAuthClaimName = "exodusAccess"
)

var oauthStateCache = struct {
	sync.Mutex
	items map[string]oauthStateEntry
}{items: map[string]oauthStateEntry{}}

type oauthStateEntry struct {
	State        string
	CodeVerifier string
	ExpiresAt    time.Time
}

type oauthAuthorizeRequest struct {
	Provider string `json:"provider"`
}

type oauthCallbackRequest struct {
	Provider string `json:"provider"`
	Code     string `json:"code"`
	State    string `json:"state"`
}

type oauthSettings struct {
	Github   oauthProviderSettings `json:"github"`
	PocketID oauthProviderSettings `json:"pocketid"`
	Yandex   oauthProviderSettings `json:"yandex"`
	Keycloak oauthProviderSettings `json:"keycloak"`
	Generic  oauthProviderSettings `json:"generic"`
	Telegram oauthProviderSettings `json:"telegram"`
}

type oauthProviderSettings struct {
	Enabled          bool     `json:"enabled"`
	ClientID         string   `json:"clientId"`
	ClientSecret     string   `json:"clientSecret"`
	PlainDomain      string   `json:"plainDomain"`
	Realm            string   `json:"realm"`
	FrontendDomain   string   `json:"frontendDomain"`
	KeycloakDomain   string   `json:"keycloakDomain"`
	WithPKCE         bool     `json:"withPkce"`
	AuthorizationURL string   `json:"authorizationUrl"`
	TokenURL         string   `json:"tokenUrl"`
	AllowedEmails    []string `json:"allowedEmails"`
	AllowedIDs       []string `json:"allowedIds"`
}

type oauthTokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func OAuth2AuthorizeHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req oauthAuthorizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		provider := strings.ToLower(strings.TrimSpace(req.Provider))
		if !isSupportedOAuthProvider(provider) {
			shared.SendError(w, http.StatusBadRequest, "OAuth2 provider not found", nil, cfg)
			return
		}
		if !isLoginAllowed(manager) {
			shared.SendError(w, http.StatusForbidden, "login is not allowed", nil, cfg)
			return
		}
		settings, err := loadOAuthSettings(manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to load OAuth2 settings", err, cfg)
			return
		}
		providerSettings := getOAuthProviderSettings(settings, provider)
		if !providerSettings.Enabled {
			shared.SendError(w, http.StatusForbidden, "OAuth2 provider is disabled", nil, cfg)
			return
		}

		state, err := security.GenerateRandomToken(32)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to generate OAuth2 state", err, cfg)
			return
		}
		codeVerifier := ""
		authURL, err := buildAuthorizationURL(provider, providerSettings, state, &codeVerifier)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "OAuth2 authorize error", err, cfg)
			return
		}
		storeOAuthState(provider, state, codeVerifier)
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"authorizationUrl": authURL},
		})
	}
}

func OAuth2CallbackHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req oauthCallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		provider := strings.ToLower(strings.TrimSpace(req.Provider))
		if !isSupportedOAuthProvider(provider) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "unsupported_provider", r)
			shared.SendError(w, http.StatusBadRequest, "OAuth2 provider not found", nil, cfg)
			return
		}
		stateEntry, ok := takeOAuthState(provider)
		if !ok || stateEntry.State != strings.TrimSpace(req.State) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "state_mismatch", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 state mismatch", nil, cfg)
			return
		}
		if !isLoginAllowed(manager) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "login_not_allowed", r)
			shared.SendError(w, http.StatusForbidden, "login is not allowed", nil, cfg)
			return
		}
		settings, err := loadOAuthSettings(manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to load OAuth2 settings", err, cfg)
			return
		}
		providerSettings := getOAuthProviderSettings(settings, provider)
		if !providerSettings.Enabled {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "provider_disabled", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 provider is disabled", nil, cfg)
			return
		}
		email, hasCustomClaim, err := exchangeOAuthCode(r.Context(), provider, providerSettings, strings.TrimSpace(req.Code), stateEntry.CodeVerifier)
		if err != nil {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "callback_error", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 callback error", err, cfg)
			return
		}
		if email == "" || !isOAuthPrincipalAllowed(provider, email, hasCustomClaim, providerSettings) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, email, "", "email_not_allowed", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 principal is not allowed", nil, cfg)
			return
		}
		token, adminUUID, err := createFirstAdminSession(w, r, manager, cfg)
		if err != nil {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, email, "", "session_create_failed", r)
			shared.SendError(w, http.StatusForbidden, "failed to create OAuth2 session", err, cfg)
			return
		}
		emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptSuccess, "oauth2", provider, email, adminUUID, "", r)
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"accessToken": token},
		})
	}
}

func isLoginAllowed(manager *dbmanager.DatabaseManager) bool {
	_, _, _, hasAdmin, err := getBootstrapData(manager)
	return err == nil && hasAdmin
}

func createFirstAdminSession(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (string, string, error) {
	var adminUUID, username, role string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT uuid, username, role
			FROM admin
			WHERE UPPER(role) = 'ADMIN'
			ORDER BY created_at ASC
			LIMIT 1
		`)
		return row.Scan(&adminUUID, &username, &role)
	})
	if err != nil {
		return "", "", err
	}

	accessToken, expiresAt, err := createAdminAccessToken(cfg, username, adminUUID, role)
	if err != nil {
		return "", "", err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Unix(expiresAt, 0).UTC(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   middleware.IsSecureRequest(r, cfg),
	})
	return accessToken, adminUUID, nil
}

func emitExternalLoginNotification(ctx context.Context, cfg *config.BackendConfig, event, method, provider, identifier, adminUUID, reason string, r *http.Request) {
	data := map[string]any{
		"method":   method,
		"provider": provider,
	}
	if identifier != "" {
		data["identifier"] = identifier
		data["username"] = identifier
	}
	if adminUUID != "" {
		data["adminUuid"] = adminUUID
	}
	if reason != "" {
		data["reason"] = reason
	}
	if r != nil {
		data["remoteAddr"] = r.RemoteAddr
		data["userAgent"] = r.UserAgent()
		data["path"] = r.URL.Path
	}
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeService,
		Event: event,
		Data:  data,
	})
}

func loadOAuthSettings(manager *dbmanager.DatabaseManager) (oauthSettings, error) {
	var out oauthSettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var raw sql.NullString
		if err := db.QueryRow(`SELECT oauth2_settings FROM exodus_settings WHERE id = 1 LIMIT 1`).Scan(&raw); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		if raw.Valid && strings.TrimSpace(raw.String) != "" {
			return json.Unmarshal([]byte(raw.String), &out)
		}
		return nil
	})
	return out, err
}

func isSupportedOAuthProvider(provider string) bool {
	switch provider {
	case "github", "pocketid", "yandex", "keycloak", "generic", "telegram":
		return true
	default:
		return false
	}
}

func getOAuthProviderSettings(settings oauthSettings, provider string) oauthProviderSettings {
	switch provider {
	case "github":
		return settings.Github
	case "pocketid":
		return settings.PocketID
	case "yandex":
		return settings.Yandex
	case "keycloak":
		return settings.Keycloak
	case "generic":
		return settings.Generic
	case "telegram":
		return settings.Telegram
	default:
		return oauthProviderSettings{}
	}
}

func buildAuthorizationURL(provider string, settings oauthProviderSettings, state string, codeVerifier *string) (string, error) {
	switch provider {
	case "github":
		if settings.ClientID == "" || settings.ClientSecret == "" {
			return "", errors.New("github OAuth2 settings are incomplete")
		}
		return makeOAuthURL("https://github.com/login/oauth/authorize", url.Values{
			"client_id": {settings.ClientID},
			"state":     {state},
			"scope":     {"user:email"},
		})
	case "yandex":
		if settings.ClientID == "" || settings.ClientSecret == "" {
			return "", errors.New("yandex OAuth2 settings are incomplete")
		}
		return makeOAuthURL("https://oauth.yandex.ru/authorize", url.Values{
			"response_type": {"code"},
			"client_id":     {settings.ClientID},
			"state":         {state},
			"scope":         {"login:email"},
		})
	case "pocketid":
		if settings.ClientID == "" || settings.ClientSecret == "" || settings.PlainDomain == "" {
			return "", errors.New("pocketid OAuth2 settings are incomplete")
		}
		return makeOAuthURL("https://"+settings.PlainDomain+"/authorize", url.Values{
			"response_type": {"code"},
			"client_id":     {settings.ClientID},
			"state":         {state},
			"scope":         {oauthScope},
		})
	case "keycloak":
		if settings.ClientID == "" || settings.ClientSecret == "" || settings.KeycloakDomain == "" || settings.Realm == "" || settings.FrontendDomain == "" {
			return "", errors.New("keycloak OAuth2 settings are incomplete")
		}
		verifier, challenge, err := generatePKCEPair()
		if err != nil {
			return "", err
		}
		*codeVerifier = verifier
		return makeOAuthURL(fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/auth", settings.KeycloakDomain, settings.Realm), url.Values{
			"response_type":         {"code"},
			"client_id":             {settings.ClientID},
			"redirect_uri":          {oauthRedirectURI(settings.FrontendDomain, provider)},
			"state":                 {state},
			"scope":                 {oauthScope},
			"code_challenge_method": {"S256"},
			"code_challenge":        {challenge},
		})
	case "generic":
		if settings.ClientID == "" || settings.ClientSecret == "" || settings.AuthorizationURL == "" || settings.TokenURL == "" || settings.FrontendDomain == "" {
			return "", errors.New("generic OAuth2 settings are incomplete")
		}
		values := url.Values{
			"response_type": {"code"},
			"client_id":     {settings.ClientID},
			"redirect_uri":  {oauthRedirectURI(settings.FrontendDomain, provider)},
			"state":         {state},
			"scope":         {oauthScope},
		}
		if settings.WithPKCE {
			verifier, challenge, err := generatePKCEPair()
			if err != nil {
				return "", err
			}
			*codeVerifier = verifier
			values.Set("code_challenge_method", "S256")
			values.Set("code_challenge", challenge)
		}
		return makeOAuthURL(settings.AuthorizationURL, values)
	case "telegram":
		if settings.ClientID == "" || settings.ClientSecret == "" || settings.FrontendDomain == "" {
			return "", errors.New("telegram OAuth2 settings are incomplete")
		}
		verifier, challenge, err := generatePKCEPair()
		if err != nil {
			return "", err
		}
		*codeVerifier = verifier
		return makeOAuthURL("https://oauth.telegram.org/auth", url.Values{
			"response_type":         {"code"},
			"client_id":             {settings.ClientID},
			"redirect_uri":          {oauthRedirectURI(settings.FrontendDomain, provider)},
			"state":                 {state},
			"scope":                 {"openid profile telegram:bot_access"},
			"code_challenge_method": {"S256"},
			"code_challenge":        {challenge},
		})
	default:
		return "", errors.New("unsupported OAuth2 provider")
	}
}

func makeOAuthURL(rawURL string, values url.Values) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, vals := range values {
		if len(vals) > 0 && strings.TrimSpace(vals[0]) != "" {
			q.Set(key, vals[0])
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func exchangeOAuthCode(ctx context.Context, provider string, settings oauthProviderSettings, code, codeVerifier string) (string, bool, error) {
	switch provider {
	case "github":
		token, err := requestOAuthToken(ctx, "https://github.com/login/oauth/access_token", url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
		})
		if err != nil {
			return "", false, err
		}
		return fetchGithubPrimaryEmail(ctx, token.AccessToken)
	case "yandex":
		token, err := requestOAuthToken(ctx, "https://oauth.yandex.ru/token", url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
		})
		if err != nil {
			return "", false, err
		}
		return fetchYandexEmail(ctx, token.AccessToken)
	case "pocketid":
		token, err := requestOAuthToken(ctx, "https://"+settings.PlainDomain+"/api/oidc/token", url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
		})
		if err != nil {
			return "", false, err
		}
		return extractEmailFromIDToken(token.IDToken)
	case "keycloak":
		token, err := requestOAuthToken(ctx, fmt.Sprintf("https://%s/realms/%s/protocol/openid-connect/token", settings.KeycloakDomain, settings.Realm), url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
			"redirect_uri":  {oauthRedirectURI(settings.FrontendDomain, provider)},
			"code_verifier": {codeVerifier},
		})
		if err != nil {
			return "", false, err
		}
		return extractEmailFromIDToken(token.IDToken)
	case "generic":
		values := url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
			"redirect_uri":  {oauthRedirectURI(settings.FrontendDomain, provider)},
		}
		if settings.WithPKCE {
			values.Set("code_verifier", codeVerifier)
		}
		token, err := requestOAuthToken(ctx, settings.TokenURL, values)
		if err != nil {
			return "", false, err
		}
		return extractEmailFromIDToken(token.IDToken)
	case "telegram":
		token, err := requestOAuthToken(ctx, "https://oauth.telegram.org/token", url.Values{
			"grant_type":    {"authorization_code"},
			"client_id":     {settings.ClientID},
			"client_secret": {settings.ClientSecret},
			"code":          {code},
			"redirect_uri":  {oauthRedirectURI(settings.FrontendDomain, provider)},
			"code_verifier": {codeVerifier},
		})
		if err != nil {
			return "", false, err
		}
		return extractTelegramIDFromIDToken(token.IDToken)
	default:
		return "", false, errors.New("unsupported OAuth2 provider")
	}
}

func requestOAuthToken(ctx context.Context, endpoint string, values url.Values) (oauthTokenResponse, error) {
	var out oauthTokenResponse
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	if err != nil {
		return out, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return out, fmt.Errorf("token endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return out, err
	}
	if out.Error != "" {
		return out, fmt.Errorf("%s: %s", out.Error, out.Description)
	}
	return out, nil
}

func fetchGithubPrimaryEmail(ctx context.Context, accessToken string) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "Exodus")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return "", false, fmt.Errorf("github email endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var emails []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", false, err
	}
	for _, item := range emails {
		if item.Primary {
			return item.Email, false, nil
		}
	}
	return "", false, errors.New("github primary email not found")
}

func fetchYandexEmail(ctx context.Context, accessToken string) (string, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://login.yandex.ru/info?format=json", nil)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("User-Agent", "Exodus")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		return "", false, fmt.Errorf("yandex info endpoint returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var payload struct {
		DefaultEmail string `json:"default_email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", false, err
	}
	if payload.DefaultEmail == "" {
		return "", false, errors.New("yandex email not found")
	}
	return payload.DefaultEmail, false, nil
}

func extractEmailFromIDToken(idToken string) (string, bool, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return "", false, errors.New("invalid id_token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false, err
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", false, err
	}
	email, _ := claims["email"].(string)
	customClaim, _ := claims[customOAuthClaimName].(bool)
	if email == "" {
		return "", customClaim, errors.New("missing email in id_token")
	}
	return email, customClaim, nil
}

func extractTelegramIDFromIDToken(idToken string) (string, bool, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return "", false, errors.New("invalid id_token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false, err
	}
	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", false, err
	}
	switch id := claims["id"].(type) {
	case string:
		id = strings.TrimSpace(id)
		if id == "" {
			return "", false, errors.New("missing telegram id in id_token")
		}
		return id, false, nil
	case float64:
		return strconv.FormatInt(int64(id), 10), false, nil
	default:
		return "", false, errors.New("missing telegram id in id_token")
	}
}

func storeOAuthState(provider, state, codeVerifier string) {
	oauthStateCache.Lock()
	defer oauthStateCache.Unlock()
	now := time.Now()
	for key, item := range oauthStateCache.items {
		if now.After(item.ExpiresAt) {
			delete(oauthStateCache.items, key)
		}
	}
	oauthStateCache.items[provider] = oauthStateEntry{
		State:        state,
		CodeVerifier: codeVerifier,
		ExpiresAt:    now.Add(oauthStateTTL),
	}
}

func takeOAuthState(provider string) (oauthStateEntry, bool) {
	oauthStateCache.Lock()
	defer oauthStateCache.Unlock()
	item, ok := oauthStateCache.items[provider]
	delete(oauthStateCache.items, provider)
	if !ok || time.Now().After(item.ExpiresAt) {
		return oauthStateEntry{}, false
	}
	return item, true
}

func generatePKCEPair() (string, string, error) {
	verifier, err := security.GenerateRandomToken(64)
	if err != nil {
		return "", "", err
	}
	sum := sha256.Sum256([]byte(verifier))
	return verifier, base64.RawURLEncoding.EncodeToString(sum[:]), nil
}

func oauthRedirectURI(frontendDomain, provider string) string {
	return "https://" + strings.TrimRight(frontendDomain, "/") + "/oauth2/callback/" + provider
}

func isEmailAllowed(email string, allowed []string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, item := range allowed {
		if strings.ToLower(strings.TrimSpace(item)) == email {
			return true
		}
	}
	return false
}

func isOAuthPrincipalAllowed(provider, principal string, hasCustomClaim bool, settings oauthProviderSettings) bool {
	if provider == "telegram" {
		return containsString(settings.AllowedIDs, principal)
	}
	return hasCustomClaim || isEmailAllowed(principal, settings.AllowedEmails)
}

func containsString(items []string, target string) bool {
	target = strings.TrimSpace(target)
	for _, item := range items {
		if strings.TrimSpace(item) == target {
			return true
		}
	}
	return false
}
