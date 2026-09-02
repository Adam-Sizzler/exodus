package system

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	subscriptionresponserules "exodus/internal/httpapi/subscriptionresponserules"

	"golang.org/x/crypto/curve25519"
)

type testSRRMatcherRequest struct {
	ResponseRules subscriptionresponserules.Config `json:"responseRules"`
}

// GenerateX25519Handler godoc
// @Summary      Generate x25519 keypairs
// @Description  Generate a batch of 30 X25519 keypairs for VLESS/Reality or WireGuard configs
// @Tags         System Controller
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Router       /system/tools/x25519/generate [get]
func GenerateX25519Handler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		keypairs := make([]map[string]string, 0, 30)
		for i := 0; i < 30; i++ {
			privateKey := make([]byte, curve25519.ScalarSize)
			if _, err := rand.Read(privateKey); err != nil {
				shared.SendAPIError(w, shared.ErrGenerateX25519KeypairFailed.WithCause(err), cfg)
				return
			}
			publicKey, err := curve25519.X25519(privateKey, curve25519.Basepoint)
			if err != nil {
				shared.SendAPIError(w, shared.ErrGenerateX25519KeypairFailed.WithCause(err), cfg)
				return
			}
			keypairs = append(keypairs, map[string]string{
				"publicKey":  base64.RawURLEncoding.EncodeToString(publicKey),
				"privateKey": base64.RawURLEncoding.EncodeToString(privateKey),
			})
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"keypairs": keypairs},
		})
	}
}

// TestSRRMatcherHandler godoc
// @Summary      Test subscription response rules matcher
// @Description  Evaluate SRR rules against request headers
// @Tags         System Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      testSRRMatcherRequest  true  "Response rules"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  shared.ErrorResponse
// @Router       /system/testers/srr-matcher [post]
func TestSRRMatcherHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req testSRRMatcherRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if err := validateResponseRules(req.ResponseRules); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid response rules", err, cfg)
			return
		}
		result := subscriptionresponserules.MatchRulesDetailed(&req.ResponseRules, r.Header, "", nil, "")
		responseType := result.ResponseType
		if responseType == "" {
			responseType = "BLOCK"
		}
		outputHeaders := map[string]string{}
		if result.MatchedRule != nil && result.MatchedRule.ResponseModifications != nil {
			for _, header := range result.MatchedRule.ResponseModifications.Headers {
				outputHeaders[header.Key] = header.Value
			}
		}
		inputHeaders := map[string]string{}
		for key, values := range r.Header {
			if len(values) > 0 {
				inputHeaders[strings.ToLower(key)] = strings.Join(values, ",")
			}
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"matched":       result.Matched,
				"responseType":  responseType,
				"matchedRule":   result.MatchedRule,
				"inputHeaders":  inputHeaders,
				"outputHeaders": outputHeaders,
			},
		})
	}
}

func validateResponseRules(config subscriptionresponserules.Config) error {
	if strings.TrimSpace(config.Version) == "" {
		return errors.New("version is required")
	}
	for _, rule := range config.Rules {
		for _, condition := range rule.Conditions {
			switch strings.ToUpper(strings.TrimSpace(condition.Operator)) {
			case "REGEX", "NOT_REGEX":
				if _, err := regexp.Compile(condition.Value); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
