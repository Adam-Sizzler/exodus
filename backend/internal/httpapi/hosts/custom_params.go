package hosts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// --- Singbox Structures ---

type SingboxVLESSCustomParams struct {
	PacketEncoding *string `json:"packet_encoding,omitempty"`
}

type SingboxAnyTLSCustomParams struct {
	IdleSessionCheckInterval *string `json:"idle_session_check_interval,omitempty"`
	IdleSessionTimeout       *string `json:"idle_session_timeout,omitempty"`
	MinIdleSession           *int    `json:"min_idle_session,omitempty"`
}

type SingboxHysteria2Obfs struct {
	Type     string `json:"type"`
	Password string `json:"password"`
}

type SingboxHysteria2CustomParams struct {
	UpMbps      *int                  `json:"up_mbps,omitempty"`
	DownMbps    *int                  `json:"down_mbps,omitempty"`
	ServerPorts []string              `json:"server_ports,omitempty"`
	HopInterval *string               `json:"hop_interval,omitempty"`
	Obfs        *SingboxHysteria2Obfs `json:"obfs,omitempty"`
}

type SingboxNaiveCustomParams struct {
	QUIC                  *bool             `json:"quic,omitempty"`
	QUICCongestionControl *string           `json:"quic_congestion_control,omitempty"`
	InsecureConcurrency   *int              `json:"insecure_concurrency,omitempty"`
	ExtraHeaders          map[string]string `json:"extra_headers,omitempty"`
}

type SingboxGRPCTransportCustomParams struct {
	IdleTimeout         *string `json:"idle_timeout,omitempty"`
	PingTimeout         *string `json:"ping_timeout,omitempty"`
	PermitWithoutStream *bool   `json:"permit_without_stream,omitempty"`
}

// Singbox Custom Params Wrapper to validate multiple blocks from a single JSON
// Wait, the chat says: "Формат хранения — просто плоский JSON (Singbox)... без обёртки с именем протокола."
// So if the host is VLESS + gRPC, the JSON contains BOTH VLESS fields and GRPC transport fields.
// Therefore, we need to decode the JSON into the specific protocol struct AND the specific transport struct.

func ValidateSingboxCustomParams(protocol string, network string, rawJSON []byte) error {
	rawJSONStr := strings.TrimSpace(string(rawJSON))
	if rawJSONStr == "" || rawJSONStr == "{}" || rawJSONStr == "null" {
		return nil
	}

	// Since json.Decoder with DisallowUnknownFields will fail if there are fields not in the struct,
	// and we have a flat JSON containing BOTH protocol fields and transport fields,
	// we must create a combined struct on the fly, or just unmarshal into a combined struct.
	// Since we know the combinations, let's build a combined type dynamically using maps, or just defined structs.

	type combinedSingbox struct {
		// VLESS
		PacketEncoding *string `json:"packet_encoding,omitempty"`
		// AnyTLS
		IdleSessionCheckInterval *string `json:"idle_session_check_interval,omitempty"`
		IdleSessionTimeout       *string `json:"idle_session_timeout,omitempty"`
		MinIdleSession           *int    `json:"min_idle_session,omitempty"`
		// Hysteria2
		UpMbps      *int                  `json:"up_mbps,omitempty"`
		DownMbps    *int                  `json:"down_mbps,omitempty"`
		ServerPorts []string              `json:"server_ports,omitempty"`
		HopInterval *string               `json:"hop_interval,omitempty"`
		Obfs        *SingboxHysteria2Obfs `json:"obfs,omitempty"`
		// Naive
		QUIC                  *bool             `json:"quic,omitempty"`
		QUICCongestionControl *string           `json:"quic_congestion_control,omitempty"`
		InsecureConcurrency   *int              `json:"insecure_concurrency,omitempty"`
		ExtraHeaders          map[string]string `json:"extra_headers,omitempty"`
		// GRPC Transport
		Transport *SingboxGRPCTransportCustomParams `json:"transport,omitempty"`
	}

	// We decode into a map first to check allowed fields, or we use a strict map approach.
	var parsed map[string]interface{}
	dec := json.NewDecoder(bytes.NewReader(rawJSON))
	if err := dec.Decode(&parsed); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}

	allowedKeys := make(map[string]bool)

	switch protocol {
	case "vless":
		allowedKeys["packet_encoding"] = true
	case "anytls":
		allowedKeys["idle_session_check_interval"] = true
		allowedKeys["idle_session_timeout"] = true
		allowedKeys["min_idle_session"] = true
	case "hysteria2":
		allowedKeys["up_mbps"] = true
		allowedKeys["down_mbps"] = true
		allowedKeys["server_ports"] = true
		allowedKeys["hop_interval"] = true
		allowedKeys["obfs"] = true
	case "naive":
		allowedKeys["quic"] = true
		allowedKeys["quic_congestion_control"] = true
		allowedKeys["insecure_concurrency"] = true
		allowedKeys["extra_headers"] = true
	}

	if network == "grpc" {
		allowedKeys["transport"] = true
	}

	for k := range parsed {
		if !allowedKeys[k] {
			return fmt.Errorf("поле '%s' не относится к протоколу %s (транспорт %s)", k, protocol, network)
		}
	}

	// Additionally parse into combined to validate types (e.g. Obfs struct)
	decStrict := json.NewDecoder(bytes.NewReader(rawJSON))
	decStrict.DisallowUnknownFields() // just to be safe
	var typed combinedSingbox
	if err := decStrict.Decode(&typed); err != nil {
		return fmt.Errorf("type validation failed: %w", err)
	}

	// Custom enum validation
	if typed.Obfs != nil && typed.Obfs.Type != "salamander" {
		return fmt.Errorf("obfs type must be 'salamander'")
	}
	if typed.PacketEncoding != nil {
		pe := *typed.PacketEncoding
		if pe != "" && pe != "xudp" && pe != "packetaddr" {
			return fmt.Errorf("packet_encoding must be '', 'xudp' or 'packetaddr'")
		}
	}

	return nil
}

