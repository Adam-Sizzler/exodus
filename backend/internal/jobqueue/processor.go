package jobqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultPrefix = "exodus:jobqueue"

type Handler func(ctx context.Context, job Job) error

type Processor struct {
	client *redis.Client
	cfg    *config.BackendConfig
	prefix string

	mu     sync.RWMutex
	queues map[string]*queueRuntime
}

type QueueOptions struct {
	Name              string
	Concurrency       int
	VisibilityTimeout time.Duration
	SchedulerInterval time.Duration
	Retention         int64
}

type JobOptions struct {
	ID       string
	DedupeID string
	Delay    time.Duration
	Attempts int
	Backoff  time.Duration
}

type Job struct {
	ID           string          `json:"id"`
	Queue        string          `json:"queue"`
	Name         string          `json:"name"`
	Payload      json.RawMessage `json:"payload"`
	DedupeKey    string          `json:"dedupeKey,omitempty"`
	Attempts     int             `json:"attempts"`
	AttemptsMade int             `json:"attemptsMade"`
	BackoffMS    int64           `json:"backoffMs"`
	CreatedAt    int64           `json:"createdAt"`
}

type queueRuntime struct {
	options  QueueOptions
	handlers map[string]Handler
}

func NewRedisClient(cfg *config.BackendConfig) (*redis.Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	if cfg.Redis.Host == "" && strings.TrimSpace(cfg.Redis.Socket) == "" {
		return nil, nil
	}

	network := "tcp"
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if socket := strings.TrimSpace(cfg.Redis.Socket); socket != "" {
		network = "unix"
		addr = socket
	}

	client := redis.NewClient(&redis.Options{
		Network:  network,
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return client, nil
}

func NewProcessor(client *redis.Client, cfg *config.BackendConfig) *Processor {
	return &Processor{
		client: client,
		cfg:    cfg,
		prefix: defaultPrefix,
		queues: make(map[string]*queueRuntime),
	}
}

func (p *Processor) RegisterQueue(options QueueOptions, handlers map[string]Handler) error {
	if p == nil || p.client == nil {
		return fmt.Errorf("job queue processor is not initialized")
	}
	if strings.TrimSpace(options.Name) == "" {
		return fmt.Errorf("queue name is required")
	}
	if len(handlers) == 0 {
		return fmt.Errorf("queue %s has no handlers", options.Name)
	}
	if options.Concurrency <= 0 {
		options.Concurrency = 1
	}
	if options.VisibilityTimeout <= 0 {
		options.VisibilityTimeout = 5 * time.Minute
	}
	if options.SchedulerInterval <= 0 {
		options.SchedulerInterval = time.Second
	}
	if options.Retention <= 0 {
		options.Retention = 500
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.queues[options.Name] = &queueRuntime{options: options, handlers: handlers}
	return nil
}

func (p *Processor) Start(ctx context.Context, wg *sync.WaitGroup) {
	if p == nil || p.client == nil {
		return
	}

	p.mu.RLock()
	queues := make([]*queueRuntime, 0, len(p.queues))
	for _, queue := range p.queues {
		queues = append(queues, queue)
	}
	p.mu.RUnlock()

	for _, queue := range queues {
		p.startQueue(ctx, wg, queue)
	}
}

func (p *Processor) Close() error {
	if p == nil || p.client == nil {
		return nil
	}
	return p.client.Close()
}

func (p *Processor) Enqueue(ctx context.Context, queueName, jobName string, payload any, options JobOptions) (bool, error) {
	if p == nil || p.client == nil {
		return false, nil
	}
	queue := p.queue(queueName)
	if queue == nil {
		return false, fmt.Errorf("queue %s is not registered", queueName)
	}
	if _, ok := queue.handlers[jobName]; !ok {
		return false, fmt.Errorf("handler %s is not registered for queue %s", jobName, queueName)
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}
	if options.ID == "" {
		options.ID = uuid.NewString()
	}
	if options.Attempts <= 0 {
		options.Attempts = 1
	}
	if options.Backoff < 0 {
		options.Backoff = 0
	}

	job := Job{
		ID:        options.ID,
		Queue:     queueName,
		Name:      jobName,
		Payload:   payloadBytes,
		Attempts:  options.Attempts,
		BackoffMS: options.Backoff.Milliseconds(),
		CreatedAt: time.Now().UnixMilli(),
	}
	if options.DedupeID != "" {
		job.DedupeKey = p.dedupeKey(queueName, options.DedupeID)
		ok, err := p.client.SetNX(ctx, job.DedupeKey, job.ID, 0).Result()
		if err != nil || !ok {
			return ok, err
		}
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		p.releaseDedupe(context.Background(), job.DedupeKey)
		return false, err
	}

	if options.Delay > 0 {
		err = p.client.ZAdd(ctx, p.delayedKey(queueName), redis.Z{
			Score:  float64(time.Now().Add(options.Delay).UnixMilli()),
			Member: string(jobBytes),
		}).Err()
	} else {
		err = p.client.RPush(ctx, p.waitingKey(queueName), string(jobBytes)).Err()
	}
	if err != nil {
		p.releaseDedupe(context.Background(), job.DedupeKey)
		return false, err
	}

	return true, nil
}

func (p *Processor) queue(name string) *queueRuntime {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.queues[name]
}

func (p *Processor) startQueue(ctx context.Context, wg *sync.WaitGroup, queue *queueRuntime) {
	if p.cfg != nil && p.cfg.Logger != nil {
		p.cfg.Logger.Info("Job queue started", "queue", queue.options.Name, "concurrency", queue.options.Concurrency)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		p.schedulerLoop(ctx, queue)
	}()

	for i := 0; i < queue.options.Concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			p.workerLoop(ctx, queue, workerID)
		}(i + 1)
	}
}

func (p *Processor) schedulerLoop(ctx context.Context, queue *queueRuntime) {
	ticker := time.NewTicker(queue.options.SchedulerInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.moveDueDelayed(ctx, queue.options.Name)
			p.recoverExpiredActive(ctx, queue)
		}
	}
}

