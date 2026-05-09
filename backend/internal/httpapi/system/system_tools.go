package system

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	subscriptionresponserules "exodus/internal/httpapi/subscriptionresponserules"

	"golang.org/x/crypto/curve25519"
)

const happCryptoV4PublicKey = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA3UZ0M3L4K+WjM3vkbQnz
ozHg/cRbEXvQ6i4A8RVN4OM3rK9kU01FdjyoIgywve8OEKsFnVwERZAQZ1Trv60B
hmaM76QQEE+EUlIOL9EpwKWGtTL5lYC1sT9XJMNP3/CI0gP5wwQI88cY/xedpOEB
W72EmOOShHUm/b/3m+HPmqwc4ugKj5zWV5SyiT829aFA5DxSjmIIFBAms7DafmSq
LFTYIQL5cShDY2u+/sqyAw9yZIOoqW2TFIgIHhLPWek/ocDU7zyOrlu1E0SmcQQb
LFqHq02fsnH6IcqTv3N5Adb/CkZDDQ6HvQVBmqbKZKf7ZdXkqsc/Zw27xhG7OfXC
tUmWsiL7zA+KoTd3avyOh93Q9ju4UQsHthL3Gs4vECYOCS9dsXXSHEY/1ngU/hjO
WFF8QEE/rYV6nA4PTyUvo5RsctSQL/9DJX7XNh3zngvif8LsCN2MPvx6X+zLouBX
zgBkQ9DFfZAGLWf9TR7KVjZC/3NsuUCDoAOcpmN8pENBbeB0puiKMMWSvll36+2M
YR1Xs0MgT8Y9TwhE2+TnnTJOhzmHi/BxiUlY/w2E0s4ax9GHAmX0wyF4zeV7kDkc
vHuEdc0d7vDmdw0oqCqWj0Xwq86HfORu6tm1A8uRATjb4SzjTKclKuoElVAVa5Jo
oh/uZMozC65SmDw+N5p6Su8CAwEAAQ==
-----END PUBLIC KEY-----
`

type encryptHappRequest struct {
	LinkToEncrypt string `json:"linkToEncrypt"`
}

type testSRRMatcherRequest struct {
	ResponseRules subscriptionresponserules.Config `json:"responseRules"`
}

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
				shared.SendError(w, http.StatusInternalServerError, "failed to generate x25519 keypair", err, cfg)
				return
			}
			publicKey, err := curve25519.X25519(privateKey, curve25519.Basepoint)
			if err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to generate x25519 keypair", err, cfg)
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

func EncryptHappCryptoLinkHandler(cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req encryptHappRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if strings.TrimSpace(req.LinkToEncrypt) == "" {
			shared.SendError(w, http.StatusBadRequest, "linkToEncrypt is required", nil, cfg)
			return
		}
		encrypted, err := encryptHappV4(req.LinkToEncrypt)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to encrypt Happ crypto link", err, cfg)
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"encryptedLink": encrypted},
		})
	}
}

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

func encryptHappV4(content string) (string, error) {
	block, _ := pem.Decode([]byte(happCryptoV4PublicKey))
	if block == nil {
		return "", errors.New("invalid Happ public key PEM")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pub, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("Happ public key is not RSA")
	}
	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(content))
	if err != nil {
		return "", err
	}
	return "happ://crypt4/" + base64.StdEncoding.EncodeToString(encrypted), nil
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
