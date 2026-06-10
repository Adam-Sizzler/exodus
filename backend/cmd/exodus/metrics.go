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

	mux := http.NewServeMux()
	mux.Handle("/metrics", metricsHandler)
	registered := map[string]struct{}{
		"/metrics": {},
	}
	for _, path := range []string{cfg.Panel.BasePath, strings.TrimSuffix(cfg.Panel.BasePath, "/")} {
		normalized := strings.TrimSpace(path)
		if normalized == "" || normalized == "/" {
			continue
		}
		endpoint := strings.TrimSuffix(normalized, "/") + "/metrics"
		if _, exists := registered[endpoint]; exists {
			continue
		}
		registered[endpoint] = struct{}{}
		mux.Handle(endpoint, metricsHandler)
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
