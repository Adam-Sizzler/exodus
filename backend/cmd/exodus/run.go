package exodus

import (
	"context"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync"
	"syscall"
	"time"

	"exodus/internal/config"
	"exodus/internal/constant"
	"exodus/internal/db"
	"exodus/internal/httpapi/panelsettings"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"
	users "exodus/internal/nodes"
	"exodus/internal/notifications"
	"exodus/internal/redisqueue"
	"exodus/internal/scheduler"
	"exodus/internal/subscriptionnodes"
	"exodus/internal/userwatchdog"
)

func Run() {
	flags := parseCLIFlags()
	if runPreConfigCLI(flags) {
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fallback, _ := logger.NewLoggerFromEnv("info", "console", "UTC", os.Stderr)
		fallback.RoleService(logger.RoleAPI, logger.ServiceConfig).Error("Environment validation failed", "error", err)
		os.Exit(1)
	}
	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceBootstrap).Info("Starting application")

	if runConfiguredCLI(flags, &cfg) {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConn, err := db.InitDatabase(&cfg)
	if err != nil {
		cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceDatabase).Fatal("Failed to initialize database", "error", err)
	}
	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceDatabase).Info("Connected to PostgreSQL")

	pools, err := db.NewPools(ctx, dbConn, &cfg)
	if err != nil {
		cfg.Logger.Fatal("Failed to create database pools", "error", err)
	}

	// Create and start node monitor (dynamically manages nodes from DB)
	nodeMonitor := users.NewNodeMonitor(pools.Interactive, &cfg)
	users.RegisterGlobalNodeMonitor(nodeMonitor)
	subNodeMonitor := subscriptionnodes.NewSubNodeMonitor(pools.Interactive, &cfg)
	subscriptionnodes.RegisterGlobalSubNodeMonitor(subNodeMonitor)

	redisWorker, err := redisqueue.NewWorker(&cfg, pools.Background)
	redisStatus := "Disabled"
	if err != nil {
		cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceRedis).Warn("Redis worker disabled", "error", err)
	} else if redisWorker != nil {
		redisStatus = "Connected"
		cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceRedis).Info("Connected to Redis")
		nodeMonitor.SetNodeUserUsageRecorder(redisWorker)
	} else {
		cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceRedis).Warn("Redis disabled: REDIS_HOST or REDIS_SOCKET is not configured")
	}

	cfg.Logger.PrintStartupBanner(logger.BannerOptions{
		Title:          "Exodus Backend",
		Version:        constant.Version,
		DocsURL:        "https://docs.exodus.dev",
		CommunityURL:   "https://t.me/exodus",
		HTTPPort:       cfg.Panel.AppPort,
		PathPrefix:     cfg.Panel.BasePath,
		DatabaseStatus: "Connected",
		RedisStatus:    redisStatus,
		RescueCLI:      "docker exec -it exodus cli",
	})

	panelsettings.LogScopeCatalog(&cfg)

	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHealthCheck).Info("Health checks initialized")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Prepare wg
	var wg sync.WaitGroup
	wg.Add(4)

	go startWebServer(ctx, pools, &cfg, &wg)
	go startMetricsServer(ctx, pools.Interactive, &cfg, &wg)
	go nodeMonitor.Start(ctx, &wg)
	go subNodeMonitor.Start(ctx, &wg)
	notifications.StartDispatcher(ctx, &wg, pools.Background, &cfg)
	userwatchdog.Start(ctx, &wg, pools.Background, &cfg)
	scheduler.Start(ctx, &wg, pools.Background, &cfg)
	if _, err := jobqueue.StartSubscriptionQueues(ctx, &wg, pools.Background, &cfg); err != nil {
		cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceUsersQueue).Warn("Subscription job queue disabled", "error", err)
	}
	if redisWorker != nil {
		redisWorker.Start(ctx, &wg)
	}

	notifications.Emit(context.Background(), &cfg, notifications.Event{
		Scope: notifications.ScopeService,
		Event: notifications.EventServicePanelStarted,
		Data: map[string]any{
			"version": constant.Version,
		},
	})

	<-sigChan
	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceBootstrap).Info("Received termination signal, saving data")
	cancel()

	// Stop node monitor
	nodeMonitor.Stop()
	subNodeMonitor.Stop()

	// Close database pools
	if err := pools.Close(); err != nil {
		cfg.Logger.Warn("Error closing database pools", "error", err)
	}

	// Wait for goroutines to finish
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

	if redisWorker != nil {
		if err := redisWorker.Close(); err != nil {
			cfg.Logger.Warn("Failed to close redis worker", "error", err)
		}
	}

	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceBootstrap).Info("Program terminated")
}
