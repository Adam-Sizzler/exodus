package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/exodus/subscription-page/backend/internal/config"
	"github.com/exodus/subscription-page/backend/internal/grpcapi"
	"github.com/exodus/subscription-page/backend/internal/grpcserver"
	"github.com/exodus/subscription-page/backend/internal/server"
)

var buildVersion = "unknown"

var semverPattern = regexp.MustCompile(`^[vV]?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z\.-]+)?$`)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] %v", err)
	}

	log.Printf("[CONFIG] SUB_PATH: %s", config.DisplayPrefix(cfg.SubPath))
	log.Printf("[CONFIG] grpc target %s:%d", cfg.GRPCAddress, cfg.GRPCPort)
	log.Printf("[CONFIG] http listening on :%s", cfg.AppPort)

	nodeService := grpcapi.NewNodeServer(resolveNodeVersion(cfg.AppVersion))

	application, appErr := server.New(cfg, nodeService)
	if appErr != nil {
		log.Fatalf("[FATAL] %v", appErr)
	}

	httpServer := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           application,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	errCh := make(chan error, 2)

	go func() {
		if err := grpcserver.Start(ctx, cfg, nodeService); err != nil {
			errCh <- err
		}
	}()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("[INFO] shutdown signal received")
	case err = <-errCh:
		log.Fatalf("[FATAL] server failed: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("[WARN] http shutdown failed: %v", err)
	}
}

func resolveNodeVersion(configVersion string) string {
	resolved := strings.TrimSpace(configVersion)
	if resolved == "" {
		resolved = strings.TrimSpace(buildVersion)
	}
	if resolved == "" {
		return "v0.0.0-dev"
	}

	lower := strings.ToLower(resolved)
	switch lower {
	case "unknown", "latest", "(devel)":
		return "v0.0.0-dev"
	}

	if semverPattern.MatchString(resolved) {
		if strings.HasPrefix(lower, "v") {
			return "v" + strings.TrimSpace(resolved[1:])
		}
		return "v" + resolved
	}

	return resolved
}
