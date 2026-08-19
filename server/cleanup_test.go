package server

import (
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupSocketFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cleanup_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sockPath1 := filepath.Join(tmpDir, "test1.sock")
	sockPath2 := filepath.Join(tmpDir, "test2.sock")
	regularPath := filepath.Join(tmpDir, "regular.txt")

	// Create unix domain sockets
	l1, err := net.Listen("unix", sockPath1)
	if err != nil {
		t.Fatalf("failed to listen on unix socket 1: %v", err)
	}
	_ = l1.Close() // closed listener leaves stale socket file

	l2, err := net.Listen("unix", sockPath2)
	if err != nil {
		t.Fatalf("failed to listen on unix socket 2: %v", err)
	}
	_ = l2.Close()

	// Create regular file
	if err := os.WriteFile(regularPath, []byte("keep me"), 0644); err != nil {
		t.Fatalf("failed to create regular file: %v", err)
	}

	// Run cleanup on *.sock
	cleanupSocketFiles(nil, []string{filepath.Join(tmpDir, "*.sock")})

	// Check that socket files were removed
	if _, err := os.Stat(sockPath1); !os.IsNotExist(err) {
		t.Errorf("expected sockPath1 to be removed, but stat returned %v", err)
	}
	if _, err := os.Stat(sockPath2); !os.IsNotExist(err) {
		t.Errorf("expected sockPath2 to be removed, but stat returned %v", err)
	}

	// Check that regular file was preserved
	if _, err := os.Stat(regularPath); err != nil {
		t.Errorf("expected regularPath to be preserved, but got %v", err)
	}
}
