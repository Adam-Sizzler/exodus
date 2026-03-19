package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"cerberus-node/config"
	"cerberus-node/constant"
	"cerberus-node/grpcserver"
	"cerberus-node/server"
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

	log.Printf("[START] cerberus-node application %s", constant.Version)

	if err := grpcserver.StartGRPCServer(&cfg, nodeServer); err != nil {
		cfg.Logger.Error("Failed to start gRPC server", "error", err)
		os.Exit(1)
	}
}
