package geocheck

import (
	"context"
	"strings"

	"exodus-node/geocheck/internal/cli"
)

// Run executes geocheck in-process with the given options and returns the raw JSON payload with base64 embedded SVG.
func Run(ctx context.Context, iface string, ip string) ([]byte, error) {
	args := []string{"--json", "--svg-base64", "--quiet"}
	if trimmed := strings.TrimSpace(iface); trimmed != "" {
		args = append(args, "--interface", trimmed)
	} else if trimmedIP := strings.TrimSpace(ip); trimmedIP != "" {
		args = append(args, "--interface", trimmedIP)
	}
	return cli.RunJSONBytes(ctx, args)
}

// RunCLI executes geocheck CLI directly with command-line arguments and returns the exit code.
func RunCLI(args []string) int {
	return cli.Run(args)
}
