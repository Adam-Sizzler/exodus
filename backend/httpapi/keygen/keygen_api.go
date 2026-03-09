package keygen

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"v2ray-stat/backend/config"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/shared"
)

type secretPayload struct {
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
}

func KeygenHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var pubKey string
		var caCert sql.NullString
		var clientCert sql.NullString
		var clientKey sql.NullString
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return db.QueryRowContext(r.Context(), `
				SELECT pub_key, ca_cert, client_cert, client_key
				FROM keygen
				ORDER BY created_at ASC
				LIMIT 1
			`).Scan(&pubKey, &caCert, &clientCert, &clientKey)
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch keygen data", err, cfg)
			return
		}

		// Remnawave returns a base64 payload for node SECRET_KEY.
		// If cert material exists, we keep the same payload shape.
		payload := pubKey
		if caCert.Valid && clientCert.Valid && clientKey.Valid && caCert.String != "" && clientCert.String != "" && clientKey.String != "" {
			raw, marshalErr := json.Marshal(secretPayload{
				NodeCertPem:  clientCert.String,
				NodeKeyPem:   clientKey.String,
				CaCertPem:    caCert.String,
				JWTPublicKey: pubKey,
			})
			if marshalErr == nil {
				payload = base64.StdEncoding.EncodeToString(raw)
			}
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"pubKey": payload,
			},
		})
	}
}
