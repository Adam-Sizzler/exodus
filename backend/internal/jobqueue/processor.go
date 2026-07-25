package jobqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/logger"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

type Handler func(ctx context.Context, job Job) error

type Processor struct {
	client   *asynq.Client
	servers  []*asynq.Server
	muxs     []*asynq.ServeMux
	cfg      *config.BackendConfig
	redisOpt asynq.RedisClientOpt

	mu sync.Mutex
}

type QueueOptions struct {
	Name              string
	Concurrency       int
	VisibilityTimeout time.Duration
	SchedulerInterval time.Duration
	BlockTimeout      time.Duration
	Retention         int64
}

type JobOptions struct {
	ID        string
	DedupeID  string
	Delay     time.Duration
	Attempts  int
	Backoff   time.Duration
	Retention int64
}

type Job struct {
	ID       string          `json:"id"`
	Queue    string          `json:"queue"`
	Name     string          `json:"name"`
	Payload  json.RawMessage `json:"payload"`
	Attempts int             `json:"attempts"`
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

func BuildAsynqRedisOpt(cfg *config.BackendConfig) asynq.RedisClientOpt {
	network := "tcp"
	addr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if socket := strings.TrimSpace(cfg.Redis.Socket); socket != "" {
		network = "unix"
		addr = socket
	}
	return asynq.RedisClientOpt{
		Network:  network,
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	}
}

func NewProcessor(client *redis.Client, cfg *config.BackendConfig) *Processor {
	redisOpt := BuildAsynqRedisOpt(cfg)
	asynqClient := asynq.NewClient(redisOpt)
	return &Processor{
		client:   asynqClient,
		cfg:      cfg,
		redisOpt: redisOpt,
	}
}

func (p *Processor) RegisterQueue(options QueueOptions, handlers map[string]Handler) error {
	if p == nil {
		return fmt.Errorf("job queue processor is not initialized")
	}
	if strings.TrimSpace(options.Name) == "" {
		return fmt.Errorf("queue name is required")
	}
	if len(handlers) == 0 {
		return fmt.Errorf("queue %s has no handlers", options.Name)
	}
	concurrency := options.Concurrency
	if concurrency <= 0 {
		concurrency = 1
	}

	srv := asynq.NewServer(
		p.redisOpt,
		asynq.Config{
			Concurrency: concurrency,
			Queues: map[string]int{
				options.Name: 10,
			},
			Logger: &asynqLogger{cfg: p.cfg},
		},
	)

	mux := asynq.NewServeMux()
	for name, handler := range handlers {
		jobName := name
		h := handler
		mux.HandleFunc(jobName, func(ctx context.Context, task *asynq.Task) error {
			attempts, _ := asynq.GetRetryCount(ctx)
			taskID, _ := asynq.GetTaskID(ctx)
			queueName, _ := asynq.GetQueueName(ctx)
			job := Job{
				ID:       taskID,
				Queue:    queueName,
				Name:     jobName,
				Payload:  task.Payload(),
				Attempts: attempts,
			}
			return h(ctx, job)
		})
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers = append(p.servers, srv)
	p.muxs = append(p.muxs, mux)

	return nil
}

func (p *Processor) Start(ctx context.Context, wg *sync.WaitGroup) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.servers {
		srv := p.servers[i]
		mux := p.muxs[i]

		wg.Add(1)
		go func(s *asynq.Server, m *asynq.ServeMux) {
			defer wg.Done()
			if err := s.Run(m); err != nil {
				if p.cfg != nil && p.cfg.Logger != nil {
					p.cfg.Logger.Error("asynq server stopped", "error", err)
				}
			}
		}(srv, mux)
	}
}

func (p *Processor) Enqueue(ctx context.Context, queue string, name string, payload json.RawMessage, options JobOptions) error {
	if p == nil || p.client == nil {
		return nil
	}
	task := asynq.NewTask(name, payload)
	var opts []asynq.Option
	opts = append(opts, asynq.Queue(queue))

	if options.ID != "" {
		opts = append(opts, asynq.TaskID(options.ID))
	} else if options.DedupeID != "" {
		opts = append(opts, asynq.TaskID(options.DedupeID))
	}

	if options.Delay > 0 {
		opts = append(opts, asynq.ProcessIn(options.Delay))
	}

	if options.Attempts > 0 {
		opts = append(opts, asynq.MaxRetry(options.Attempts))
	}

	if options.Retention > 0 {
		opts = append(opts, asynq.Retention(time.Duration(options.Retention)*time.Second)) // Note: asynq uses time.Duration
	}

	_, err := p.client.EnqueueContext(ctx, task, opts...)
	return err
}

func (p *Processor) Close() error {
	if p == nil {
		return nil
	}
	var errs []string
	if p.client != nil {
		if err := p.client.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	p.mu.Lock()
	for _, srv := range p.servers {
		srv.Shutdown()
	}
	p.mu.Unlock()

	if len(errs) > 0 {
		return fmt.Errorf("errors closing processor: %s", strings.Join(errs, ", "))
	}
	return nil
}

type asynqLogger struct {
	cfg *config.BackendConfig
}

func (l *asynqLogger) Debug(args ...interface{}) {
	if l.cfg != nil && l.cfg.Logger != nil {
		l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Debug(fmt.Sprint(args...))
	}
}

func (l *asynqLogger) Info(args ...interface{}) {
	if l.cfg != nil && l.cfg.Logger != nil {
		msg := fmt.Sprint(args...)
		if strings.HasPrefix(msg, "Send signal") {
			l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Debug(msg)
			return
		}
		l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Info(msg)
	}
}

func (l *asynqLogger) Warn(args ...interface{}) {
	if l.cfg != nil && l.cfg.Logger != nil {
		l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Warn(fmt.Sprint(args...))
	}
}

func (l *asynqLogger) Error(args ...interface{}) {
	if l.cfg != nil && l.cfg.Logger != nil {
		l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Error(fmt.Sprint(args...))
	}
}

func (l *asynqLogger) Fatal(args ...interface{}) {
	if l.cfg != nil && l.cfg.Logger != nil {
		l.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Error("FATAL: " + fmt.Sprint(args...))
	}
}
