package srslists

import (
	"context"
	"sync"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
)

func StartPeriodicChecker(ctx context.Context, wg *sync.WaitGroup, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, interval time.Duration) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if interval <= 0 {
			interval = 5 * time.Minute
		}
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		cfg.Logger.Info("SRS availability checker started", "interval", interval.String())
		_, _ = CheckAndUpdateAvailability(ctx, manager, cfg)

		for {
			select {
			case <-ctx.Done():
				cfg.Logger.Info("SRS availability checker stopped")
				return
			case <-ticker.C:
				if _, err := CheckAndUpdateAvailability(ctx, manager, cfg); err != nil {
					cfg.Logger.Warn("Periodic SRS availability check failed", "error", err)
				}
			}
		}
	}()
}
