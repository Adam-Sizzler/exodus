package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"exodus-node/config"
	"exodus-node/constant"
	"exodus-node/geocheck"
	"exodus-node/grpcserver"
	"exodus-node/server"
)

func main() {
	execName := filepath.Base(os.Args[0])

	// 1. Direct invocation via symlink "cli" or "geocheck" (e.g. `cli -j ...` or `geocheck --demo`):
	if execName == "cli" || execName == "geocheck" {
		os.Exit(geocheck.RunCLI(os.Args[1:]))
	}

	// 2. Subcommand invocation via `exodus-node cli ...` or `exodus-node geocheck ...`:
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cli", "geocheck", "check":
			os.Exit(geocheck.RunCLI(os.Args[2:]))
		}
	}

	// 3. Normal node daemon bootstrap:
	var versionFlag = flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *versionFlag {
		fmt.Println(constant.GetBuildInfo())
		os.Exit(0)
	}

	cfg, err := config.LoadNodeConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load node config: %v\n", err)
		os.Exit(1)
	}

	nodeServer, err := server.NewNodeServer(&cfg)
	if err != nil {
		cfg.LoggerFor("Bootstrap").Error("Failed to create node server", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := nodeServer.Close(); err != nil {
			cfg.LoggerFor("Bootstrap").Warn("Failed to close node server", "error", err)
		}
	}()

	cfg.Logger.Log(server.GetStartMessage(&cfg))

	if err := grpcserver.StartGRPCServer(&cfg, nodeServer); err != nil {
		cfg.LoggerFor("Bootstrap").Error("Failed to start gRPC server", "error", err)
		os.Exit(1)
	}
}