func (p *Processor) workerLoop(ctx context.Context, queue *queueRuntime, workerID int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, ok, err := p.reserveJob(ctx, queue)
		if err != nil {
			p.warn("Failed to reserve job", "queue", queue.options.Name, "worker", workerID, "error", err)
			sleepContext(ctx, 250*time.Millisecond)
			continue
		}
		if !ok {
			sleepContext(ctx, 250*time.Millisecond)
			continue
		}

		p.processJob(ctx, queue, job)
	}
}

func (p *Processor) reserveJob(ctx context.Context, queue *queueRuntime) (Job, bool, error) {
	script := redis.NewScript(`
local job = redis.call("LPOP", KEYS[1])
if not job then
  return nil
end
local decoded = cjson.decode(job)
redis.call("HSET", KEYS[2], decoded["id"], job)
redis.call("ZADD", KEYS[3], ARGV[1], decoded["id"])
return job
`)
	res, err := script.Run(ctx, p.client, []string{
		p.waitingKey(queue.options.Name),
		p.activeHashKey(queue.options.Name),
		p.activeZSetKey(queue.options.Name),
	}, time.Now().Add(queue.options.VisibilityTimeout).UnixMilli()).Result()
	if err == redis.Nil || res == nil {
		return Job{}, false, nil
	}
	if err != nil {
		return Job{}, false, err
	}

	jobStr, ok := res.(string)
	if !ok {
		return Job{}, false, fmt.Errorf("unexpected redis job payload type %T", res)
	}

	var job Job
	if err := json.Unmarshal([]byte(jobStr), &job); err != nil {
		return Job{}, false, err
	}
	return job, true, nil
}

func (p *Processor) processJob(ctx context.Context, queue *queueRuntime, job Job) {
	handler, ok := queue.handlers[job.Name]
	if !ok {
		_ = p.failJob(context.Background(), queue, job, fmt.Errorf("handler %s is not registered", job.Name))
		return
	}

	if err := handler(ctx, job); err != nil {
		if job.AttemptsMade+1 < job.Attempts {
			job.AttemptsMade++
			_ = p.retryJob(context.Background(), queue, job, err)
			return
		}
		_ = p.failJob(context.Background(), queue, job, err)
		return
	}

	if err := p.completeJob(context.Background(), queue, job); err != nil {
		p.warn("Failed to complete job", "queue", job.Queue, "job", job.Name, "id", job.ID, "error", err)
	}
}

func (p *Processor) retryJob(ctx context.Context, queue *queueRuntime, job Job, cause error) error {
	if err := p.removeActive(ctx, job.Queue, job.ID); err != nil {
		return err
	}

	jobBytes, err := json.Marshal(job)
	if err != nil {
		return err
	}
	delay := time.Duration(job.BackoffMS) * time.Millisecond
	if delay <= 0 {
		delay = time.Second
	}
	if err := p.client.ZAdd(ctx, p.delayedKey(job.Queue), redis.Z{
		Score:  float64(time.Now().Add(delay).UnixMilli()),
		Member: string(jobBytes),
	}).Err(); err != nil {
		return err
	}
	p.warn("Job scheduled for retry", "queue", job.Queue, "job", job.Name, "id", job.ID, "attempt", job.AttemptsMade, "error", cause)
	return nil
}

func (p *Processor) completeJob(ctx context.Context, queue *queueRuntime, job Job) error {
	if err := p.removeActive(ctx, job.Queue, job.ID); err != nil {
		return err
	}
	p.releaseDedupe(ctx, job.DedupeKey)
	return p.recordRetention(ctx, p.completedKey(job.Queue), job.ID, queue.options.Retention)
}

