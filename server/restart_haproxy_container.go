package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	haproxyContainerName = "haproxy"
	dockerSocketPath     = "/var/run/docker.sock"
)

func restartHaproxyContainer() error {
	container := haproxyContainerName
	socketPath := dockerSocketPath

	endpoint := fmt.Sprintf("http://docker/containers/%s/restart?t=10", url.PathEscape(container))
	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				d := net.Dialer{Timeout: 10 * time.Second}
				return d.DialContext(ctx, "unix", socketPath)
			},
		},
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build docker restart request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("call docker restart api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotModified {
		return fmt.Errorf("docker restart api returned status %d", resp.StatusCode)
	}

	return nil
}
