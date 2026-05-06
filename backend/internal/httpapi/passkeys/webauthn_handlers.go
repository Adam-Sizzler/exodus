package passkeys

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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
	"exodus/internal/security"

	"github.com/go-webauthn/webauthn/protocol"
	gowebauthn "github.com/go-webauthn/webauthn/webauthn"
)

const (
	passkeyRPName              = "Exodus"
	passkeySessionCookieName   = "exodus_session"
	passkeyRegistrationTimeout = 5 * time.Minute
	passkeyLoginTimeout        = time.Minute
)

var (
	errPasskeysNotEnabled    = errors.New("passkeys not enabled")
	errPasskeysNotConfigured = errors.New("passkeys not configured")
	errAdminNotFound         = errors.New("admin not found")
	errChallengeNotFound     = errors.New("challenge not found or expired")

	passkeySessions = newPasskeySessionStore()
)

type passkeySettings struct {
	Enabled bool    `json:"enabled"`
	RPID    *string `json:"rpId"`
	Origin  *string `json:"origin"`
}

type resolvedPasskeySettings struct {
	RPID    string
	Origins []string
}

type webAuthnAdmin struct {
	uuid              string
	username          string
	sessionTTLMinutes int
	credentials       []gowebauthn.Credential
}

func (a *webAuthnAdmin) WebAuthnID() []byte {
	return []byte(a.uuid)
}

func (a *webAuthnAdmin) WebAuthnName() string {
	return a.username
}

func (a *webAuthnAdmin) WebAuthnDisplayName() string {
	return "Exodus Administrator"
}

func (a *webAuthnAdmin) WebAuthnCredentials() []gowebauthn.Credential {
	return a.credentials
}

type passkeySessionStore struct {
	mu             sync.Mutex
	registration   map[string]storedPasskeySession
	authentication map[string]storedPasskeySession
}

type storedPasskeySession struct {
	session gowebauthn.SessionData
	expires time.Time
}

func newPasskeySessionStore() *passkeySessionStore {
	return &passkeySessionStore{
		registration:   make(map[string]storedPasskeySession),
		authentication: make(map[string]storedPasskeySession),
	}
}

func (s *passkeySessionStore) setRegistration(adminUUID string, session gowebauthn.SessionData) {
	s.set(s.registration, adminUUID, session, passkeyRegistrationTimeout)
}

func (s *passkeySessionStore) popRegistration(adminUUID string) (gowebauthn.SessionData, bool) {
	return s.pop(s.registration, adminUUID)
}

func (s *passkeySessionStore) setAuthentication(adminUUID string, session gowebauthn.SessionData) {
	s.set(s.authentication, adminUUID, session, passkeyLoginTimeout)
}

func (s *passkeySessionStore) popAuthentication(adminUUID string) (gowebauthn.SessionData, bool) {
	return s.pop(s.authentication, adminUUID)
}

func (s *passkeySessionStore) set(store map[string]storedPasskeySession, adminUUID string, session gowebauthn.SessionData, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expires := session.Expires
	if expires.IsZero() {
		expires = time.Now().Add(ttl)
	}
	store[adminUUID] = storedPasskeySession{
		session: session,
		expires: expires,
	}
}

func (s *passkeySessionStore) pop(store map[string]storedPasskeySession, adminUUID string) (gowebauthn.SessionData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stored, ok := store[adminUUID]
	if !ok {
		return gowebauthn.SessionData{}, false
	}
	delete(store, adminUUID)
	if time.Now().After(stored.expires) {
		return gowebauthn.SessionData{}, false
	}
	return stored.session, true
}

type verifyRegistrationRequest struct {
	Response json.RawMessage `json:"response"`
}

type verifyAuthenticationRequest struct {
	Response json.RawMessage `json:"response"`
}

func RegistrationOptionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		adminUUID, ok := currentAdminUUID(r)
		if !ok {
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		admin, err := loadWebAuthnAdmin(r.Context(), manager, adminUUID)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to load admin", err, cfg)
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), manager, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		wa, err := newWebAuthn(resolved)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to initialize passkeys", err, cfg)
			return
		}

		exclusions := make([]protocol.CredentialDescriptor, 0, len(admin.credentials))
		for i := range admin.credentials {
			exclusions = append(exclusions, admin.credentials[i].Descriptor())
		}

		creation, session, err := wa.BeginRegistration(
			admin,
			gowebauthn.WithExclusions(exclusions),
		)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to generate passkey registration options", err, cfg)
			return
		}

		passkeySessions.setRegistration(admin.uuid, *session)
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": creation.Response})
	}
}

func VerifyRegistrationHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		adminUUID, ok := currentAdminUUID(r)
		if !ok {
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req verifyRegistrationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if len(req.Response) == 0 {
			shared.WriteJSONError(w, http.StatusBadRequest, "response is required")
			return
		}

		admin, err := loadWebAuthnAdmin(r.Context(), manager, adminUUID)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to load admin", err, cfg)
			return
		}

		session, ok := passkeySessions.popRegistration(admin.uuid)
		if !ok {
			sendPasskeySetupError(w, errChallengeNotFound, cfg)
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), manager, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		wa, err := newWebAuthn(resolved)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to initialize passkeys", err, cfg)
			return
		}

		parsed, err := protocol.ParseCredentialCreationResponseBytes(req.Response)
		if err != nil {
			sendPasskeyError(w, http.StatusBadRequest, "invalid passkey registration response", err, cfg)
			return
		}

		credential, err := wa.CreateCredential(admin, session, parsed)
		if err != nil {
			sendPasskeyError(w, http.StatusForbidden, "failed to verify passkey registration", err, cfg)
			return
		}

		if err := saveNewCredential(r.Context(), manager, admin.uuid, credential); err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to save passkey", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"verified": true,
			},
		})
	}
}

func AuthenticationOptionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), manager, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		admin, err := loadFirstWebAuthnAdmin(r.Context(), manager)
		if err != nil {
			sendPasskeyError(w, http.StatusForbidden, "passkey authentication is not available", err, cfg)
			return
		}
		if len(admin.credentials) == 0 {
			sendPasskeySetupError(w, errPasskeysNotConfigured, cfg)
			return
		}

		wa, err := newWebAuthn(resolved)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to initialize passkeys", err, cfg)
			return
		}

		assertion, session, err := wa.BeginLogin(admin)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to generate passkey authentication options", err, cfg)
			return
		}

		passkeySessions.setAuthentication(admin.uuid, *session)
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": assertion.Response})
	}
}

func VerifyAuthenticationHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req verifyAuthenticationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if len(req.Response) == 0 {
			shared.WriteJSONError(w, http.StatusBadRequest, "response is required")
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), manager, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		admin, err := loadFirstWebAuthnAdmin(r.Context(), manager)
		if err != nil {
			sendPasskeyError(w, http.StatusForbidden, "passkey authentication is not available", err, cfg)
			return
		}

		session, ok := passkeySessions.popAuthentication(admin.uuid)
		if !ok {
			sendPasskeySetupError(w, errChallengeNotFound, cfg)
			return
		}

		wa, err := newWebAuthn(resolved)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to initialize passkeys", err, cfg)
			return
		}

		parsed, err := protocol.ParseCredentialRequestResponseBytes(req.Response)
		if err != nil {
			sendPasskeyError(w, http.StatusBadRequest, "invalid passkey authentication response", err, cfg)
			return
		}

		credential, err := wa.ValidateLogin(admin, session, parsed)
		if err != nil {
			sendPasskeyError(w, http.StatusForbidden, "failed to verify passkey authentication", err, cfg)
			return
		}

		if err := updateCredentialUsage(r.Context(), manager, admin.uuid, credential); err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to update passkey usage", err, cfg)
			return
		}

		sessionToken, expiresAt, err := createAdminSession(r.Context(), manager, admin)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to create session", err, cfg)
			return
		}

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     passkeySessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Unix(expiresAt, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
		})

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"accessToken": sessionToken,
			},
		})
	}
}

func newWebAuthn(settings resolvedPasskeySettings) (*gowebauthn.WebAuthn, error) {
	return gowebauthn.New(&gowebauthn.Config{
		RPID:          settings.RPID,
		RPDisplayName: passkeyRPName,
		RPOrigins:     settings.Origins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementPreferred,
			UserVerification: protocol.VerificationRequired,
		},
		AttestationPreference: protocol.PreferNoAttestation,
		Timeouts: gowebauthn.TimeoutsConfig{
			Login: gowebauthn.TimeoutConfig{
				Enforce: true,
				Timeout: passkeyLoginTimeout,
			},
			Registration: gowebauthn.TimeoutConfig{
				Enforce: true,
				Timeout: passkeyRegistrationTimeout,
			},
		},
	})
}

