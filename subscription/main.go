package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"v2ray-stat/backend/api"
	"v2ray-stat/common"
	"v2ray-stat/constant"
	"v2ray-stat/subscription/config"
	"v2ray-stat/subscription/grpcserver"
	"v2ray-stat/subscription/handler"
	"v2ray-stat/subscription/templates"
)

func LoggingMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		cfg.Logger.Debug("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
			"remote", r.RemoteAddr,
		)
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")

		// Если куки нет, редиректим на Google
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "https://www.google.com", http.StatusFound)
			return
		}

		// Предотвращаем кэширование, чтобы проверка куки срабатывала всегда
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		next.ServeHTTP(w, r)
	})
}

// startAPIServer starts the API server.
func startAPIServer(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := fmt.Sprintf("%s:%d", cfg.V2RSSub.Address, cfg.V2RSSub.Port)

	mux := http.NewServeMux()

	staticFS := http.FileServer(http.Dir("/app/assets"))
	mux.Handle("/api/v1/assets/", http.StripPrefix("/api/v1/assets/", AuthMiddleware(staticFS)))

	mux.HandleFunc("/api/v1/sub", handler.SubscriptionHandler)

	// Оборачиваем все в логгер и стандартные заголовки
	finalHandler := LoggingMiddleware(cfg, api.WithServerHeader(mux))

	server := &http.Server{
		Addr:    addr,
		Handler: finalHandler,
	}

	cfg.Logger.Debug("Starting API server", "address", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Fatal("Failed to start server", "error", err)
		}
	}()

	<-ctx.Done()

	cfg.Logger.Debug("Shutting down API server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.Error("Error shutting down server", "error", err)
	}
}

func main() {
	var versionFlag = flag.Bool("version", false, "Show version information")
	flag.Parse()
	if *versionFlag {
		fmt.Println(constant.GetBuildInfo())
		os.Exit(0)
	}

	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	if cfg.Logger == nil {
		log.Fatalf("cfg.Logger is nil after LoadConfig")
	}

	cfg.Logger.Debug("Global config initialized", "config", fmt.Sprintf("%+v", cfg))

	common.InitTimezone(cfg.TZ, cfg.Logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initial load of templates
	if err := templates.LoadTemplates(&cfg); err != nil {
		cfg.Logger.Fatal("Failed to load templates", "error", err)
	} else {
		cfg.Logger.Info("Templates loaded successfully")
	}

	var wg sync.WaitGroup
	wg.Add(4)

	go grpcserver.StartGrpcServer(ctx, &cfg, &wg)
	go startAPIServer(ctx, &cfg, &wg)
	go config.WatchConfig(ctx, &cfg, &wg)
	go templates.WatchTemplates(ctx, &cfg, &wg)

	log.Printf("[START] v2ray-stat-subscription application %s", constant.Version)

	<-sigChan
	cfg.Logger.Info("Received termination signal")
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		cfg.Logger.Debug("All goroutines completed")
	case <-time.After(10 * time.Second):
		cfg.Logger.Warn("Timeout waiting for goroutines to complete, forcing shutdown")
		pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	}

	cfg.Logger.Info("Program terminated")
}
