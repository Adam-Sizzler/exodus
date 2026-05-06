package redisqueue

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"

	"github.com/redis/go-redis/v9"
)

const (
	scheduledKey          = "exodus:push_to_db:scheduled"
	recordUserUsagePrefix = "recordUserUsage:"
	processingPostfix     = ":processing"
)

type Worker struct {
	client       *redis.Client
	manager      *dbmanager.DatabaseManager
	cfg          *config.BackendConfig
	pollInterval time.Duration
	delay        time.Duration
}

func NewWorker(cfg *config.BackendConfig, manager *dbmanager.DatabaseManager) (*Worker, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if cfg.Redis.Host == "" {
		return nil, nil
	}

	opts := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
	}

	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Worker{
		client:       client,
		manager:      manager,
		cfg:          cfg,
		pollInterval: 2 * time.Second,
		delay:        2 * time.Minute,
	}, nil
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	if w == nil {
		return
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.cfg.Logger.Info("Redis worker started")
		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.cfg.Logger.Info("Redis worker stopped")
				return
			case <-ticker.C:
				w.processDueJobs(ctx)
			}
		}
	}()
}

func (w *Worker) Close() error {
	if w == nil || w.client == nil {
		return nil
	}
	return w.client.Close()
}

// RecordUserUsageDelayed schedules a recordUserUsage job for the given redisKey.
func (w *Worker) RecordUserUsageDelayed(ctx context.Context, redisKey string) error {
	if w == nil {
		return fmt.Errorf("redis worker is not initialized")
	}
	member := recordUserUsagePrefix + redisKey
	score := float64(time.Now().Add(w.delay).Unix())
	return w.client.ZAddArgs(ctx, scheduledKey, redis.ZAddArgs{
		NX:      true,
		Members: []redis.Z{{Score: score, Member: member}},
	}).Err()
}

func (w *Worker) processDueJobs(ctx context.Context) {
	if w == nil {
		return
	}

	now := fmt.Sprintf("%d", time.Now().Unix())
	jobs, err := w.client.ZRangeByScore(ctx, scheduledKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   now,
		Count: 100,
	}).Result()
	if err != nil {
		w.cfg.Logger.Warn("Failed to fetch redis jobs", "error", err)
		return
	}
	for _, job := range jobs {
		removed, err := w.client.ZRem(ctx, scheduledKey, job).Result()
		if err != nil || removed == 0 {
			continue
		}
		if strings.HasPrefix(job, recordUserUsagePrefix) {
			redisKey := strings.TrimPrefix(job, recordUserUsagePrefix)
			if err := w.handleRecordUserUsage(ctx, redisKey); err != nil {
				w.cfg.Logger.Warn("Failed to handle recordUserUsage job", "error", err, "redis_key", redisKey)
			}
		}
	}
}

func (w *Worker) handleRecordUserUsage(ctx context.Context, redisKey string) error {
	processingKey := redisKey + processingPostfix

	exists, err := w.client.Exists(ctx, redisKey).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return nil
	}

	if err := w.client.Rename(ctx, redisKey, processingKey).Err(); err != nil {
		return err
	}
	defer func() {
		_ = w.client.Del(context.Background(), processingKey).Err()
	}()

	data, err := w.client.HGetAll(ctx, processingKey).Result()
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}

	nodeID, err := parseNodeID(redisKey)
	if err != nil {
		return err
	}

	return w.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO nodes_user_usage_history (node_id, user_id, total_bytes)
			VALUES (?, ?, ?)
			ON CONFLICT (node_id, created_at, user_id)
			DO UPDATE SET total_bytes = EXCLUDED.total_bytes, updated_at = now()`

		for userIDStr, totalBytesStr := range data {
			userID, err := strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				continue
			}
			totalBytes, err := strconv.ParseInt(totalBytesStr, 10, 64)
			if err != nil {
				continue
			}
			if _, err := tx.ExecContext(ctx, query, nodeID, userID, totalBytes); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
}

func parseNodeID(redisKey string) (int64, error) {
	parts := strings.Split(redisKey, ":")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid redis key: %s", redisKey)
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid node id in key %s", redisKey)
	}
	return id, nil
}
