package server

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractSingboxLogReason(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "current")

	content := "@4000000066c8919114f828a4 2026-08-23 12:48:01 INFO sing-box starting\n@4000000066c8919114f828a5 \x1b[31mFATAL[0000] create service: listen tcp 127.0.0.1:10086: bind: address already in use\x1b[0m\n"
	if err := os.WriteFile(logFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	reason := ExtractSingboxLogReason(logFile, 10)
	want := "create service: listen tcp 127.0.0.1:10086: bind: address already in use"
	if reason != want {
		t.Fatalf("reason got %q, want %q", reason, want)
	}
}

func TestExtractSingboxLogReasonPanic(t *testing.T) {
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "current")

	content := `2026-08-23 12:48:01 INFO starting
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation]
`
	if err := os.WriteFile(logFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write log file: %v", err)
	}

	reason := ExtractSingboxLogReason(logFile, 10)
	if reason != "panic: runtime error: invalid memory address or nil pointer dereference" {
		t.Fatalf("unexpected reason: %q", reason)
	}
}
