package server

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"exodus-node/config"
)

func (s *NodeServer) diagnoseCoreFailure(ctx context.Context, baseErr error) string {
	message := fmt.Sprintf("core stats unavailable: %v", baseErr)

	checkCtx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	output, err := runSingboxCheck(checkCtx)
	if err == nil {
		trimmed := strings.TrimSpace(output)
		if trimmed != "" {
			message = fmt.Sprintf("%s; sing-box check: %s", message, trimmed)
		}
		s.Cfg.LoggerFor("SingboxService").Error("Core health-check failed", "diagnostic", message)
		return message
	}

	trimmedErr := strings.TrimSpace(err.Error())
	if trimmedErr != "" {
		message = fmt.Sprintf("%s; sing-box check failed: %s", message, trimmedErr)
	}
	s.Cfg.LoggerFor("SingboxService").Error("Core health-check failed", "diagnostic", message)
	return message
}

func runSingboxCheck(ctx context.Context) (string, error) {
	commands := [][]string{
		{"/usr/local/bin/sing-box", "check", "-c", config.FixedSingboxConfigPath},
		{"sing-box", "check", "-c", config.FixedSingboxConfigPath},
	}

	var lastErr error
	for _, cmdArgs := range commands {
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		out, err := cmd.CombinedOutput()
		if err == nil {
			return string(out), nil
		}
		lastErr = fmt.Errorf("%s: %w: %s", strings.Join(cmdArgs, " "), err, strings.TrimSpace(string(out)))
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("sing-box check command not found")
	}
	return "", lastErr
}
