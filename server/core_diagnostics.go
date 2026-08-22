package server

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"exodus-node/config"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

func stripANSI(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

// RunSingboxCheck executes `sing-box check -c <configPath>` and returns the cleaned command output.
func RunSingboxCheck(ctx context.Context, configPath string) (string, error) {
	if configPath == "" {
		configPath = config.FixedSingboxConfigPath
	}

	commands := [][]string{
		{"/usr/local/bin/sing-box", "check", "-c", configPath},
		{"sing-box", "check", "-c", configPath},
	}

	var lastErr error
	for _, cmdArgs := range commands {
		cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
		out, err := cmd.CombinedOutput()
		outputStr := strings.TrimSpace(stripANSI(string(out)))
		if err == nil {
			return outputStr, nil
		}
		if outputStr != "" {
			return outputStr, fmt.Errorf("%s", outputStr)
		}
		lastErr = fmt.Errorf("%s: %w", strings.Join(cmdArgs, " "), err)
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("sing-box check command not found")
	}
	return "", lastErr
}
