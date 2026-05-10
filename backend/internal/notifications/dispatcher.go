package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
)

const (
	notificationQueueMaxAttempts = 8
	notificationQueuePollEvery   = 5 * time.Second
)

type Dispatcher struct {
	manager  *dbmanager.DatabaseManager
	cfg      *config.BackendConfig
	notifier *Notifier
	wake     chan struct{}
}

var (
	globalDispatcherMu sync.RWMutex
	globalDispatcher   *Dispatcher
)

func StartDispatcher(ctx context.Context, wg *sync.WaitGroup, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	if manager == nil || cfg == nil {
		return
	}
	dispatcher := &Dispatcher{
		manager:  manager,
		cfg:      cfg,
		notifier: New(cfg),
		wake:     make(chan struct{}, 1),
	}
	if err := dispatcher.ensureQueueTable(ctx); err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Warn("Notification delivery queue disabled", "error", err)
		}
		return
	}

	globalDispatcherMu.Lock()
	globalDispatcher = dispatcher
	globalDispatcherMu.Unlock()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatcher.run(ctx)
	}()
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
	if d == nil || d.manager == nil {
		return nil
	}
	if strings.TrimSpace(event.Timestamp) == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	err = d.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(ctx, `
			INSERT INTO notification_delivery_queue (event, attempts, next_attempt_at)
			VALUES (?::jsonb, 0, CURRENT_TIMESTAMP)
		`, string(payload))
		return execErr
	})
	if err != nil {
		return err
	}
	select {
	case d.wake <- struct{}{}:
	default:
	}
	return nil
}

func (d *Dispatcher) ensureQueueTable(ctx context.Context) error {
	return d.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			CREATE TABLE IF NOT EXISTS notification_delivery_queue (
				id bigserial PRIMARY KEY,
				event jsonb NOT NULL,
				attempts integer NOT NULL DEFAULT 0,
				next_attempt_at timestamp(3) without time zone NOT NULL DEFAULT now(),
				last_error text,
				failed_at timestamp(3) without time zone,
				created_at timestamp(3) without time zone NOT NULL DEFAULT now(),
				updated_at timestamp(3) without time zone NOT NULL DEFAULT now()
			)
		`)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(ctx, `
			CREATE INDEX IF NOT EXISTS notification_delivery_queue_due_idx
			ON notification_delivery_queue (failed_at, next_attempt_at, id)
		`)
		return err
	})
}

func (d *Dispatcher) run(ctx context.Context) {
	if d.cfg != nil && d.cfg.Logger != nil {
		d.cfg.Logger.Info("Notification delivery dispatcher started")
	}
	ticker := time.NewTicker(notificationQueuePollEvery)
	defer ticker.Stop()

	for {
		d.drain(ctx)
		select {
		case <-ctx.Done():
			if d.cfg != nil && d.cfg.Logger != nil {
				d.cfg.Logger.Info("Notification delivery dispatcher stopped")
			}
			return
		case <-d.wake:
		case <-ticker.C:
		}
	}
}

func (d *Dispatcher) drain(ctx context.Context) {
	for {
		delivered, err := d.processOne(ctx)
		if err != nil {
			if d.cfg != nil && d.cfg.Logger != nil {
				d.cfg.Logger.Warn("Notification delivery attempt failed", "error", err)
			}
			return
		}
		if !delivered {
			return
		}
	}
}

func (d *Dispatcher) processOne(ctx context.Context) (bool, error) {
	item, found, err := d.claimOne(ctx)
	if err != nil || !found {
		return found, err
	}

	sendCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	if err := d.notifier.Send(sendCtx, item.Event); err != nil {
		return true, d.markFailed(ctx, item.ID, item.Attempts+1, err)
	}
	return true, d.delete(ctx, item.ID)
}

type queuedNotification struct {
	ID       int64
	Event    Event
	Attempts int
}

func (d *Dispatcher) claimOne(ctx context.Context) (queuedNotification, bool, error) {
	var item queuedNotification
	var raw string
	found := false
	err := d.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		row := tx.QueryRowContext(ctx, `
			SELECT id, event::text, attempts
			FROM notification_delivery_queue
			WHERE failed_at IS NULL
				AND next_attempt_at <= CURRENT_TIMESTAMP
			ORDER BY next_attempt_at ASC, id ASC
			LIMIT 1
			FOR UPDATE SKIP LOCKED
		`)
		if scanErr := row.Scan(&item.ID, &raw, &item.Attempts); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return nil
			}
			return scanErr
		}
		found = true
		return tx.Commit()
	})
	if err != nil || !found {
		return item, found, err
	}
	if err := json.Unmarshal([]byte(raw), &item.Event); err != nil {
		return item, true, d.markFailed(ctx, item.ID, notificationQueueMaxAttempts, err)
	}
	return item, true, nil
}

func (d *Dispatcher) delete(ctx context.Context, id int64) error {
	return d.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `DELETE FROM notification_delivery_queue WHERE id = ?`, id)
		return err
	})
}

func (d *Dispatcher) markFailed(ctx context.Context, id int64, attempts int, sendErr error) error {
	delay := notificationRetryDelay(attempts, sendErr)
	message := ""
	if sendErr != nil {
		message = sendErr.Error()
	}
	if attempts >= notificationQueueMaxAttempts {
		return d.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
			_, err := db.ExecContext(ctx, `
				UPDATE notification_delivery_queue
				SET attempts = ?, last_error = ?, failed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
				WHERE id = ?
			`, attempts, message, id)
			return err
		})
	}
	return d.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			UPDATE notification_delivery_queue
			SET attempts = ?, last_error = ?, next_attempt_at = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ?
		`, attempts, message, time.Now().UTC().Add(delay), id)
		return err
	})
}

func notificationRetryDelay(attempts int, err error) time.Duration {
	var rateLimit RateLimitError
	if errors.As(err, &rateLimit) && rateLimit.RetryAfter > 0 {
		return clampDuration(rateLimit.RetryAfter, time.Second, time.Hour)
	}
	if attempts < 1 {
		attempts = 1
	}
	if attempts > 6 {
		attempts = 6
	}
	delay := time.Duration(1<<uint(attempts-1)) * time.Minute
	return clampDuration(delay, time.Minute, 30*time.Minute)
}

func clampDuration(value, min, max time.Duration) time.Duration {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