func (p *Processor) failJob(ctx context.Context, queue *queueRuntime, job Job, cause error) error {
	if err := p.removeActive(ctx, job.Queue, job.ID); err != nil {
		return err
	}
	p.releaseDedupe(ctx, job.DedupeKey)
	p.warn("Job failed", "queue", job.Queue, "job", job.Name, "id", job.ID, "attempts", job.AttemptsMade+1, "error", cause)
	return p.recordRetention(ctx, p.failedKey(job.Queue), job.ID, queue.options.Retention)
}

func (p *Processor) removeActive(ctx context.Context, queueName, jobID string) error {
	pipe := p.client.Pipeline()
	pipe.HDel(ctx, p.activeHashKey(queueName), jobID)
	pipe.ZRem(ctx, p.activeZSetKey(queueName), jobID)
	_, err := pipe.Exec(ctx)
	return err
}

func (p *Processor) moveDueDelayed(ctx context.Context, queueName string) {
	now := fmt.Sprintf("%d", time.Now().UnixMilli())
	jobs, err := p.client.ZRangeByScore(ctx, p.delayedKey(queueName), &redis.ZRangeBy{
		Min:   "-inf",
		Max:   now,
		Count: 500,
	}).Result()
	if err != nil {
		p.warn("Failed to read delayed jobs", "queue", queueName, "error", err)
		return
	}
	for _, job := range jobs {
		removed, err := p.client.ZRem(ctx, p.delayedKey(queueName), job).Result()
		if err != nil || removed == 0 {
			continue
		}
		if err := p.client.RPush(ctx, p.waitingKey(queueName), job).Err(); err != nil {
			p.warn("Failed to move delayed job to waiting", "queue", queueName, "error", err)
		}
	}
}

func (p *Processor) recoverExpiredActive(ctx context.Context, queue *queueRuntime) {
	now := fmt.Sprintf("%d", time.Now().UnixMilli())
	jobIDs, err := p.client.ZRangeByScore(ctx, p.activeZSetKey(queue.options.Name), &redis.ZRangeBy{
		Min:   "-inf",
		Max:   now,
		Count: 100,
	}).Result()
	if err != nil {
		p.warn("Failed to read active jobs", "queue", queue.options.Name, "error", err)
		return
	}
	for _, jobID := range jobIDs {
		jobData, err := p.client.HGet(ctx, p.activeHashKey(queue.options.Name), jobID).Result()
		if err == redis.Nil {
			_ = p.client.ZRem(ctx, p.activeZSetKey(queue.options.Name), jobID).Err()
			continue
		}
		if err != nil {
			continue
		}
		if err := p.removeActive(ctx, queue.options.Name, jobID); err != nil {
			continue
		}
		if err := p.client.RPush(ctx, p.waitingKey(queue.options.Name), jobData).Err(); err != nil {
			p.warn("Failed to requeue expired active job", "queue", queue.options.Name, "job_id", jobID, "error", err)
		}
	}
}

func (p *Processor) recordRetention(ctx context.Context, key, jobID string, retention int64) error {
	if retention <= 0 {
		return nil
	}
	if err := p.client.ZAdd(ctx, key, redis.Z{Score: float64(time.Now().UnixMilli()), Member: jobID}).Err(); err != nil {
		return err
	}
	count, err := p.client.ZCard(ctx, key).Result()
	if err != nil || count <= retention {
		return err
	}
	oldIDs, err := p.client.ZRange(ctx, key, 0, count-retention-1).Result()
	if err != nil || len(oldIDs) == 0 {
		return err
	}
	members := make([]any, 0, len(oldIDs))
	for _, id := range oldIDs {
		members = append(members, id)
	}
	return p.client.ZRem(ctx, key, members...).Err()
}

func (p *Processor) releaseDedupe(ctx context.Context, key string) {
	if key == "" {
		return
	}
	_ = p.client.Del(ctx, key).Err()
}

func (p *Processor) waitingKey(queue string) string {
	return p.prefix + ":" + queue + ":waiting"
}

func (p *Processor) delayedKey(queue string) string {
	return p.prefix + ":" + queue + ":delayed"
}

func (p *Processor) activeHashKey(queue string) string {
	return p.prefix + ":" + queue + ":active"
}

func (p *Processor) activeZSetKey(queue string) string {
	return p.prefix + ":" + queue + ":active_deadline"
}

func (p *Processor) completedKey(queue string) string {
	return p.prefix + ":" + queue + ":completed"
}

func (p *Processor) failedKey(queue string) string {
	return p.prefix + ":" + queue + ":failed"
}

func (p *Processor) dedupeKey(queue, id string) string {
	return p.prefix + ":" + queue + ":dedupe:" + id
}

func (p *Processor) warn(msg string, args ...any) {
	if p != nil && p.cfg != nil && p.cfg.Logger != nil {
		p.cfg.Logger.Warn(msg, args...)
	}
}

func sleepContext(ctx context.Context, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
