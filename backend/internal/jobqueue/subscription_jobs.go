package jobqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/logger"
)

const (
	subscriptionQueueName     = "subscription_requests"
	jobUpdateUserSubscription = "update_user_subscription"
	jobAddSubscriptionRecord  = "add_subscription_request_record"
	jobUpsertHwidDevice       = "upsert_hwid_device"
)

type subscriptionDispatcher struct {
	processor *Processor
}

type UpdateUserSubscriptionPayload struct {
	UserUUID  string `json:"userUuid"`
	UserAgent string `json:"userAgent"`
}

type AddSubscriptionRequestRecordPayload struct {
	UserID    int64  `json:"userId"`
	RequestIP string `json:"requestIp"`
	UserAgent string `json:"userAgent"`
}

type UpsertHwidDevicePayload struct {
	UserID      int64   `json:"userId"`
	Hwid        string  `json:"hwid"`
	Platform    *string `json:"platform,omitempty"`
	OsVersion   *string `json:"osVersion,omitempty"`
	DeviceModel *string `json:"deviceModel,omitempty"`
	UserAgent   *string `json:"userAgent,omitempty"`
	RequestIP   *string `json:"requestIp,omitempty"`
}

var (
	subscriptionDispatcherMu sync.RWMutex
	subscriptionJobs         *subscriptionDispatcher
)

func StartSubscriptionQueues(ctx context.Context, wg *sync.WaitGroup, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (*Processor, error) {
	client, err := NewRedisClient(cfg)
	if err != nil || client == nil {
		return nil, err
	}

	processor := NewProcessor(client, cfg)
	visibility := time.Duration(cfg.Redis.JobQueueVisibilitySeconds) * time.Second
	if visibility <= 0 {
		visibility = 5 * time.Minute
	}

	if err := processor.RegisterQueue(QueueOptions{
		Name:              subscriptionQueueName,
		Concurrency:       cfg.Redis.SubscriptionQueueConcurrency,
		VisibilityTimeout: visibility,
		SchedulerInterval: time.Second,
		Retention:         500,
	}, map[string]Handler{
		jobUpdateUserSubscription: func(ctx context.Context, job Job) error {
			var payload UpdateUserSubscriptionPayload
			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				return err
			}
			return updateUserSubscription(ctx, manager, payload)
		},
		jobAddSubscriptionRecord: func(ctx context.Context, job Job) error {
			var payload AddSubscriptionRequestRecordPayload
			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				return err
			}
			return addSubscriptionRequestRecord(ctx, manager, payload)
		},
		jobUpsertHwidDevice: func(ctx context.Context, job Job) error {
			var payload UpsertHwidDevicePayload
			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				return err
			}
			return upsertHwidDevice(ctx, manager, payload)
		},
	}); err != nil {
		_ = client.Close()
		return nil, err
	}

	subscriptionDispatcherMu.Lock()
	subscriptionJobs = &subscriptionDispatcher{processor: processor}
	subscriptionDispatcherMu.Unlock()

	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceUsersQueue).Info("1 queues connected", "queue", subscriptionQueueName, "concurrency", cfg.Redis.SubscriptionQueueConcurrency)
	}

	processor.Start(ctx, wg)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = processor.Close()
	}()

	return processor, nil
}

func EnqueueUpdateUserSubscription(ctx context.Context, payload UpdateUserSubscriptionPayload) (bool, error) {
	return enqueueSubscriptionJob(ctx, jobUpdateUserSubscription, payload, JobOptions{
		ID:       fmt.Sprintf("%s:USS", payload.UserUUID),
		DedupeID: fmt.Sprintf("%s:USS", payload.UserUUID),
		Attempts: 3,
		Backoff:  time.Second,
	})
}

func EnqueueAddSubscriptionRequestRecord(ctx context.Context, payload AddSubscriptionRequestRecordPayload) (bool, error) {
	return enqueueSubscriptionJob(ctx, jobAddSubscriptionRecord, payload, JobOptions{
		ID:       fmt.Sprintf("%d:AR", payload.UserID),
		DedupeID: fmt.Sprintf("%d:AR", payload.UserID),
		Attempts: 3,
		Backoff:  time.Second,
	})
}

func EnqueueUpsertHwidDevice(ctx context.Context, payload UpsertHwidDevicePayload) (bool, error) {
	if payload.Hwid == "" || payload.UserID <= 0 {
		return false, nil
	}
	jobID := fmt.Sprintf("%d:%s:CAUHD", payload.UserID, payload.Hwid)
	return enqueueSubscriptionJob(ctx, jobUpsertHwidDevice, payload, JobOptions{
		ID:       jobID,
		DedupeID: jobID,
		Attempts: 3,
		Backoff:  time.Second,
	})
}

func enqueueSubscriptionJob(ctx context.Context, jobName string, payload any, options JobOptions) (bool, error) {
	subscriptionDispatcherMu.RLock()
	dispatcher := subscriptionJobs
	subscriptionDispatcherMu.RUnlock()
	if dispatcher == nil || dispatcher.processor == nil {
		return false, nil
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}
	err = dispatcher.processor.Enqueue(ctx, subscriptionQueueName, jobName, rawPayload, options)
	return err == nil, err
}

func updateUserSubscription(ctx context.Context, manager *dbmanager.DatabaseManager, payload UpdateUserSubscriptionPayload) error {
	if payload.UserUUID == "" {
		return nil
	}
	return nil
}

func addSubscriptionRequestRecord(ctx context.Context, manager *dbmanager.DatabaseManager, payload AddSubscriptionRequestRecordPayload) error {
	if payload.UserID <= 0 {
		return nil
	}
	return manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		if _, err := db.ExecContext(ctx, `
			INSERT INTO user_subscription_request_history (user_id, request_ip, user_agent)
			VALUES (?, ?, ?)
		`, payload.UserID, payload.RequestIP, payload.UserAgent); err != nil {
			return err
		}

		_, err := db.ExecContext(ctx, `
			DELETE FROM user_subscription_request_history
			WHERE user_id = ?
			  AND id NOT IN (
				  SELECT id
				  FROM user_subscription_request_history
				  WHERE user_id = ?
				  ORDER BY request_at DESC, id DESC
				  LIMIT 24
			  )
		`, payload.UserID, payload.UserID)
		return err
	})
}

func upsertHwidDevice(ctx context.Context, manager *dbmanager.DatabaseManager, payload UpsertHwidDevicePayload) error {
	if payload.UserID <= 0 || payload.Hwid == "" {
		return nil
	}
	payload.Platform = lowerStringPtr(payload.Platform)
	return manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT (hwid, user_id)
			DO UPDATE SET
				platform = EXCLUDED.platform,
				os_version = EXCLUDED.os_version,
				device_model = EXCLUDED.device_model,
				user_agent = EXCLUDED.user_agent,
				request_ip = COALESCE(EXCLUDED.request_ip, hwid_user_devices.request_ip),
				updated_at = now()
		`, payload.Hwid, payload.UserID, payload.Platform, payload.OsVersion, payload.DeviceModel, payload.UserAgent, payload.RequestIP)
		return err
	})
}

func lowerStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	lowered := strings.ToLower(strings.TrimSpace(*value))
	if lowered == "" {
		return nil
	}
	return &lowered
}
