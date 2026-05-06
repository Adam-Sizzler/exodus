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

	cfg.Logger.Debug("Starting metrics server", "address", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Error("Failed to start metrics server", "error", err)
		}
	}()

	<-ctx.Done()

	cfg.Logger.Debug("Shutting down metrics server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.Error("Error shutting down metrics server", "error", err)
	}
}
