package keygen

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	"exodus/internal/security"
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

		var (
			pubKey    string
			caCert    string
			caKey     string
			grpcToken string
		)
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return db.QueryRowContext(r.Context(), `
				SELECT pub_key, ca_cert, ca_key, grpc_auth_token
				FROM keygen
				ORDER BY created_at ASC
				LIMIT 1
			`).Scan(&pubKey, &caCert, &caKey, &grpcToken)
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch keygen data", err, cfg)
			return
		}

		nodeCert, err := security.GenerateNodeCert(caCert, caKey)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to generate node certificate", err, cfg)
			return
		}

		raw, err := json.Marshal(secretPayload{
			NodeCertPem:  nodeCert.NodeCertPEM,
			NodeKeyPem:   nodeCert.NodeKeyPEM,
			CaCertPem:    caCert,
			JWTPublicKey: pubKey,
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to encode secret payload", err, cfg)
			return
		}
		payload := base64.StdEncoding.EncodeToString(raw)

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"pubKey":    payload,
				"grpcToken": grpcToken,
			},
		})
	}
}
