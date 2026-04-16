package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"exodus/backend/config"
	"exodus/backend/db"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/httpapi"
	"exodus/backend/httpapi/middleware"
	"exodus/backend/httpapi/system"
	users "exodus/backend/nodes"
	"exodus/backend/redisqueue"
	"exodus/backend/srslists"
	"exodus/backend/subscriptionnodes"
	"exodus/constant"
)

// startWebServer serves both panel UI and API on a single APP_PORT.
func startWebServer(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", cfg.EXODUS.Address, cfg.Panel.AppPort)
	panelBasePath := cfg.Panel.BasePath
	panelBasePathNoTrailing := strings.TrimSuffix(panelBasePath, "/")
	if panelBasePathNoTrailing == "" {
		panelBasePathNoTrailing = "/"
	}

	apiHandler := httpapi.NewAPIHandler(manager, cfg)
	metricsHandler := system.MetricsHandler(manager, cfg)

	mux := http.NewServeMux()
	uiDir := cfg.Panel.StaticDir
	indexPath := filepath.Join(uiDir, "index.html")
	staticFS := http.FileServer(http.Dir(uiDir))
	if _, err := os.Stat(indexPath); err != nil {
		cfg.Logger.Warn("Panel UI index not found; static UI disabled", "path", indexPath, "error", err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestPath := r.URL.Path

		if panelBasePath != "/" && requestPath == panelBasePathNoTrailing {
			http.Redirect(w, r, panelBasePath, http.StatusPermanentRedirect)
			return
		}

		if !strings.HasPrefix(requestPath, panelBasePath) {
			http.NotFound(w, r)
			return
		}

		relativePath := strings.TrimPrefix(requestPath, panelBasePath)
		if relativePath == "metrics" {
			metricsHandler.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(relativePath, "api/") {
			apiReq := r.Clone(r.Context())
			apiReq.URL.Path = "/" + relativePath
			apiReq.URL.RawPath = ""
			apiHandler.ServeHTTP(w, apiReq)
			return
		}

		if relativePath == "app-config.js" {
			serveAppConfigJS(w, panelBasePathNoTrailing)
			return
		}

		cleanPath := filepath.Clean(relativePath)
		if strings.HasPrefix(cleanPath, "..") {
			http.NotFound(w, r)
			return
		}

		if cleanPath != "." && cleanPath != "" {
			targetPath := filepath.Join(uiDir, cleanPath)
			if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
				staticReq := r.Clone(r.Context())
				staticReq.URL.Path = "/" + cleanPath
				staticReq.URL.RawPath = ""
				staticFS.ServeHTTP(w, staticReq)
				return
			}
		}

		if _, err := os.Stat(indexPath); err != nil {
			http.NotFound(w, r)
			return
		}

		servePanelIndex(w, indexPath, panelBasePath, panelBasePathNoTrailing)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: middleware.WithCORS(cfg, middleware.WithRequestLogging(cfg, "web", mux)),
	}

	cfg.Logger.Debug("Starting web/API server", "address", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Fatal("Failed to start web server", "error", err)
		}
	}()

	<-ctx.Done()

	cfg.Logger.Debug("Shutting down web/API server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.Error("Error shutting down web server", "error", err)
	}
}

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

func serveAppConfigJS(w http.ResponseWriter, basePath string) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = fmt.Fprintf(
		w,
		"window.__EXODUS_RUNTIME__={basePath:%q};\n",
		basePath,
	)
}

func servePanelIndex(w http.ResponseWriter, indexPath string, basePathWithSlash, basePath string) {
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		http.Error(w, "panel index not found", http.StatusNotFound)
		return
	}

	injected := fmt.Sprintf(
		"<base href=\"%s\" />\n<script>window.__EXODUS_RUNTIME__={basePath:%q};</script>",
		html.EscapeString(basePathWithSlash),
		basePath,
	)

	page := string(indexBytes)
	// Normalize accidental ".//" asset prefixes from the built index template.
	page = strings.ReplaceAll(page, ".//", "./")
	if strings.Contains(page, "<head>") {
		page = strings.Replace(page, "<head>", "<head>\n"+injected, 1)
	} else {
		page = injected + page
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page))
}

func main() {
	var versionFlag = flag.Bool("version", false, "Show version information")
	var resetAdminPasswordFlag = flag.Bool("reset-admin-password", false, "Interactively reset admin credentials and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(constant.GetBuildInfo())
		os.Exit(0)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	if *resetAdminPasswordFlag {
		if err := runAdminCredentialReset(&cfg); err != nil {
			cfg.Logger.Error("Failed to reset admin credentials", "error", err)
			os.Exit(1)
		}
		return
	}

	fmt.Println("Hint: run backend with --reset-admin-password for emergency admin credential reset.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := db.InitDatabase(&cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to initialize database", "error", err)
	}

	manager, err := dbmanager.NewDatabaseManager(dbConn, ctx, 1, 300, 500, &cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to create DatabaseManager", "error", err)
	}

	// Create and start node monitor (dynamically manages nodes from DB)
	nodeMonitor := users.NewNodeMonitor(manager, &cfg)
	users.RegisterGlobalNodeMonitor(nodeMonitor)
	subNodeMonitor := subscriptionnodes.NewSubNodeMonitor(manager, &cfg)
	subscriptionnodes.RegisterGlobalSubNodeMonitor(subNodeMonitor)

	redisWorker, err := redisqueue.NewWorker(&cfg, manager)
	if err != nil {
		cfg.Logger.Warn("Redis worker disabled", "error", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Prepare wg
	var wg sync.WaitGroup
	wg.Add(4)

	go startWebServer(ctx, manager, &cfg, &wg)
	go startMetricsServer(ctx, manager, &cfg, &wg)
	go nodeMonitor.Start(ctx, &wg)
	go subNodeMonitor.Start(ctx, &wg)
	srslists.StartPeriodicChecker(ctx, &wg, manager, &cfg, 5*time.Minute)
	if redisWorker != nil {
		redisWorker.Start(ctx, &wg)
	}

	log.Printf("[START] exodus application %s", constant.Version)

	<-sigChan
	cfg.Logger.Info("Received termination signal, saving data")
	cancel()

	// Stop node monitor
	nodeMonitor.Stop()
	subNodeMonitor.Stop()

	// Закрываем DatabaseManager
	manager.Close()

	// Ждём завершения горутин
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
		// Профилирование горутин при тайм-ауте
		pprof.Lookup("goroutine").WriteTo(os.Stderr, 1)
	}

	// Закрываем базу данных
	if err := dbConn.Close(); err != nil {
		cfg.Logger.Error("Failed to close database", "error", err)
	}
	if redisWorker != nil {
		if err := redisWorker.Close(); err != nil {
			cfg.Logger.Warn("Failed to close redis worker", "error", err)
		}
	}

	log.Printf("[STOP] Program terminated")
}