func resolvePasskeySettings(ctx context.Context, manager *dbmanager.DatabaseManager, r *http.Request) (resolvedPasskeySettings, error) {
	settings, err := loadPasskeySettings(ctx, manager)
	if err != nil {
		return resolvedPasskeySettings{}, err
	}
	if !settings.Enabled {
		return resolvedPasskeySettings{}, errPasskeysNotEnabled
	}

	requestHost := firstHeaderValue(r.Header.Get("X-Forwarded-Host"))
	if requestHost == "" {
		requestHost = r.Host
	}
	host, port := splitHostPortLoose(requestHost)
	localRequest := isLocalHost(host)

	currentOrigin := currentRequestOrigin(r)
	origins := make([]string, 0, 4)

	if localRequest {
		rpID := normalizeRPID(host)
		if rpID == "" {
			return resolvedPasskeySettings{}, errPasskeysNotConfigured
		}
		origins = appendOrigin(origins, currentOrigin)
		origins = appendLocalOrigins(origins, rpID, port)
		if settings.Origin != nil {
			origins = appendOrigin(origins, normalizeOrigin(*settings.Origin))
		}
		if len(origins) == 0 {
			return resolvedPasskeySettings{}, errPasskeysNotConfigured
		}
		return resolvedPasskeySettings{RPID: rpID, Origins: origins}, nil
	}

	rpID := ""
	if settings.RPID != nil {
		rpID = normalizeRPID(*settings.RPID)
	}
	if rpID == "" || settings.Origin == nil {
		return resolvedPasskeySettings{}, errPasskeysNotConfigured
	}

	origin := normalizeOrigin(*settings.Origin)
	if origin == "" {
		return resolvedPasskeySettings{}, errPasskeysNotConfigured
	}
	origins = appendOrigin(origins, origin)
	if sameOriginHost(origin, currentOrigin) {
		origins = appendOrigin(origins, currentOrigin)
	}
	return resolvedPasskeySettings{RPID: rpID, Origins: origins}, nil
}

func loadPasskeySettings(ctx context.Context, manager *dbmanager.DatabaseManager) (passkeySettings, error) {
	settings := passkeySettings{}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT passkey_settings
			FROM exodus_settings
			WHERE id = 1
			LIMIT 1
		`)

		var raw sql.NullString
		if err := row.Scan(&raw); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		if raw.Valid && strings.TrimSpace(raw.String) != "" {
			if err := json.Unmarshal([]byte(raw.String), &settings); err != nil {
				return err
			}
		}
		return nil
	})
	return settings, err
}

func loadFirstWebAuthnAdmin(ctx context.Context, manager *dbmanager.DatabaseManager) (*webAuthnAdmin, error) {
	return loadWebAuthnAdmin(ctx, manager, "")
}

func loadWebAuthnAdmin(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string) (*webAuthnAdmin, error) {
	admin := &webAuthnAdmin{sessionTTLMinutes: 60}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var row *sql.Row
		if strings.TrimSpace(adminUUID) == "" {
			row = db.QueryRowContext(ctx, `
				SELECT uuid, username, COALESCE(session_ttl_minutes, 60)
				FROM admin
				ORDER BY created_at ASC
				LIMIT 1
			`)
		} else {
			row = db.QueryRowContext(ctx, `
				SELECT uuid, username, COALESCE(session_ttl_minutes, 60)
				FROM admin
				WHERE uuid = ?
				LIMIT 1
			`, adminUUID)
		}

		if err := row.Scan(&admin.uuid, &admin.username, &admin.sessionTTLMinutes); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errAdminNotFound
			}
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT id, public_key, counter, COALESCE(transports, ''), COALESCE(device_type, ''), backed_up
			FROM passkeys
			WHERE admin_uuid = ?
			ORDER BY created_at ASC
		`, admin.uuid)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				id        string
				publicKey []byte
				counter   int64
				transRaw  string
				device    string
				backedUp  bool
			)
			if err := rows.Scan(&id, &publicKey, &counter, &transRaw, &device, &backedUp); err != nil {
				return err
			}

			rawID, err := decodeCredentialID(id)
			if err != nil {
				return fmt.Errorf("decode credential id: %w", err)
			}
			backupEligible := strings.EqualFold(device, "multiDevice") || backedUp
			admin.credentials = append(admin.credentials, gowebauthn.Credential{
				ID:        rawID,
				PublicKey: publicKey,
				Transport: parseTransports(transRaw),
				Flags: gowebauthn.CredentialFlags{
					BackupEligible: backupEligible,
					BackupState:    backedUp,
				},
				Authenticator: gowebauthn.Authenticator{
					SignCount: uint32Counter(counter),
				},
			})
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(admin.uuid) == "" {
		return nil, errAdminNotFound
	}
	return admin, nil
}

