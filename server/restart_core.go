package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"exodus-node/config"
)

func reloadCoreProcess(cfg *config.NodeConfig) error {
	cfg.Logger.Debug("Reloading sing-box via supervisorctl signal HUP")
	return reloadViaSupervisorHUP(cfg)
}

func reloadViaSupervisorHUP(cfg *config.NodeConfig) error {
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
	args = append(args, "signal", "HUP", "singbox")

	out, err := exec.CommandContext(ctx, "supervisorctl", args...).CombinedOutput()
	if err != nil {
		cfg.Logger.Error("supervisorctl signal HUP failed", "error", err, "output", strings.TrimSpace(string(out)))
		return fmt.Errorf("reload singbox via supervisorctl signal HUP failed: %w: %s", err, strings.TrimSpace(string(out)))
	}
	cfg.Logger.Info("Core reloaded via supervisorctl signal HUP", "output", strings.TrimSpace(string(out)))
	return nil
}
