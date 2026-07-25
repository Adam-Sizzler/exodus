package exodus

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/health"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/system"
	"exodus/internal/logger"
)

func startMetricsServer(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	if cfg == nil || cfg.Metrics.Port <= 0 {
		return
	}

	metricsAddress := strings.TrimSpace(cfg.Metrics.Address)
	if metricsAddress == "" {
		metricsAddress = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", metricsAddress, cfg.Metrics.Port)
	metricsHandler := system.MetricsHandler(manager, cfg)
	healthHandler := health.HealthHandler()

	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/health", healthHandler)
	registered := map[string]struct{}{
		"/metrics":    {},
		"/health":     {},
		"/api/health": {},
	}
	for _, path := range []string{cfg.Panel.BasePath, strings.TrimSuffix(cfg.Panel.BasePath, "/")} {
		normalized := strings.TrimSpace(path)
		if normalized == "" || normalized == "/" {
			continue
		}
		endpointMetrics := strings.TrimSuffix(normalized, "/") + "/metrics"
		if _, exists := registered[endpointMetrics]; !exists {
			registered[endpointMetrics] = struct{}{}
			mux.Handle(endpointMetrics, metricsHandler)
		}
		endpointHealth := strings.TrimSuffix(normalized, "/") + "/health"
		if _, exists := registered[endpointHealth]; !exists {
			registered[endpointHealth] = struct{}{}
			mux.HandleFunc(endpointHealth, healthHandler)
		}
	}

	server := &http.Server{
		Addr:    addr,
		Handler: middleware.WithRequestLogging(cfg, "metrics", mux),
	}

	cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceMetrics).Info("Metrics reporter started", "address", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceMetrics).Error("Failed to start metrics server", "error", err)
		}
	}()

	<-ctx.Done()

	cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceMetrics).Debug("Shutting down metrics server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceMetrics).Error("Error shutting down metrics server", "error", err)
	}
}