func saveNewCredential(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string, credential *gowebauthn.Credential) error {
	credentialID := encodeCredentialID(credential.ID)
	transports := transportsToCSV(credential.Transport)
	deviceType := "singleDevice"
	if credential.Flags.BackupEligible {
		deviceType = "multiDevice"
	}
	provider := passkeyProviderName(credential.Transport)

	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO passkeys (
				id, admin_uuid, public_key, counter, device_type, backed_up, transports, passkey_provider
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, credentialID, adminUUID, credential.PublicKey, int64(credential.Authenticator.SignCount), deviceType, credential.Flags.BackupState, transports, provider)
		return err
	})
}

func updateCredentialUsage(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string, credential *gowebauthn.Credential) error {
	credentialID := encodeCredentialID(credential.ID)
	deviceType := "singleDevice"
	if credential.Flags.BackupEligible {
		deviceType = "multiDevice"
	}

	var rowsAffected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `
			UPDATE passkeys
			SET counter = ?, device_type = ?, backed_up = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND admin_uuid = ?
		`, int64(credential.Authenticator.SignCount), deviceType, credential.Flags.BackupState, credentialID, adminUUID)
		if err != nil {
			return err
		}
		rowsAffected, err = result.RowsAffected()
		return err
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errAdminNotFound
	}
	return nil
}

func createAdminSession(ctx context.Context, manager *dbmanager.DatabaseManager, admin *webAuthnAdmin) (string, int64, error) {
	sessionToken, err := security.GenerateRandomToken(48)
	if err != nil {
		return "", 0, err
	}

	ttlMinutes := admin.sessionTTLMinutes
	if ttlMinutes <= 0 {
		ttlMinutes = 60
	}
	expiresAt := time.Now().UTC().Add(time.Duration(ttlMinutes) * time.Minute).Unix()

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if _, execErr := db.ExecContext(ctx, "DELETE FROM admin_sessions WHERE admin_uuid = ?", admin.uuid); execErr != nil {
			return execErr
		}
		_, execErr := db.ExecContext(ctx, `
			INSERT INTO admin_sessions (session_token, admin_uuid, expires_at)
			VALUES (?, ?, ?)
		`, sessionToken, admin.uuid, expiresAt)
		return execErr
	})
	if err != nil {
		return "", 0, err
	}
	return sessionToken, expiresAt, nil
}

func sendPasskeySetupError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errPasskeysNotEnabled):
		shared.SendError(w, http.StatusForbidden, "passkeys not enabled", err, cfg)
	case errors.Is(err, errPasskeysNotConfigured):
		shared.SendError(w, http.StatusForbidden, "passkeys not configured", err, cfg)
	case errors.Is(err, errChallengeNotFound):
		shared.SendError(w, http.StatusForbidden, "challenge not found or expired", err, cfg)
	default:
		shared.SendError(w, http.StatusInternalServerError, "passkey setup error", err, cfg)
	}
}

func sendPasskeyError(w http.ResponseWriter, status int, msg string, err error, cfg *config.BackendConfig) {
	shared.SendError(w, status, msg, err, cfg)
}

func encodeCredentialID(raw []byte) string {
	return base64.RawURLEncoding.EncodeToString(raw)
}

func decodeCredentialID(id string) ([]byte, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("empty credential id")
	}
	raw, err := base64.RawURLEncoding.DecodeString(id)
	if err == nil && len(raw) > 0 {
		return raw, nil
	}
	return []byte(id), nil
}

