package config

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestResolveExodusLogLevelDevelopment(t *testing.T) {
	t.Setenv("NODE_ENV", "development")
	if got := ResolveExodusLogLevel(""); got != "debug" {
		t.Fatalf("ResolveExodusLogLevel() = %q, want debug", got)
	}
}

func TestResolveExodusLogLevelDebugFlag(t *testing.T) {
	t.Setenv("NODE_ENV", "production")
	t.Setenv("ENABLE_DEBUG_LOGS", "true")
	if got := ResolveExodusLogLevel(""); got != "debug" {
		t.Fatalf("ResolveExodusLogLevel() = %q, want debug", got)
	}
}

func TestResolveExodusLogLevelDefaultsToInfo(t *testing.T) {
	t.Setenv("NODE_ENV", "")
	t.Setenv("ENV", "")
	t.Setenv("ENABLE_DEBUG_LOGS", "")
	if got := ResolveExodusLogLevel(""); got != "info" {
		t.Fatalf("ResolveExodusLogLevel() = %q, want info", got)
	}
}

func TestResolveExodusLogLevelDevelopmentWinsOverConfiguredLevel(t *testing.T) {
	t.Setenv("NODE_ENV", "development")
	if got := ResolveExodusLogLevel("error"); got != "debug" {
		t.Fatalf("ResolveExodusLogLevel(error) = %q, want debug", got)
	}
}

func TestResolveExodusLogLevelAllowsConfiguredDebug(t *testing.T) {
	t.Setenv("NODE_ENV", "production")
	t.Setenv("ENABLE_DEBUG_LOGS", "")
	if got := ResolveExodusLogLevel("debug"); got != "debug" {
		t.Fatalf("ResolveExodusLogLevel(debug) = %q, want debug", got)
	}
}

func TestExodusLoggerFormatIncludesContext(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	t.Setenv("LOG_FORMAT", "console")
	var buf bytes.Buffer
	logger := NewExodusLogger(&buf, "debug").WithContext("SingboxService")
	logger.Log("Sing-box Core configuration is up-to-date - no restart required")

	line := strings.TrimSpace(buf.String())
	pattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3} INFO \[SingboxService\] Sing-box Core configuration is up-to-date - no restart required$`)
	if !pattern.MatchString(line) {
		t.Fatalf("unexpected log line: %q", line)
	}
}
