package notifications

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
)

type Dispatcher struct {
	cfg      *config.BackendConfig
	notifier *Notifier
	worker   *Worker
}

var (
	globalDispatcherMu sync.RWMutex
	globalDispatcher   *Dispatcher
)

func StartDispatcher(ctx context.Context, wg *sync.WaitGroup, db *sql.DB, cfg *config.BackendConfig) {
	if cfg == nil {
		return
	}
	worker, err := NewWorker(cfg)
	if err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Warn("Notification delivery queue disabled", "error", err)
		}
		return
	}

	dispatcher := &Dispatcher{
		cfg:      cfg,
		notifier: New(cfg),
		worker:   worker,
	}

	globalDispatcherMu.Lock()
	globalDispatcher = dispatcher
	globalDispatcherMu.Unlock()

	worker.Start(ctx, wg)
}

func enqueueWithGlobalDispatcher(ctx context.Context, event Event) bool {
	globalDispatcherMu.RLock()
	dispatcher := globalDispatcher
	globalDispatcherMu.RUnlock()
	if dispatcher == nil {
		return false
	}
	if err := dispatcher.Enqueue(ctx, event); err != nil {
		if dispatcher.cfg != nil && dispatcher.cfg.Logger != nil {
			dispatcher.cfg.Logger.Warn("Failed to enqueue notification", "event", event.Event, "error", err)
		}
		return false
	}
	return true
}

func (d *Dispatcher) Enqueue(ctx context.Context, event Event) error {
	if d == nil || d.worker == nil {
		return nil
	}
	if strings.TrimSpace(event.Timestamp) == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}
	if ctx == nil {
		ctx = context.Background()
	}

	if d.notifier.webhookEnabled() && d.cfg.Notifications.EventChannelEnabled(event.Event, "webhook") {
		_ = d.worker.EnqueueWebhook(ctx, event)
	}

	if d.notifier.telegramEnabled() && d.cfg.Notifications.EventChannelEnabled(event.Event, "telegram") {
		_ = d.worker.EnqueueTelegram(ctx, event)
	}

	return nil
}
