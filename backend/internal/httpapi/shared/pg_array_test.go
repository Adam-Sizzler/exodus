package shared

import (
	"reflect"
	"testing"
)

func TestPostgresTextArrayLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: "{}",
		},
		{
			name:     "single element",
			input:    []string{"announce"},
			expected: `{"announce"}`,
		},
		{
			name:     "multiple elements",
			input:    []string{"announce", "routing", "support-url"},
			expected: `{"announce","routing","support-url"}`,
		},
		{
			name:     "elements with quotes and backslashes",
			input:    []string{`test"val`, `test\val`},
			expected: `{"test\"val","test\\val"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := PostgresTextArrayLiteral(tt.input)
			if actual != tt.expected {
				t.Fatalf("expected %s, got %s", tt.expected, actual)
			}
		})
	}
}

func TestParsePgTextArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "empty pg array",
			input:    "{}",
			expected: []string{},
		},
		{
			name:     "empty json array",
			input:    "[]",
			expected: []string{},
		},
		{
			name:     "pg array literal",
			input:    "{announce,routing}",
			expected: []string{"announce", "routing"},
		},
		{
			name:     "quoted pg array literal",
			input:    `{"announce","routing","support-url"}`,
			expected: []string{"announce", "routing", "support-url"},
		},
		{
			name:     "json array string",
			input:    `["announce","routing","support-url"]`,
			expected: []string{"announce", "routing", "support-url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParsePgTextArray(tt.input)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
