package main

import (
	"flag"
	"fmt"
	"os"

	"exodus-node/config"
	"exodus-node/constant"
	"exodus-node/grpcserver"
	"exodus-node/server"
)

func main() {
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
		cfg.Logger.Error("Failed to create node server", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := nodeServer.Close(); err != nil {
			cfg.Logger.Warn("Failed to close node server", "error", err)
		}
	}()

	cfg.Logger.Info("Starting exodus-node", "version", constant.Version)

	if err := grpcserver.StartGRPCServer(&cfg, nodeServer); err != nil {
		cfg.Logger.Error("Failed to start gRPC server", "error", err)
		os.Exit(1)
	}
}
