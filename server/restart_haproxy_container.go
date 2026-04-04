package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const (
	haproxyRuntimeSocketPath = "/var/run/haproxy/haproxy.sock"
	haproxyReloadUsersCmd    = "lua reload users\n"
	haproxySocketReadLimit   = 64 * 1024
)

type haproxyReloadResult struct {
	Reloaded bool
	Skipped  bool
	Warning  string
	Output   string
}

func reloadHaproxyUsers() haproxyReloadResult {
	result, err := runHaproxyRuntimeCommand(haproxyRuntimeSocketPath, haproxyReloadUsersCmd)
	if err == nil {
		return haproxyReloadResult{
			Reloaded: true,
			Output:   result,
		}
	}

	if errors.Is(err, os.ErrNotExist) {
		return haproxyReloadResult{
			Skipped: true,
			Warning: err.Error(),
		}
	}

	return haproxyReloadResult{
		Warning: err.Error(),
	}
}

func runHaproxyRuntimeCommand(socketPath string, command string) (string, error) {
	if strings.TrimSpace(socketPath) == "" {
		return "", fmt.Errorf("haproxy runtime socket path is empty")
	}
	if strings.TrimSpace(command) == "" {
		return "", fmt.Errorf("haproxy runtime command is empty")
	}

	if _, err := os.Stat(socketPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("haproxy runtime socket not found at %s: %w", socketPath, os.ErrNotExist)
		}
		return "", fmt.Errorf("stat haproxy runtime socket %s: %w", socketPath, err)
	}

	conn, err := net.DialTimeout("unix", socketPath, 3*time.Second)
	if err != nil {
		return "", fmt.Errorf("connect haproxy runtime socket %s: %w", socketPath, err)
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))

	if _, err := io.WriteString(conn, command); err != nil {
		return "", fmt.Errorf("send command to haproxy runtime socket: %w", err)
	}

	raw, err := io.ReadAll(io.LimitReader(conn, haproxySocketReadLimit))
	if err != nil {
		return "", fmt.Errorf("read response from haproxy runtime socket: %w", err)
	}

	response := strings.TrimSpace(string(raw))
	if response == "" {
		return "", fmt.Errorf("empty response from haproxy runtime command")
	}

	firstLine := response
	if idx := strings.IndexByte(response, '\n'); idx >= 0 {
		firstLine = response[:idx]
	}
	if strings.HasPrefix(firstLine, "ERR") {
		return "", fmt.Errorf("haproxy runtime command failed: %s", firstLine)
	}
	if !strings.HasPrefix(firstLine, "OK") {
		return "", fmt.Errorf("unexpected haproxy runtime response: %s", firstLine)
	}

	return response, nil
}
