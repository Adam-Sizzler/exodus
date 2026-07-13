package passkeys

import (
	"encoding/base64"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
)

func encodeCredentialID(raw []byte) string {
	return base64.RawURLEncoding.EncodeToString(raw)
}

func decodeCredentialID(id string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(strings.TrimSpace(id))
}

func parseTransports(raw string) []protocol.AuthenticatorTransport {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []protocol.AuthenticatorTransport{}
	}
	parts := strings.Split(raw, ",")
	res := make([]protocol.AuthenticatorTransport, 0, len(parts))
	for _, p := range parts {
		t := protocol.AuthenticatorTransport(strings.ToLower(strings.TrimSpace(p)))
		if t != "" {
			res = append(res, t)
		}
	}
	return res
}

func transportsToCSV(transports []protocol.AuthenticatorTransport) string {
	if len(transports) == 0 {
		return ""
	}
	parts := make([]string, len(transports))
	for i, t := range transports {
		parts[i] = string(t)
	}
	return strings.Join(parts, ",")
}

func passkeyProviderName(transports []protocol.AuthenticatorTransport) string {
	if len(transports) == 0 {
		return "Passkey"
	}
	for _, t := range transports {
		switch t {
		case protocol.USB, protocol.NFC, protocol.BLE:
			return "Hardware Key (YubiKey / Titan)"
		case protocol.Internal:
			return "Platform Authenticator (TouchID / FaceID / Windows Hello)"
		case protocol.Hybrid:
			return "Mobile Device / Passkey"
		}
	}
	return "Passkey"
}

func uint32Counter(counter int64) uint32 {
	if counter < 0 {
		return 0
	}
	if counter > 4294967295 {
		return 4294967295
	}
	return uint32(counter)
}

func currentRequestOrigin(r *http.Request) string {
	if r == nil {
		return ""
	}

	proto := "https"
	if r.TLS == nil {
		proto = "http"
	}
	if forwardedProto := firstHeaderValue(r.Header.Get("X-Forwarded-Proto")); forwardedProto != "" {
		proto = strings.ToLower(forwardedProto)
	}

	host := strings.TrimSpace(r.Host)
	if host == "" {
		return ""
	}

	return proto + "://" + host
}

func normalizeOrigin(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return ""
	}

	host, port, _ := net.SplitHostPort(parsed.Host)
	if host == "" {
		host = parsed.Host
	}

	if (parsed.Scheme == "http" && port == "80") || (parsed.Scheme == "https" && port == "443") {
		port = ""
	}

	result := parsed.Scheme + "://" + host
	if port != "" {
		result += ":" + port
	}
	return result
}

func normalizeRPID(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.Contains(raw, "://") {
		parsed, err := url.Parse(raw)
		if err == nil && parsed.Host != "" {
			raw = parsed.Host
		}
	}

	host, _, _ := net.SplitHostPort(raw)
	if host != "" {
		return host
	}
	return raw
}

func appendOrigin(origins []string, origin string) []string {
	norm := normalizeOrigin(origin)
	if norm == "" {
		return origins
	}

	for _, o := range origins {
		if o == norm {
			return origins
		}
	}
	return append(origins, norm)
}

func appendLocalOrigins(origins []string, host, port string) []string {
	if isLocalHost(host) {
		origins = appendOrigin(origins, "http://localhost:"+port)
		origins = appendOrigin(origins, "http://127.0.0.1:"+port)
		origins = appendOrigin(origins, "https://localhost:"+port)
		origins = appendOrigin(origins, "https://127.0.0.1:"+port)
	}
	return origins
}

func sameOriginHost(a, b string) bool {
	parsedA, errA := url.Parse(normalizeOrigin(a))
	parsedB, errB := url.Parse(normalizeOrigin(b))
	if errA != nil || errB != nil {
		return false
	}

	hostA, _, _ := net.SplitHostPort(parsedA.Host)
	if hostA == "" {
		hostA = parsedA.Host
	}

	hostB, _, _ := net.SplitHostPort(parsedB.Host)
	if hostB == "" {
		hostB = parsedB.Host
	}

	return strings.EqualFold(hostA, hostB)
}

func isLocalHost(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func splitHostPortLoose(value string) (host string, port string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", ""
	}

	h, p, err := net.SplitHostPort(value)
	if err != nil {
		return value, ""
	}
	return h, p
}

func firstHeaderValue(value string) string {
	parts := strings.Split(value, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}