func parseTransports(raw string) []protocol.AuthenticatorTransport {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	transports := make([]protocol.AuthenticatorTransport, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		transports = append(transports, protocol.AuthenticatorTransport(part))
	}
	return transports
}

func transportsToCSV(transports []protocol.AuthenticatorTransport) string {
	if len(transports) == 0 {
		return ""
	}
	parts := make([]string, 0, len(transports))
	for _, transport := range transports {
		value := strings.TrimSpace(string(transport))
		if value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, ",")
}

func passkeyProviderName(transports []protocol.AuthenticatorTransport) string {
	for _, transport := range transports {
		switch transport {
		case protocol.Internal:
			return "Platform authenticator"
		case protocol.USB:
			return "USB security key"
		case protocol.NFC:
			return "NFC security key"
		case protocol.BLE:
			return "Bluetooth security key"
		case protocol.Hybrid:
			return "Hybrid passkey"
		}
	}
	return "Passkey"
}

func uint32Counter(counter int64) uint32 {
	if counter <= 0 {
		return 0
	}
	max := int64(^uint32(0))
	if counter > max {
		return ^uint32(0)
	}
	return uint32(counter)
}

func currentRequestOrigin(r *http.Request) string {
	if origin := normalizeOrigin(r.Header.Get("Origin")); origin != "" {
		return origin
	}

	scheme := firstHeaderValue(r.Header.Get("X-Forwarded-Proto"))
	if scheme == "" {
		scheme = firstHeaderValue(r.Header.Get("X-Forwarded-Scheme"))
	}
	if scheme == "" {
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := firstHeaderValue(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = r.Host
	}
	return normalizeOrigin(scheme + "://" + host)
}

func normalizeOrigin(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ""
	}
	return scheme + "://" + strings.ToLower(parsed.Host)
}

func normalizeRPID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if parsed, err := url.Parse(raw); err == nil && parsed.Hostname() != "" {
		return strings.ToLower(parsed.Hostname())
	}
	if idx := strings.IndexByte(raw, '/'); idx >= 0 {
		raw = raw[:idx]
	}
	host, _ := splitHostPortLoose(raw)
	return strings.ToLower(strings.Trim(host, "[]"))
}

func appendOrigin(origins []string, origin string) []string {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return origins
	}
	for _, existing := range origins {
		if strings.EqualFold(existing, origin) {
			return origins
		}
	}
	return append(origins, origin)
}

func appendLocalOrigins(origins []string, host, port string) []string {
	host = strings.ToLower(strings.Trim(host, "[]"))
	if host == "" {
		return origins
	}
	for _, scheme := range []string{"https", "http"} {
		origins = appendOrigin(origins, scheme+"://"+host)
		if port != "" {
			origins = appendOrigin(origins, scheme+"://"+net.JoinHostPort(host, port))
		}
	}
	return origins
}

func sameOriginHost(a, b string) bool {
	if a == "" || b == "" {
		return false
	}
	originA, errA := url.Parse(a)
	originB, errB := url.Parse(b)
	if errA != nil || errB != nil {
		return false
	}
	return strings.EqualFold(originA.Hostname(), originB.Hostname())
}

func isLocalHost(host string) bool {
	host = strings.ToLower(strings.Trim(host, "[]"))
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func splitHostPortLoose(value string) (host string, port string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ""
	}

	if h, p, err := net.SplitHostPort(value); err == nil {
		return strings.Trim(h, "[]"), p
	}

	if strings.HasPrefix(value, "[") {
		if end := strings.Index(value, "]"); end > 0 {
			host = value[1:end]
			rest := strings.TrimPrefix(value[end+1:], ":")
			if _, err := strconv.Atoi(rest); err == nil {
				return host, rest
			}
			return host, ""
		}
	}

	if idx := strings.LastIndex(value, ":"); idx > 0 && strings.Count(value, ":") == 1 {
		candidatePort := value[idx+1:]
		if _, err := strconv.Atoi(candidatePort); err == nil {
			return value[:idx], candidatePort
		}
	}

	return strings.Trim(value, "[]"), ""
}

func firstHeaderValue(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	return strings.TrimSpace(parts[0])
}
