package common

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"v2ray-stat/logger"
	"v2ray-stat/backend/node/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func RestartService(serviceName string, cfg *config.NodeConfig) error {
	mgr, err := NewServiceManager(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return mgr.Restart(ctx, serviceName)
}

type ServiceManager interface {
	Restart(ctx context.Context, serviceName string) error
}

func NewServiceManager(cfg *config.NodeConfig) (ServiceManager, error) {
	switch cfg.ServiceManager {
	case "docker":
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			return nil, err
		}
		return &DockerManager{client: cli, logger: cfg.Logger}, nil
	default:
		return &SystemdManager{logger: cfg.Logger}, nil
	}
}

type DockerManager struct {
	client *client.Client
	logger *logger.Logger
}

func (d *DockerManager) Restart(ctx context.Context, serviceName string) error {
	containerName := serviceName
	switch serviceName {
	case "sing-box", "singbox":
		containerName = "singbox"
	case "xray":
		containerName = "xray"
	}

	d.logger.Info("Restarting Docker service", "container", containerName)

	timeout := 10
	err := d.client.ContainerRestart(ctx, containerName, container.StopOptions{Timeout: &timeout})
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)
	return nil
}

type SystemdManager struct {
	logger *logger.Logger
}

func (s *SystemdManager) Restart(ctx context.Context, serviceName string) error {
	s.logger.Info("Restarting Systemd service", "service", serviceName)
	cmd := exec.CommandContext(ctx, "systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("systemctl failed: %s, output: %s", err, string(output))
	}
	return nil
}
