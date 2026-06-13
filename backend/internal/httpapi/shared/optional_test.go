package shared

import (
	"encoding/json"
	"testing"
)

func TestOptionalStringUnmarshalJSON(t *testing.T) {
	type payload struct {
		Value OptionalString `json:"value,omitempty"`
	}

	tests := []struct {
		name      string
		raw       string
		wantSet   bool
		wantValue *string
		wantErr   bool
	}{
		{name: "missing field", raw: `{}`, wantSet: false, wantValue: nil},
		{name: "null value", raw: `{"value":null}`, wantSet: true, wantValue: nil},
		{name: "empty string", raw: `{"value":""}`, wantSet: true, wantValue: stringPtr("")},
		{name: "string value", raw: `{"value":"alpha"}`, wantSet: true, wantValue: stringPtr("alpha")},
		{name: "invalid type", raw: `{"value":123}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got payload
			err := json.Unmarshal([]byte(tt.raw), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value.Set != tt.wantSet {
				t.Fatalf("Set = %v, want %v", got.Value.Set, tt.wantSet)
			}
			assertStringPtr(t, got.Value.Value, tt.wantValue)
		})
	}
}

func TestOptionalIntUnmarshalJSON(t *testing.T) {
	type payload struct {
		Value OptionalInt `json:"value,omitempty"`
	}

	tests := []struct {
		name      string
		raw       string
		wantSet   bool
		wantValue *int
		wantErr   bool
	}{
		{name: "missing field", raw: `{}`, wantSet: false, wantValue: nil},
		{name: "null value", raw: `{"value":null}`, wantSet: true, wantValue: nil},
		{name: "zero value", raw: `{"value":0}`, wantSet: true, wantValue: intPtr(0)},
		{name: "int value", raw: `{"value":7}`, wantSet: true, wantValue: intPtr(7)},
		{name: "invalid type", raw: `{"value":"7"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got payload
			err := json.Unmarshal([]byte(tt.raw), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value.Set != tt.wantSet {
				t.Fatalf("Set = %v, want %v", got.Value.Set, tt.wantSet)
			}
			assertIntPtr(t, got.Value.Value, tt.wantValue)
		})
	}
}

func TestOptionalInt64UnmarshalJSON(t *testing.T) {
	type payload struct {
		Value OptionalInt64 `json:"value,omitempty"`
	}

	tests := []struct {
		name      string
		raw       string
		wantSet   bool
		wantValue *int64
		wantErr   bool
	}{
		{name: "missing field", raw: `{}`, wantSet: false, wantValue: nil},
		{name: "null value", raw: `{"value":null}`, wantSet: true, wantValue: nil},
		{name: "zero value", raw: `{"value":0}`, wantSet: true, wantValue: int64Ptr(0)},
		{name: "int64 value", raw: `{"value":9223372036854775807}`, wantSet: true, wantValue: int64Ptr(9223372036854775807)},
		{name: "invalid type", raw: `{"value":"7"}`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got payload
			err := json.Unmarshal([]byte(tt.raw), &got)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value.Set != tt.wantSet {
				t.Fatalf("Set = %v, want %v", got.Value.Set, tt.wantSet)
			}
			assertInt64Ptr(t, got.Value.Value, tt.wantValue)
		})
	}
}

func TestOptionalJSONUnmarshalJSON(t *testing.T) {
	type payload struct {
		Value OptionalJSON `json:"value,omitempty"`
	}

	tests := []struct {
		name    string
		raw     string
		wantSet bool
		wantRaw string
	}{
		{name: "missing field", raw: `{}`, wantSet: false, wantRaw: ""},
		{name: "null value", raw: `{"value":null}`, wantSet: true, wantRaw: "null"},
		{name: "object value", raw: `{"value":{"enabled":true}}`, wantSet: true, wantRaw: `{"enabled":true}`},
		{name: "array value", raw: `{"value":[1,2]}`, wantSet: true, wantRaw: `[1,2]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got payload
			if err := json.Unmarshal([]byte(tt.raw), &got); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Value.Set != tt.wantSet {
				t.Fatalf("Set = %v, want %v", got.Value.Set, tt.wantSet)
			}
			if string(got.Value.Raw) != tt.wantRaw {
				t.Fatalf("Raw = %q, want %q", string(got.Value.Raw), tt.wantRaw)
			}
		})
	}
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}

func assertStringPtr(t *testing.T, got *string, want *string) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("pointer = %v, want %v", got, want)
		}
		return
	}
	if *got != *want {
		t.Fatalf("value = %q, want %q", *got, *want)
	}
}

func assertIntPtr(t *testing.T, got *int, want *int) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("pointer = %v, want %v", got, want)
		}
		return
	}
	if *got != *want {
		t.Fatalf("value = %d, want %d", *got, *want)
	}
}

func assertInt64Ptr(t *testing.T, got *int64, want *int64) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("pointer = %v, want %v", got, want)
		}
		return
	}
	if *got != *want {
		t.Fatalf("value = %d, want %d", *got, *want)
	}
}
