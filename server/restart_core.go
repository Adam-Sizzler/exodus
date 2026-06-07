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
	cfg.Logger.Debug("Reloading sing-box via supervisorctl")

	if out, err := runSupervisorctl(cfg, 20*time.Second, "signal", "HUP", "singbox"); err == nil {
		cfg.Logger.Info("Core reloaded via supervisorctl signal HUP", "output", out)
		return nil
	} else {
		cfg.Logger.Warn("Core HUP reload failed, trying restart", "error", err, "output", out)
	}

	if out, err := runSupervisorctl(cfg, 45*time.Second, "restart", "singbox"); err == nil {
		cfg.Logger.Info("Core restarted via supervisorctl restart", "output", out)
		return nil
	} else {
		cfg.Logger.Warn("Core restart failed, trying start", "error", err, "output", out)
	}

	if out, err := runSupervisorctl(cfg, 30*time.Second, "start", "singbox"); err == nil {
		cfg.Logger.Info("Core started via supervisorctl start", "output", out)
		return nil
	} else {
		cfg.Logger.Error("Core start failed", "error", err, "output", out)
		return fmt.Errorf("reload singbox failed: %w: %s", err, out)
	}
}

func runSupervisorctl(cfg *config.NodeConfig, timeout time.Duration, command ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := supervisorctlArgs(command...)
	cfg.Logger.Trace("Executing supervisorctl", "args", strings.Join(args, " "))

	out, err := exec.CommandContext(ctx, "supervisorctl", args...).CombinedOutput()
	trimmed := strings.TrimSpace(string(out))
	if err != nil {
		return trimmed, fmt.Errorf("supervisorctl %s failed: %w", strings.Join(command, " "), err)
	}
	return trimmed, nil
}

func supervisorctlArgs(command ...string) []string {
	args := make([]string, 0, len(command)+6)
	if sock := strings.TrimSpace(os.Getenv("SUPERVISORD_SOCKET_PATH")); sock != "" {
		args = append(args, "-s", "unix://"+sock)
	}
	if user := strings.TrimSpace(os.Getenv("SUPERVISORD_USER")); user != "" {
		args = append(args, "-u", user)
	}
	if pass := strings.TrimSpace(os.Getenv("SUPERVISORD_PASSWORD")); pass != "" {
		args = append(args, "-p", pass)
	}
	args = append(args, command...)
	return args
}
