package server

import (
	"context"
	"fmt"
	"os"
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

const (
	DefaultSingboxLogPath = "/var/log/singbox/current"
)

var (
	tai64nRegexp   = regexp.MustCompile(`^@[0-9a-fA-F]{16,24}\s*`)
	fatalPfxRegexp = regexp.MustCompile(`^(?:FATAL|ERROR)(?:\[[0-9]+\])?\s*`)
)

// TailSingboxLogLines reads the last N lines from the Sing-box log file.
func TailSingboxLogLines(logPath string, n int) []string {
	if logPath == "" {
		logPath = DefaultSingboxLogPath
	}
	if n <= 0 {
		n = 10
	}

	// Try tail command first
	out, err := exec.Command("tail", "-n", fmt.Sprintf("%d", n), logPath).Output()
	var rawLines []string
	if err == nil {
		rawLines = strings.Split(string(out), "\n")
	} else {
		// Fallback to reading file directly
		content, readErr := os.ReadFile(logPath)
		if readErr != nil {
			return nil
		}
		allLines := strings.Split(string(content), "\n")
		if len(allLines) > n {
			rawLines = allLines[len(allLines)-n:]
		} else {
			rawLines = allLines
		}
	}

	cleaned := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = stripANSI(line)
		line = tai64nRegexp.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return cleaned
}

// ExtractSingboxLogReason scans recent log lines to find the root-cause error reason.
func ExtractSingboxLogReason(logPath string, maxLines int) string {
	lines := TailSingboxLogLines(logPath, maxLines)
	if len(lines) == 0 {
		return ""
	}

	// Scan in reverse (newest to oldest) looking for fatal, error, or panic lines
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		upper := strings.ToUpper(line)
		if strings.Contains(upper, "FATAL") ||
			strings.Contains(upper, "ERROR") ||
			strings.Contains(upper, "PANIC:") ||
			strings.Contains(upper, "BIND: ADDRESS ALREADY IN USE") ||
			strings.Contains(upper, "CREATE SERVICE:") ||
			strings.Contains(upper, "START SERVICE") {
			cleaned := fatalPfxRegexp.ReplaceAllString(line, "")
			cleaned = strings.TrimSpace(cleaned)
			if len(cleaned) > 500 {
				cleaned = cleaned[:500]
			}
			return cleaned
		}
	}

	// Fallback to the very last non-empty line
	last := lines[len(lines)-1]
	last = fatalPfxRegexp.ReplaceAllString(last, "")
	last = strings.TrimSpace(last)
	if len(last) > 500 {
		last = last[:500]
	}
	return last
}

