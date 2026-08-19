package keygen

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/security"
)

type secretPayload struct {
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
}

type KeygenResponse struct {
	Response KeygenPayload `json:"response"`
}

type KeygenPayload struct {
	SecretKey string `json:"secretKey"`
	GrpcToken string `json:"grpcToken"`
}

// KeygenHandler godoc
// @Summary      Generate node keys and certs
// @Description  Generate node TLS certificates, public keys, and gRPC auth tokens for node provisioning
// @Tags         Keygen Controller
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  KeygenResponse
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /keygen [get]
func KeygenHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var (
			pubKey string
			caCert string
			caKey  string
		)
		err := db.QueryRowContext(r.Context(), `
			SELECT pub_key, ca_cert, ca_key
			FROM keygen
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&pubKey, &caCert, &caKey)
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
		grpcToken, err := security.GenerateGRPCAuthToken()
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to generate grpc auth token", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"secretKey": payload,
				"grpcToken": grpcToken,
			},
		})
	}
}
