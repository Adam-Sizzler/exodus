package nodes

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestDecodeOptionalRestartNodesRequest(t *testing.T) {
	tests := []struct {
		name      string
		body      io.Reader
		wantForce bool
		wantErr   bool
	}{
		{name: "nil body", body: nil, wantForce: false},
		{name: "empty body", body: strings.NewReader(""), wantForce: false},
		{name: "empty object", body: strings.NewReader(`{}`), wantForce: false},
		{name: "force restart true", body: strings.NewReader(`{"forceRestart":true}`), wantForce: true},
		{name: "force restart false", body: strings.NewReader(`{"forceRestart":false}`), wantForce: false},
		{name: "invalid json", body: strings.NewReader(`{`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/api/nodes/actions/restart-all", tt.body)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}
			if tt.body == nil {
				req.Body = nil
			}

			got, err := decodeOptionalRestartNodesRequest(req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if isForceRestartRequested(got) != tt.wantForce {
				t.Fatalf("force restart = %v, want %v", isForceRestartRequested(got), tt.wantForce)
			}
		})
	}
}

func TestNormalizeAPISchema(t *testing.T) {
	tests := map[string]string{
		"":      "mtls",
		"mtls":  "mtls",
		"MTLS":  "mtls",
		"grpc":  "mtls",
		"grpcs": "mtls",
		"https": "mtls",
		"tls":   "tls",
		"bad":   "mtls",
	}

	for raw, want := range tests {
		t.Run(raw, func(t *testing.T) {
			if got := normalizeAPISchema(&raw); got != want {
				t.Fatalf("normalizeAPISchema(%q) = %q, want %q", raw, got, want)
			}
		})
	}

	if got := normalizeAPISchema(nil); got != "mtls" {
		t.Fatalf("normalizeAPISchema(nil) = %q, want mtls", got)
	}
}

func TestNormalizeAPIPath(t *testing.T) {
	tests := map[string]string{
		"":           "/",
		"/":          "/",
		"node":       "/node",
		"/node":      "/node",
		"/node/":     "/node",
		" /grpc/v1 ": "/grpc/v1",
	}

	for raw, want := range tests {
		t.Run(raw, func(t *testing.T) {
			if got := normalizeAPIPath(&raw); got != want {
				t.Fatalf("normalizeAPIPath(%q) = %q, want %q", raw, got, want)
			}
		})
	}

	if got := normalizeAPIPath(nil); got != "/" {
		t.Fatalf("normalizeAPIPath(nil) = %q, want /", got)
	}
}

func TestNodeMultiplierRoundTrip(t *testing.T) {
	values := []float64{0, 1, 1.25, 2.5, 100}
	for _, value := range values {
		t.Run("", func(t *testing.T) {
			if got := fromNanoMultiplier(toNanoMultiplier(value)); got != value {
				t.Fatalf("round trip = %v, want %v", got, value)
			}
		})
	}
}

func TestValidateProxyURLMatchesFrontendContract(t *testing.T) {
	valid := []string{
		"socks5://127.0.0.1:1080",
		"socks5://user:pass@example.com:1080",
		"socks5://user@example.com:1080",
	}
	for _, value := range valid {
		t.Run(value, func(t *testing.T) {
			if err := validateProxyURL(&value); err != nil {
				t.Fatalf("expected valid proxyUrl, got %v", err)
			}
		})
	}

	invalid := []string{
		"http://127.0.0.1:1080",
		"socks5://127.0.0.1",
		"socks5://user:pass@example.com",
	}
	for _, value := range invalid {
		t.Run(value, func(t *testing.T) {
			if err := validateProxyURL(&value); err == nil {
				t.Fatal("expected invalid proxyUrl")
			}
		})
	}
}

func TestValidateNode28Fields(t *testing.T) {
	note := strings.Repeat("a", 255)
	if err := validateNote(&note); err != nil {
		t.Fatalf("expected note at 255 chars to be valid, got %v", err)
	}

	tooLongNote := strings.Repeat("a", 256)
	if err := validateNote(&tooLongNote); err == nil {
		t.Fatal("expected note over 255 chars to be invalid")
	}

	multiplier := 100.1
	err := validateCreateRequest(createNodeRequest{
		Name:                      "node-a",
		Address:                   "127.0.0.1",
		NodeConsumptionMultiplier: &multiplier,
		ConfigProfile: configProfileRefRequest{
			ActiveConfigProfileUUID: "00000000-0000-4000-8000-000000000000",
			ActiveInbounds:          []string{"00000000-0000-4000-8000-000000000001"},
		},
	})
	if err == nil {
		t.Fatal("expected nodeConsumptionMultiplier over 100 to be invalid")
	}
}

func TestBuildNodeVersions(t *testing.T) {
	if got := buildNodeVersions(nil, nil); got != nil {
		t.Fatalf("buildNodeVersions(nil, nil) = %#v, want nil", got)
	}

	singbox := "1.13.5"
	empty := ""
	got := buildNodeVersions(&singbox, &empty)
	if got == nil {
		t.Fatal("expected versions")
	}
	if got.Singbox != "1.13.5" || got.Node != "unknown" {
		t.Fatalf("versions = %#v", got)
	}
}
