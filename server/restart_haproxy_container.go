package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	haproxyContainerName = "haproxy"
	dockerSocketPath     = "/var/run/docker.sock"
)

type haproxyRestartResult struct {
	Reloaded  bool
	Restarted bool
	Warning   string
}

func restartHaproxyContainer() haproxyRestartResult {
	container := haproxyContainerName
	socketPath := dockerSocketPath

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{Timeout: 10 * time.Second}
				return d.DialContext(ctx, "unix", socketPath)
			},
		},
	}

	// Prefer soft reload first: send SIGHUP to the running HAProxy process.
	// If it fails (container missing/stopped), fallback to full restart.
	reloadEndpoint := fmt.Sprintf("http://docker/containers/%s/kill?signal=HUP", url.PathEscape(container))
	reloadErr := doDockerPost(client, reloadEndpoint)
	if reloadErr == nil {
		return haproxyRestartResult{Reloaded: true}
	}

	restartEndpoint := fmt.Sprintf("http://docker/containers/%s/restart?t=10", url.PathEscape(container))
	restartErr := doDockerPost(client, restartEndpoint)
	if restartErr == nil {
		return haproxyRestartResult{
			Restarted: true,
			Warning:   fmt.Sprintf("soft reload failed, fallback to restart: %v", reloadErr),
		}
	}

	return haproxyRestartResult{
		Warning: fmt.Sprintf("skip HAProxy reload/restart: reload_err=%v restart_err=%v", reloadErr, restartErr),
	}
}

func doDockerPost(client *http.Client, endpoint string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build docker request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("call docker api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("docker api returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