// --- Mihomo Structures ---

type MihomoVLESSCustomParams struct {
	PacketEncoding *string `yaml:"packet-encoding,omitempty"`
}

type MihomoAnyTLSCustomParams struct {
	IdleSessionCheckInterval *int `yaml:"idle-session-check-interval,omitempty"`
	IdleSessionTimeout       *int `yaml:"idle-session-timeout,omitempty"`
	MinIdleSession           *int `yaml:"min-idle-session,omitempty"`
}

type MihomoHysteria2CustomParams struct {
	Up           *string `yaml:"up,omitempty"`
	Down         *string `yaml:"down,omitempty"`
	Ports        *string `yaml:"ports,omitempty"`
	HopInterval  *int    `yaml:"hop-interval,omitempty"`
	Obfs         *string `yaml:"obfs,omitempty"`
	ObfsPassword *string `yaml:"obfs-password,omitempty"`
}

type MihomoGRPCTransportCustomParams struct {
	PingInterval   *int `yaml:"ping-interval,omitempty"`
	MaxConnections *int `yaml:"max-connections,omitempty"`
	MinStreams     *int `yaml:"min-streams,omitempty"`
	MaxStreams     *int `yaml:"max-streams,omitempty"`
}

func ValidateMihomoCustomParams(protocol string, network string, rawYAML []byte) error {
	rawYAMLStr := strings.TrimSpace(string(rawYAML))
	if rawYAMLStr == "" || rawYAMLStr == "{}" || rawYAMLStr == "null" {
		return nil
	}

	type combinedMihomo struct {
		// VLESS
		PacketEncoding *string `yaml:"packet-encoding,omitempty"`
		// AnyTLS
		IdleSessionCheckInterval *int `yaml:"idle-session-check-interval,omitempty"`
		IdleSessionTimeout       *int `yaml:"idle-session-timeout,omitempty"`
		MinIdleSession           *int `yaml:"min-idle-session,omitempty"`
		// Hysteria2
		Up           *string `yaml:"up,omitempty"`
		Down         *string `yaml:"down,omitempty"`
		Ports        *string `yaml:"ports,omitempty"`
		HopInterval  *int    `yaml:"hop-interval,omitempty"`
		Obfs         *string `yaml:"obfs,omitempty"`
		ObfsPassword *string `yaml:"obfs-password,omitempty"`
		// GRPC Transport
		GrpcOpts *MihomoGRPCTransportCustomParams `yaml:"grpc-opts,omitempty"`
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal(rawYAML, &parsed); err != nil {
		return fmt.Errorf("invalid yaml: %w", err)
	}

	allowedKeys := make(map[string]bool)

	switch protocol {
	case "vless":
		allowedKeys["packet-encoding"] = true
	case "anytls":
		allowedKeys["idle-session-check-interval"] = true
		allowedKeys["idle-session-timeout"] = true
		allowedKeys["min-idle-session"] = true
	case "hysteria2":
		allowedKeys["up"] = true
		allowedKeys["down"] = true
		allowedKeys["ports"] = true
		allowedKeys["hop-interval"] = true
		allowedKeys["obfs"] = true
		allowedKeys["obfs-password"] = true
	}

	if network == "grpc" {
		allowedKeys["grpc-opts"] = true
	}

	for k := range parsed {
		if !allowedKeys[k] {
			return fmt.Errorf("поле '%s' не относится к протоколу %s (транспорт %s)", k, protocol, network)
		}
	}

	var typed combinedMihomo
	dec := yaml.NewDecoder(bytes.NewReader(rawYAML))
	dec.KnownFields(true)
	if err := dec.Decode(&typed); err != nil {
		return fmt.Errorf("type validation failed: %w", err)
	}

	if typed.Obfs != nil && *typed.Obfs != "salamander" {
		return fmt.Errorf("obfs must be 'salamander'")
	}
	if typed.PacketEncoding != nil {
		pe := *typed.PacketEncoding
		if pe != "" && pe != "xudp" && pe != "packetaddr" {
			return fmt.Errorf("packet-encoding must be '', 'xudp' or 'packetaddr'")
		}
	}

	return nil
}
