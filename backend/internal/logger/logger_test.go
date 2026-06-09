package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestResolveLevelDebugEnvEnablesDebug(t *testing.T) {
	if got := ResolveLevel("production", "true", ""); got != LevelDebug {
		t.Fatalf("ResolveLevel() = %v, want debug", got)
	}
}

func TestResolveLevelDefaultsToInfo(t *testing.T) {
	if got := ResolveLevel("", "", ""); got != LevelInfo {
		t.Fatalf("ResolveLevel() = %v, want info", got)
	}
}

func TestLoggerFormatIncludesContextAndTimestamp(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Options{Writer: &buf, Level: "debug"}).WithContext("RootService")
	logger.Log("hello")

	line := buf.String()
	if !strings.Contains(line, " LOG [RootService] hello") {
		t.Fatalf("unexpected line: %q", line)
	}
	if !strings.Contains(line, "+0ms") {
		t.Fatalf("expected ms suffix, got: %q", line)
	}
}

func TestLoggerFiltersNestBootstrapContexts(t *testing.T) {
	var buf bytes.Buffer
	logger := New(Options{Writer: &buf, Level: "debug"}).WithContext("RouterExplorer")
	logger.Log("hidden")
	if buf.Len() != 0 {
		t.Fatalf("expected filtered log, got: %q", buf.String())
	}
}
