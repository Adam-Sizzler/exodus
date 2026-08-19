package keygen

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
)

func TestKeygenHandler_MethodNotAllowed(t *testing.T) {
	cfg := &config.BackendConfig{}
	handler := KeygenHandler(nil, cfg)

	req := httptest.NewRequest(http.MethodPost, "/api/keygen", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestKeygenResponseJSONSchema(t *testing.T) {
	resp := KeygenResponse{
		Response: KeygenPayload{
			SecretKey: "base64-encoded-secret",
			GrpcToken: "grpc-token-sample",
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal KeygenResponse: %v", err)
	}

	var parsed map[string]map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	res, ok := parsed["response"]
	if !ok {
		t.Fatalf("expected response field, got %v", parsed)
	}

	if res["secretKey"] != "base64-encoded-secret" {
		t.Fatalf("expected secretKey 'base64-encoded-secret', got %v", res["secretKey"])
	}

	if res["grpcToken"] != "grpc-token-sample" {
		t.Fatalf("expected grpcToken 'grpc-token-sample', got %v", res["grpcToken"])
	}

	if _, exists := res["pubKey"]; exists {
		t.Fatalf("expected pubKey to NOT exist in response schema, got %v", res["pubKey"])
	}
}
