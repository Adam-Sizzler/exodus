package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/exodus/subscription-page/backend/internal/config"
	"github.com/exodus/subscription-page/backend/internal/grpcapi"
	"github.com/exodus/subscription-page/backend/internal/grpcserver"
	"github.com/exodus/subscription-page/backend/internal/logger"
	"github.com/exodus/subscription-page/backend/internal/server"
)

var buildVersion = "unknown"

var semverPattern = regexp.MustCompile(`^[vV]?\d+\.\d+\.\d+(?:[-+][0-9A-Za-z\.-]+)?$`)

func main() {
	logger.Configure(logger.Options{
		NodeEnv:   os.Getenv("NODE_ENV"),
		DebugLogs: os.Getenv("ENABLE_DEBUG_LOGS"),
		Colors:    true,
	})
	bootstrapLogger := logger.WithContext("Bootstrap")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		bootstrapLogger.Fatal("[FATAL] configuration failed", err)
	}

	if cfg.SubPath == "" {
		bootstrapLogger.Log("[CONFIG] CUSTOM_SUB_PREFIX: not set")
	} else {
		bootstrapLogger.Log("[CONFIG] CUSTOM_SUB_PREFIX: " + cfg.SubPath)
	}

	resolvedVersion := resolveNodeVersion(cfg.AppVersion)
	nodeService := grpcapi.NewNodeServer(resolvedVersion)

	application, appErr := server.New(cfg, nodeService)
	if appErr != nil {
		bootstrapLogger.Fatal("[FATAL] application bootstrap failed", appErr)
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

	listener, err := net.Listen("tcp", httpServer.Addr)
	if err != nil {
		bootstrapLogger.Fatal("[FATAL] listen http", err)
	}
	go func() {
		if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	bootstrapLogger.Log("\n" + server.GetStartMessage(resolvedVersion) + "\n")

	select {
	case <-ctx.Done():
		bootstrapLogger.Log("shutdown signal received")
	case err = <-errCh:
		bootstrapLogger.Fatal("[FATAL] server failed", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		bootstrapLogger.Warn("http shutdown failed", logger.String("error", err.Error()))
	}
}

func resolveNodeVersion(configVersion string) string {
	if resolved := normalizeNodeVersion(configVersion); resolved != "" {
		return resolved
	}
	if resolved := normalizeNodeVersion(buildVersion); resolved != "" {
		return resolved
	}

	return "v0.0.0-dev"
}

func normalizeNodeVersion(raw string) string {
	resolved := strings.TrimSpace(raw)
	if resolved == "" {
		return ""
	}

	lower := strings.ToLower(resolved)
	switch lower {
	case "unknown", "latest", "(devel)":
		return ""
	}

	if semverPattern.MatchString(resolved) {
		if strings.HasPrefix(lower, "v") {
			return "v" + strings.TrimSpace(resolved[1:])
		}
		return "v" + resolved
	}

	return resolved
}
