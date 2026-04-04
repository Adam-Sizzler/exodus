package server

import (
	"bufio"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func startUnixRuntimeServer(t *testing.T, socketPath string, reply string, received chan<- string) {
	t.Helper()

	if err := os.Remove(socketPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("remove stale socket: %v", err)
	}

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		t.Fatalf("listen on unix socket: %v", err)
	}

	t.Cleanup(func() {
		_ = ln.Close()
		_ = os.Remove(socketPath)
	})

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
		line, _ := bufio.NewReader(conn).ReadString('\n')
		received <- line
		_, _ = conn.Write([]byte(reply))
	}()
}

func TestRunHaproxyRuntimeCommandOK(t *testing.T) {
	tmpDir := t.TempDir()
	socketPath := filepath.Join(tmpDir, "haproxy.sock")
	received := make(chan string, 1)

	startUnixRuntimeServer(t, socketPath, "OK reload users: users=2\n", received)

	out, err := runHaproxyRuntimeCommand(socketPath, "lua reload users\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "OK reload users") {
		t.Fatalf("unexpected output: %q", out)
	}

	select {
	case cmd := <-received:
		if cmd != "lua reload users\n" {
			t.Fatalf("unexpected command: %q", cmd)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive command on unix socket")
	}
}

func TestRunHaproxyRuntimeCommandErrResponse(t *testing.T) {
	tmpDir := t.TempDir()
	socketPath := filepath.Join(tmpDir, "haproxy.sock")
	received := make(chan string, 1)

	startUnixRuntimeServer(t, socketPath, "ERR reload users failed: parse error\n", received)

	_, err := runHaproxyRuntimeCommand(socketPath, "lua reload users\n")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "ERR reload users failed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunHaproxyRuntimeCommandSocketMissing(t *testing.T) {
	tmpDir := t.TempDir()
	socketPath := filepath.Join(tmpDir, "missing.sock")

	_, err := runHaproxyRuntimeCommand(socketPath, "lua reload users\n")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected os.ErrNotExist, got: %v", err)
	}
}
