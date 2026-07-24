package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"
)

const (
	notificationsQueueName = "notifications"
	webhookJobName         = "sendWebhook"
	telegramJobName        = "sendTelegram"
)

type Worker struct {
	processor *jobqueue.Processor
	cfg       *config.BackendConfig
	notifier  *Notifier
}

func NewWorker(cfg *config.BackendConfig) (*Worker, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	
	client, err := jobqueue.NewRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	
	processor := jobqueue.NewProcessor(client, cfg)
	worker := &Worker{
		processor: processor,
		cfg:       cfg,
		notifier:  New(cfg),
	}
	
	err = processor.RegisterQueue(jobqueue.QueueOptions{
		Name:              notificationsQueueName,
		Concurrency:       5,
		VisibilityTimeout: 10 * time.Minute,
		Retention:         2000,
	}, map[string]jobqueue.Handler{
		webhookJobName:  worker.handleWebhook,
		telegramJobName: worker.handleTelegram,
	})
	
	if err != nil {
		processor.Close()
		return nil, err
	}
	
	return worker, nil
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
	if w.cfg != nil && w.cfg.Logger != nil {
		w.cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Info("Starting notifications queue worker")
	}
	w.processor.Start(ctx, wg)
}

func (w *Worker) handleWebhook(ctx context.Context, job jobqueue.Job) error {
	var event Event
	if err := json.Unmarshal(job.Payload, &event); err != nil {
		return err
	}
	return w.notifier.sendWebhook(ctx, event)
}

func (w *Worker) handleTelegram(ctx context.Context, job jobqueue.Job) error {
	var event Event
	if err := json.Unmarshal(job.Payload, &event); err != nil {
		return err
	}
	return w.notifier.sendTelegram(ctx, event)
}

func (w *Worker) EnqueueWebhook(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return w.processor.Enqueue(ctx, notificationsQueueName, webhookJobName, payload, jobqueue.JobOptions{
		Attempts: 8,
		Backoff:  time.Second,
	})
}

func (w *Worker) EnqueueTelegram(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return w.processor.Enqueue(ctx, notificationsQueueName, telegramJobName, payload, jobqueue.JobOptions{
		Attempts: 8,
		Backoff:  time.Second,
	})
}
