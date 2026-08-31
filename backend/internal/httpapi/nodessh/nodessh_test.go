package nodessh

import (
"bytes"
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"exodus/internal/config"
"exodus/internal/httpapi/auth"
)

func TestEvaluateBlindedElement(t *testing.T) {
appSecret := "test-secret"
blindedB64 := "yI6YVYxJJ/ifiuszVbf149kx6v/mfv21+qTGKJ8exQY="
expectedEvaluatedB64 := "4krWufU/eG5zjcGzyEK/Qubc7ksxyZuMvCv63UF3wEE="

evaluated, err := EvaluateBlindedElement(appSecret, blindedB64)
if err != nil {
t.Fatalf("EvaluateBlindedElement failed: %v", err)
}

if evaluated != expectedEvaluatedB64 {
t.Fatalf("expected evaluated %s, got %s", expectedEvaluatedB64, evaluated)
}
}

func TestNodeSSHVaultEvaluateHandler(t *testing.T) {
cfg := &config.BackendConfig{
JWT: config.JWTConfig{
AuthSecret: "test-secret",
},
}

reqBody, _ := json.Marshal(EvaluateVaultRequestBody{
Blinded: "yI6YVYxJJ/ifiuszVbf149kx6v/mfv21+qTGKJ8exQY=",
})

req := httptest.NewRequest(http.MethodPost, "/api/node-ssh/vault/evaluate", bytes.NewReader(reqBody))
req.Header.Set("Content-Type", "application/json")

// Set admin principal context
ctx := auth.WithAuthPrincipal(req.Context(), &auth.AuthPrincipal{
AdminUUID: "00000000-0000-0000-0000-000000000001",
Role:      "ADMIN",
TokenType: "jwt_auth",
})
req = req.WithContext(ctx)

rr := httptest.NewRecorder()
handler := NodeSSHVaultEvaluateHandler(nil, cfg)
handler.ServeHTTP(rr, req)

if rr.Code != http.StatusOK {
t.Fatalf("expected status 200, got %d: %s", rr.Code, rr.Body.String())
}

var res struct {
Response EvaluateVaultResponse `json:"response"`
}
if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
t.Fatalf("failed to decode response: %v", err)
}

expected := "4krWufU/eG5zjcGzyEK/Qubc7ksxyZuMvCv63UF3wEE="
if res.Response.Evaluated != expected {
t.Fatalf("expected evaluated %s, got %s", expected, res.Response.Evaluated)
}
}
