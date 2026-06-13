package hosts

import (
	"encoding/json"
	"testing"
)

func hostStringPtr(value string) *string {
	return &value
}

func hostBoolPtr(value bool) *bool {
	return &value
}

func TestNormalizeSecurityLayer(t *testing.T) {
	tests := []struct {
		name string
		raw  *string
		want string
	}{
		{name: "nil default", raw: nil, want: "DEFAULT"},
		{name: "empty default", raw: hostStringPtr(""), want: "DEFAULT"},
		{name: "tls upper", raw: hostStringPtr("tls"), want: "TLS"},
		{name: "none upper", raw: hostStringPtr(" NONE "), want: "NONE"},
		{name: "invalid default", raw: hostStringPtr("invalid"), want: "DEFAULT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeSecurityLayer(tt.raw); got != tt.want {
				t.Fatalf("normalizeSecurityLayer() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeJSONField(t *testing.T) {
	tests := []struct {
		name      string
		raw       *json.RawMessage
		emptyNull bool
		wantSet   bool
		wantValue string
		wantErr   bool
	}{
		{name: "nil", raw: nil, emptyNull: true, wantSet: false},
		{name: "null", raw: rawJSON("null"), emptyNull: true, wantSet: true, wantValue: ""},
		{name: "empty object as null", raw: rawJSON(`{}`), emptyNull: true, wantSet: true, wantValue: ""},
		{name: "empty object kept", raw: rawJSON(`{}`), emptyNull: false, wantSet: true, wantValue: `{}`},
		{name: "object", raw: rawJSON(`{"enabled":true}`), emptyNull: true, wantSet: true, wantValue: `{"enabled":true}`},
		{name: "invalid", raw: rawJSON(`{`), emptyNull: true, wantSet: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, gotValue, err := normalizeJSONField(tt.raw, tt.emptyNull)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotSet != tt.wantSet {
				t.Fatalf("set = %v, want %v", gotSet, tt.wantSet)
			}
			if string(gotValue) != tt.wantValue {
				t.Fatalf("value = %q, want %q", string(gotValue), tt.wantValue)
			}
		})
	}
}

func TestValidateProtocolCredentialCreate(t *testing.T) {
	if err := validateProtocolCredentialCreate(hostBoolPtr(true), nil); err == nil {
		t.Fatal("expected required credential error")
	}
	if err := validateProtocolCredentialCreate(hostBoolPtr(true), hostStringPtr("secret")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateProtocolCredentialCreate(hostBoolPtr(false), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateProtocolCredentialUpdate(t *testing.T) {
	if err := validateProtocolCredentialUpdate(hostBoolPtr(true), OptionalString{Set: true, Value: nil}); err == nil {
		t.Fatal("expected required credential error")
	}
	if err := validateProtocolCredentialUpdate(hostBoolPtr(true), OptionalString{Set: true, Value: hostStringPtr("secret")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateProtocolCredentialUpdate(nil, OptionalString{Set: true, Value: hostStringPtr("")}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateTemplateTypes(t *testing.T) {
	if err := validateTemplateTypes([]string{"SINGBOX", "CLASH"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validateTemplateTypes([]string{"UNKNOWN"}); err == nil {
		t.Fatal("expected invalid type error")
	}
}

func TestMapHostRecordToAPIEnsuresSlices(t *testing.T) {
	api := mapHostRecordToAPI(hostRecord{UUID: "uuid", Remark: "remark", Address: "example.com", SecurityLayer: "DEFAULT"}, nil, nil)
	if api.Nodes == nil {
		t.Fatal("Nodes must be empty slice, got nil")
	}
	if api.ExcludedInternalSquads == nil {
		t.Fatal("ExcludedInternalSquads must be empty slice, got nil")
	}
	if api.ExcludeFromSubscription == nil {
		t.Fatal("ExcludeFromSubscription must be empty slice, got nil")
	}
}

func rawJSON(value string) *json.RawMessage {
	raw := json.RawMessage(value)
	return &raw
}
