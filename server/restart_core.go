package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"cerberus-node/config"
)

func restartCoreProcess(cfg *config.NodeConfig) error {
	cfg.Logger.Debug("Restarting sing-box via supervisorctl")
	return restartViaSupervisor(cfg)
}

func restartViaSupervisor(cfg *config.NodeConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	args := make([]string, 0, 10)
	if sock := strings.TrimSpace(os.Getenv("SUPERVISORD_SOCKET_PATH")); sock != "" {
		args = append(args, "-s", "unix://"+sock)
	}
	if user := strings.TrimSpace(os.Getenv("SUPERVISORD_USER")); user != "" {
		args = append(args, "-u", user)
	}
	if pass := strings.TrimSpace(os.Getenv("SUPERVISORD_PASSWORD")); pass != "" {
		args = append(args, "-p", pass)
	}
	args = append(args, "restart", "singbox")

	out, err := exec.CommandContext(ctx, "supervisorctl", args...).CombinedOutput()
	if err != nil {
		cfg.Logger.Error("supervisorctl restart failed", "error", err, "output", strings.TrimSpace(string(out)))
		return fmt.Errorf("restart singbox via supervisorctl failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	cfg.Logger.Info("Core restarted via supervisorctl", "output", strings.TrimSpace(string(out)))
	return nil
}
