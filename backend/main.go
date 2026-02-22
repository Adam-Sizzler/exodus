package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"v2ray-stat/backend/api"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/grpcclient"
	"v2ray-stat/backend/tasks"
	"v2ray-stat/backend/users"
	"v2ray-stat/constant"

	_ "github.com/mattn/go-sqlite3"
)

// startAPIServer starts the API server.
func startAPIServer(ctx context.Context, manager *manager.DatabaseManager, taskManager *tasks.TaskManager, nodeClients []*db.NodeClient, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := fmt.Sprintf("%s:%d", cfg.V2RS.Address, cfg.V2RS.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: api.WithServerHeader(http.DefaultServeMux),
	}

	// Full working API endpoints
	http.HandleFunc("/", api.Answer())
	http.HandleFunc("/api/v1/users", api.UsersHandler(manager, cfg))
	http.HandleFunc("/api/v1/client_stats", api.ClientStatsHandler(manager, cfg))
	http.HandleFunc("/api/v1/server_stats", api.ServerStatsHandler(manager, cfg))

	http.HandleFunc("/api/v1/nodes", api.NodesHandler(manager, cfg))
	http.HandleFunc("/api/v1/nodes/summary", api.NodesSummaryHandler(manager, cfg))
	http.HandleFunc("/api/v1/nodes/", api.NodeByUUIDHandler(manager, cfg))
	http.HandleFunc("/api/v1/users-list", api.UsersAPIHandler(manager, cfg))
	http.HandleFunc("/api/v1/users-list/", api.UserByUUIDHandler(manager, cfg))
	http.HandleFunc("/api/v1/users-list/create", api.UsersCreateHandler(manager, cfg))

	http.HandleFunc("/api/v1/task_status", api.TokenAuthMiddleware(cfg, api.TaskStatusHandler(taskManager, cfg)))

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbDir := filepath.Dir(cfg.Paths.Database)
	if err := os.MkdirAll(dbDir, 0775); err != nil {
		cfg.Logger.Error("CRITICAL: Permission denied. Cannot create database directory",
			"path", dbDir,
			"error", err)
		os.Exit(1)
	}

	memDB, fileDB, err := db.InitDatabase(&cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to initialize database", "error", err)
	}

	manager, err := manager.NewDatabaseManager(memDB, ctx, 1, 300, 500, &cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to create DatabaseManager", "error", err)
	}

	// Node clients are now managed dynamically by NodeMonitor
	// Legacy InitNodeClients returns empty slice
	nodeClients, err := db.InitNodeClients(&cfg)
	if err != nil {
		cfg.Logger.Debug("InitNodeClients returned", "error", err)
	}

	// Initialize task manager and worker
	taskManager := tasks.NewTaskManager(manager, &cfg)
	taskWorker := tasks.NewTaskWorker(taskManager, manager, nodeClients, &cfg)
	taskWorker.Start()

	// Create and start node monitor (dynamically manages nodes from DB)
	nodeMonitor := users.NewNodeMonitor(manager, &cfg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Prepare wg
	var wg sync.WaitGroup
	wg.Add(3)

	go startAPIServer(ctx, manager, taskManager, nodeClients, &cfg, &wg)
	go db.MonitorSubscriptionsAndSync(ctx, manager, fileDB, &cfg, &wg)
	go nodeMonitor.Start(ctx, &wg)
	if cfg.Subscription != nil {
		wg.Add(1)
		go grpcclient.StartGrpcClient(ctx, &cfg, manager, &wg)
	} else {
		cfg.Logger.Info("Subscription is disabled, skipping gRPC client")
	}

	log.Printf("[START] v2ray-stat-backend application %s", constant.Version)

	<-sigChan
	cfg.Logger.Info("Received termination signal, saving data")
	cancel()

	// Stop task worker
	taskWorker.Stop()

	// Stop node monitor
	nodeMonitor.Stop()

	// Закрываем все gRPC-соединения
	for _, nc := range nodeClients {
		if err := nc.Close(); err != nil {
			cfg.Logger.Error("Failed to close node client", "node_name", nc.NodeName, "error", err)
		}
	}

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

	// Финальная синхронизация базы данных
	syncCtx, syncCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer syncCancel()
	if _, err := os.Stat(cfg.Paths.Database); os.IsNotExist(err) {
		cfg.Logger.Warn("File database does not exist, recreating", "path", cfg.Paths.Database)
		if fileDB, err = db.OpenAndInitDB(cfg.Paths.Database, "file", &cfg); err != nil {
			cfg.Logger.Error("Failed to recreate file database", "path", cfg.Paths.Database, "error", err)
		} else {
			defer fileDB.Close()
			if err := manager.SyncDBWithContext(syncCtx, fileDB, "memory to file"); err != nil {
				cfg.Logger.Error("Failed to perform final database synchronization", "error", err)
			} else {
				cfg.Logger.Info("Database synchronized successfully (memory to file)")
			}
		}
	} else {
		if err := manager.SyncDBWithContext(syncCtx, fileDB, "memory to file"); err != nil {
			cfg.Logger.Error("Failed to perform final database synchronization", "error", err)
		} else {
			cfg.Logger.Info("Database synchronized successfully (memory to file)")
		}
	}

	// Закрываем базы данных после синхронизации
	if err := memDB.Close(); err != nil {
		cfg.Logger.Error("Failed to close in-memory database", "error", err)
	}
	if err := fileDB.Close(); err != nil {
		cfg.Logger.Error("Failed to close file database", "error", err)
	}

	log.Printf("[STOP] Program terminated")
}
