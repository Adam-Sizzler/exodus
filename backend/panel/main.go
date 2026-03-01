package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"syscall"
	"time"

	"v2ray-stat/backend/panel/api"
	"v2ray-stat/backend/panel/config"
	"v2ray-stat/backend/panel/db"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/redisqueue"
	"v2ray-stat/backend/panel/users"
	"v2ray-stat/constant"
)

// startAPIServer starts the API server.
func startAPIServer(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", cfg.V2RS.Address, cfg.V2RS.Port)
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: api.WithCORS(cfg, api.WithPanelAuth(manager, cfg, mux)),
	}

	mux.HandleFunc("/api/v1/auth/bootstrap", api.AuthBootstrapHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/setup", api.AuthSetupHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/login", api.AuthLoginHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/logout", api.AuthLogoutHandler(manager, cfg))
	mux.HandleFunc("/api/v1/auth/me", api.AuthMeHandler(manager, cfg))

	mux.HandleFunc("/api/v1/panel/settings", api.PanelSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/panel/api-tokens", api.PanelAPITokensHandler(manager, cfg))
	mux.HandleFunc("/api/v1/panel/api-tokens/", api.PanelAPITokenByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/v1/nodes", api.NodesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/summary", api.NodesSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/reorder", api.NodesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes/", api.NodeByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/reorder", api.ConfigProfilesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/reorder", api.InternalSquadsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hosts/reorder", api.HostsReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hosts", api.HostsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hosts/", api.HostByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/hosts-to-nodes", api.HostNodeAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list", api.UsersAPIHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/reorder", api.UsersReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/summary", api.UsersListSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/", api.UserByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/users-list/create", api.UsersCreateHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles", api.ConfigProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles/", api.ConfigProfileByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads", api.InternalSquadsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/internal-squads/", api.InternalSquadByUUIDHandler(manager, cfg))

	// Squad bindings and queries
	mux.HandleFunc("/api/v1/inbound-assignments", api.InboundAssignmentsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-inbounds", api.SquadInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-members", api.SquadMembersHandler(manager, cfg))
	mux.HandleFunc("/api/v1/config-profiles-with-inbounds", api.ConfigProfilesWithInboundsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squad-details/", api.SquadDetailsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/squads-summary", api.AllSquadsSummaryHandler(manager, cfg))
	mux.HandleFunc("/api/v1/nodes-with-config", api.NodesWithConfigHandler(manager, cfg))
	mux.HandleFunc("/api/v1/inbounds-with-profiles", api.InboundsWithProfilesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/subscription-settings", api.SubscriptionSettingsHandler(manager, cfg))
	mux.HandleFunc("/api/v1/subscription-settings/", api.SubscriptionSettingsByUUIDHandler(manager, cfg))
	mux.HandleFunc("/api/v1/templates/reorder", api.SubscriptionTemplatesReorderHandler(manager, cfg))
	mux.HandleFunc("/api/v1/templates", api.SubscriptionTemplatesHandler(manager, cfg))
	mux.HandleFunc("/api/v1/templates/", api.SubscriptionTemplateByUUIDHandler(manager, cfg))

	mux.HandleFunc("/api/health", api.HealthHandler())
	mux.HandleFunc("/api/v1/health", api.HealthHandler())

	mux.Handle("/api/", http.NotFoundHandler())

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

// startWebServer serves the panel UI and proxies API requests to the backend.
func startWebServer(ctx context.Context, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", cfg.V2RS.Address, cfg.Panel.WebPort)
	targetURL, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", cfg.V2RS.Port))
	if err != nil {
		cfg.Logger.Fatal("Failed to parse backend URL", "error", err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		if req.Header.Get("X-Forwarded-Proto") == "" {
			if req.TLS != nil {
				req.Header.Set("X-Forwarded-Proto", "https")
			} else {
				req.Header.Set("X-Forwarded-Proto", "http")
			}
		}
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		cfg.Logger.Warn("API proxy error", "error", err)
		http.Error(w, "backend unavailable", http.StatusBadGateway)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	uiDir := cfg.Panel.StaticDir
	indexPath := filepath.Join(uiDir, "index.html")
	staticFS := http.FileServer(http.Dir(uiDir))
	if _, err := os.Stat(indexPath); err != nil {
		cfg.Logger.Warn("Panel UI index not found; static UI disabled", "path", indexPath, "error", err)
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		if _, err := os.Stat(indexPath); err != nil {
			http.NotFound(w, r)
			return
		}

		requestPath := filepath.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		targetPath := filepath.Join(uiDir, requestPath)
		if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
			staticFS.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, indexPath)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: api.WithCORS(cfg, mux),
	}

	cfg.Logger.Debug("Starting web server", "address", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			cfg.Logger.Fatal("Failed to start web server", "error", err)
		}
	}()

	<-ctx.Done()

	cfg.Logger.Debug("Shutting down web server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		cfg.Logger.Error("Error shutting down web server", "error", err)
	}
}

func main() {
	var versionFlag = flag.Bool("version", false, "Show version information")
	var resetAdminPasswordFlag = flag.Bool("reset-admin-password", false, "Interactively reset admin credentials and exit")
	flag.Parse()

	if *versionFlag {
		fmt.Println(constant.GetBuildInfo())
		os.Exit(0)
	}

	cfg, err := config.LoadConfig("config.yml")
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

	redisWorker, err := redisqueue.NewWorker(&cfg, manager)
	if err != nil {
		cfg.Logger.Warn("Redis worker disabled", "error", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Prepare wg
	var wg sync.WaitGroup
	wg.Add(3)

	go startAPIServer(ctx, manager, &cfg, &wg)
	go startWebServer(ctx, &cfg, &wg)
	go nodeMonitor.Start(ctx, &wg)
	if redisWorker != nil {
		redisWorker.Start(ctx, &wg)
	}

	log.Printf("[START] v2rs-panel application %s", constant.Version)

	<-sigChan
	cfg.Logger.Info("Received termination signal, saving data")
	cancel()

	// Stop node monitor
	nodeMonitor.Stop()

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
