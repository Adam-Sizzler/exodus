package passkeys

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/auth"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"

	"github.com/go-webauthn/webauthn/protocol"
	gowebauthn "github.com/go-webauthn/webauthn/webauthn"
)

func PasskeysHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/passkeys" && r.URL.Path != "/api/passkeys/" {
			shared.WriteJSONError(w, http.StatusNotFound, "not found")
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetPasskeys(w, r, db, cfg)
		case http.MethodPatch:
			handlePatchPasskey(w, r, db, cfg)
		case http.MethodDelete:
			handleDeletePasskey(w, r, db, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetPasskeys(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), db, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"passkeys": items,
		},
	})
}

func handlePatchPasskey(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req passkeyWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	req.Name = strings.TrimSpace(req.Name)
	if req.ID == "" {
		shared.SendError(w, http.StatusBadRequest, "id is required", nil, cfg)
		return
	}
	if len(req.Name) < 2 || len(req.Name) > 30 || !passkeyNameRegexp.MatchString(req.Name) {
		shared.SendError(w, http.StatusBadRequest, "invalid passkey name", nil, cfg)
		return
	}

	result, execErr := db.ExecContext(r.Context(), `
		UPDATE passkeys
		SET passkey_provider = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND admin_uuid = $3
	`, req.Name, req.ID, adminUUID)
	if execErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update passkey name", execErr, cfg)
		return
	}
	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to read rows affected", rowsErr, cfg)
		return
	}
	if rows == 0 {
		shared.SendError(w, http.StatusNotFound, "passkey not found", nil, cfg)
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), db, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"passkeys": items,
		},
	})
}

func handleDeletePasskey(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req passkeyWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		shared.SendError(w, http.StatusBadRequest, "id is required", nil, cfg)
		return
	}

	result, execErr := db.ExecContext(r.Context(), `DELETE FROM passkeys WHERE id = $1 AND admin_uuid = $2`, req.ID, adminUUID)
	if execErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete passkey", execErr, cfg)
		return
	}
	rows, rowsErr := result.RowsAffected()
	if rowsErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to read rows affected", rowsErr, cfg)
		return
	}
	if rows == 0 {
		shared.SendError(w, http.StatusNotFound, "passkey not found", nil, cfg)
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), db, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"passkeys": items}})
}

func RegistrationOptionsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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

		admin, err := loadWebAuthnAdmin(r.Context(), db, adminUUID)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to load admin", err, cfg)
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), db, r)
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

func VerifyRegistrationHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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

		admin, err := loadWebAuthnAdmin(r.Context(), db, adminUUID)
		if err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to load admin", err, cfg)
			return
		}

		session, ok := passkeySessions.popRegistration(admin.uuid)
		if !ok {
			sendPasskeySetupError(w, errChallengeNotFound, cfg)
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), db, r)
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

		if err := saveNewCredential(r.Context(), db, admin.uuid, credential); err != nil {
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

func AuthenticationOptionsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		resolved, err := resolvePasskeySettings(r.Context(), db, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		admin, err := loadWebAuthnAdmin(r.Context(), db, "")
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

func VerifyAuthenticationHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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

		resolved, err := resolvePasskeySettings(r.Context(), db, r)
		if err != nil {
			sendPasskeySetupError(w, err, cfg)
			return
		}

		admin, err := loadWebAuthnAdmin(r.Context(), db, "")
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

		if err := updateCredentialUsage(r.Context(), db, admin.uuid, credential); err != nil {
			sendPasskeyError(w, http.StatusInternalServerError, "failed to update passkey usage", err, cfg)
			return
		}

		sessionToken, expiresAt, err := createAdminSession(r.Context(), db, cfg, admin)
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

func currentAdminUUID(r *http.Request) (string, bool) {
	principal, ok := auth.CurrentAuthPrincipal(r.Context())
	if !ok || principal == nil {
		return "", false
	}
	adminUUID := strings.TrimSpace(principal.AdminUUID)
	if adminUUID == "" {
		return "", false
	}
	return adminUUID, true
}
