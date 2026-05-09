package main

import "testing"

func TestResolveNodeVersion(t *testing.T) {
	prev := buildVersion
	defer func() { buildVersion = prev }()

	buildVersion = "26.1.13"

	testCases := []struct {
		name          string
		configVersion string
		expected      string
	}{
		{name: "config plain semver", configVersion: "26.1.13", expected: "v26.1.13"},
		{name: "config prefixed semver", configVersion: "v26.1.13", expected: "v26.1.13"},
		{name: "config uppercase prefixed semver", configVersion: "V26.1.13", expected: "v26.1.13"},
		{name: "fallback to build version", configVersion: "", expected: "v26.1.13"},
		{name: "unknown sentinel falls back to build version", configVersion: "unknown", expected: "v26.1.13"},
		{name: "release tag suffix is preserved", configVersion: "v26.5.8.s", expected: "v26.5.8.s"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := resolveNodeVersion(testCase.configVersion)
			if got != testCase.expected {
				t.Fatalf("unexpected version: got %q want %q", got, testCase.expected)
			}
		})
	}
}

func TestResolveNodeVersionDevFallback(t *testing.T) {
	prev := buildVersion
	defer func() { buildVersion = prev }()

	buildVersion = "unknown"
	if got := resolveNodeVersion("unknown"); got != "v0.0.0-dev" {
		t.Fatalf("unexpected dev fallback: got %q want %q", got, "v0.0.0-dev")
	}
}
