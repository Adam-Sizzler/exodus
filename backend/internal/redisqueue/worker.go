package redisqueue

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"

	"github.com/redis/go-redis/v9"
)

const (
	scheduledKey           = "exodus:push_to_db:scheduled"
	recordUserUsagePrefix  = "recordUserUsage:"
	nodeUserUsagePrefix    = "node_user_usage:"
	processingPostfix      = ":processing"
	nodeUserUsageBatchSize = 10000
)

type Worker struct {
	client       *redis.Client
	manager      *dbmanager.DatabaseManager
	cfg          *config.BackendConfig
	pollInterval time.Duration
	delay        time.Duration
	usageTTL     time.Duration
}

type nodeUsageEntry struct {
	UserID     int64
	TotalBytes int64
}

func NewWorker(cfg *config.BackendConfig, manager *dbmanager.DatabaseManager) (*Worker, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if cfg.Redis.Host == "" && cfg.Redis.Socket == "" {
		return nil, nil
	}

	network := "tcp"
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if strings.TrimSpace(cfg.Redis.Socket) != "" {
		network = "unix"
		addr = strings.TrimSpace(cfg.Redis.Socket)
	}

	opts := &redis.Options{
		Network:  network,
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
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
		delay:        time.Duration(cfg.Redis.UserUsageHistoryDelaySeconds) * time.Second,
		usageTTL:     time.Duration(cfg.Redis.UserUsageHistoryTTLSeconds) * time.Second,
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
		if w.cfg.Redis.DisableUserUsageRecords {
			w.cfg.Logger.Warn("SERVICE_DISABLE_USER_USAGE_RECORDS is enabled, node user usage history will not be recorded")
		}
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

// RecordNodeUserUsage accumulates node user usage in Redis and schedules a delayed flush to PostgreSQL.
func (w *Worker) RecordNodeUserUsage(ctx context.Context, nodeID int64, userBytes map[int64]int64) error {
	if w == nil {
		return fmt.Errorf("redis worker is not initialized")
	}
	if w.cfg != nil && w.cfg.Redis.DisableUserUsageRecords {
		return nil
	}
	if nodeID <= 0 || len(userBytes) == 0 {
		return nil
	}

	redisKey := nodeUserUsageRedisKey(nodeID)
	pipe := w.client.Pipeline()
	increments := 0
	for userID, totalBytes := range userBytes {
		if userID <= 0 || totalBytes <= 0 {
			continue
		}
		pipe.HIncrBy(ctx, redisKey, strconv.FormatInt(userID, 10), totalBytes)
		increments++
	}
	if increments == 0 {
		return nil
	}
	if w.usageTTL > 0 {
		pipe.Expire(ctx, redisKey, w.usageTTL)
	}
	w.recordUserUsageDelayedPipeline(ctx, pipe, redisKey)
	_, err := pipe.Exec(ctx)
	return err
}

// RecordUserUsageDelayed schedules a recordUserUsage job for the given redisKey.
func (w *Worker) RecordUserUsageDelayed(ctx context.Context, redisKey string) error {
	if w == nil {
		return fmt.Errorf("redis worker is not initialized")
	}
	return w.recordUserUsageDelayed(ctx, redisKey)
}

func (w *Worker) recordUserUsageDelayed(ctx context.Context, redisKey string) error {
	member := recordUserUsagePrefix + redisKey
	score := float64(time.Now().Add(w.delay).Unix())
	return w.client.ZAddArgs(ctx, scheduledKey, redis.ZAddArgs{
		NX:      true,
		Members: []redis.Z{{Score: score, Member: member}},
	}).Err()
}

func (w *Worker) recordUserUsageDelayedPipeline(ctx context.Context, pipe redis.Pipeliner, redisKey string) {
	member := recordUserUsagePrefix + redisKey
	score := float64(time.Now().Add(w.delay).Unix())
	pipe.ZAddArgs(ctx, scheduledKey, redis.ZAddArgs{
		NX:      true,
		Members: []redis.Z{{Score: score, Member: member}},
	})
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
				_ = w.recordUserUsageDelayed(context.Background(), redisKey)
			}
		}
	}
}

func (w *Worker) handleRecordUserUsage(ctx context.Context, redisKey string) error {
	if w.cfg != nil && w.cfg.Redis.DisableUserUsageRecords {
		return nil
	}
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

	entries := make([]nodeUsageEntry, 0, len(data))
	for userIDStr, totalBytesStr := range data {
		userID, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil || userID <= 0 {
			continue
		}
		totalBytes, err := strconv.ParseInt(totalBytesStr, 10, 64)
		if err != nil || totalBytes <= 0 {
			continue
		}
		entries = append(entries, nodeUsageEntry{
			UserID:     userID,
			TotalBytes: totalBytes,
		})
	}
	if len(entries) == 0 {
		return nil
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].UserID < entries[j].UserID
	})

	return w.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		for start := 0; start < len(entries); start += nodeUserUsageBatchSize {
			end := start + nodeUserUsageBatchSize
			if end > len(entries) {
				end = len(entries)
			}
			if err := bulkUpsertNodeUserUsageHistory(ctx, db, nodeID, entries[start:end]); err != nil {
				return err
			}
		}
		return nil
	})
}

func bulkUpsertNodeUserUsageHistory(ctx context.Context, db dbmanager.DBExecutor, nodeID int64, entries []nodeUsageEntry) error {
	if nodeID <= 0 || len(entries) == 0 {
		return nil
	}

	var query strings.Builder
	args := make([]any, 0, len(entries)*3)
	query.WriteString(`
		INSERT INTO nodes_user_usage_history (
			node_id,
			user_id,
			total_bytes,
			created_at,
			updated_at
		)
		SELECT
			v.node_id,
			v.user_id,
			v.total_bytes,
			CURRENT_DATE,
			now()
		FROM (VALUES `)
	for i, entry := range entries {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString("(?::bigint, ?::bigint, ?::bigint)")
		args = append(args, nodeID, entry.UserID, entry.TotalBytes)
	}
	query.WriteString(`) AS v(node_id, user_id, total_bytes)
		WHERE EXISTS (SELECT 1 FROM nodes WHERE id = v.node_id)
		  AND EXISTS (SELECT 1 FROM users WHERE t_id = v.user_id)
		ON CONFLICT ON CONSTRAINT nodes_user_usage_history_pkey
		DO UPDATE SET
			total_bytes = nodes_user_usage_history.total_bytes + EXCLUDED.total_bytes,
			updated_at = EXCLUDED.updated_at
	`)

	_, err := db.ExecContext(ctx, query.String(), args...)
	return err
}

func nodeUserUsageRedisKey(nodeID int64) string {
	return fmt.Sprintf("%s%d", nodeUserUsagePrefix, nodeID)
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
