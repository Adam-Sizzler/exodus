package jobqueue

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/logger"
)

const (
	UsersUpdateUsageQueue  = "USERS_UPDATE_USERS_USAGE_QUEUE"
	UsersWatchdogQueue     = "USERS_WATCHDOG_QUEUE"
	UsersResetTrafficQueue = "USERS_RESET_USER_TRAFFIC_QUEUE"
	UsersSerialOpsQueue    = "USERS_SERIAL_OPERATIONS_QUEUE"
	SquadsActionsQueue     = "SQUADS_ACTIONS_QUEUE"
	UserEventsQueue        = "USER_EVENTS_QUEUE"

	JobUpdateUserUsage  = "update_user_usage"
	JobWatchdogReview   = "watchdog_review"
	JobResetUserTraffic = "reset_user_traffic"
	JobSerialOperation  = "serial_operation"
	JobSquadAction      = "squad_action"
	JobUserEvent        = "user_event"
)

type MiscQueuesDispatcher struct {
	processor *Processor
}

var (
	miscDispatcherMu sync.RWMutex
	miscJobs         *MiscQueuesDispatcher
)

type UserUsagePayload struct {
	NodeUUID string            `json:"nodeUuid"`
	Deltas   map[string]uint64 `json:"deltas"`
}

type WatchdogReviewPayload struct {
	ReviewType string `json:"reviewType"`
}

type ResetTrafficPayload struct {
	UserUUID string `json:"userUuid"`
}

type SquadActionPayload struct {
	SquadUUID string   `json:"squadUuid"`
	Action    string   `json:"action"`
	UserUUIDs []string `json:"userUuids"`
}

type UserEventPayload struct {
	EventType string          `json:"eventType"`
	Data      json.RawMessage `json:"data"`
}

type MiscQueueHandlers struct {
	HandleUserUsage        func(ctx context.Context, payload UserUsagePayload) error
	HandleWatchdogReview   func(ctx context.Context, payload WatchdogReviewPayload) error
	HandleResetUserTraffic func(ctx context.Context, payload ResetTrafficPayload) error
	HandleSquadAction      func(ctx context.Context, payload SquadActionPayload) error
	HandleUserEvent        func(ctx context.Context, payload UserEventPayload) error
}

func StartMiscQueues(ctx context.Context, wg *sync.WaitGroup, dbConn *sql.DB, cfg *config.BackendConfig, handlers MiscQueueHandlers) (*Processor, error) {
	client, err := NewRedisClient(cfg)
	if err != nil || client == nil {
		return nil, err
	}

	processor := NewProcessor(client, cfg)
	visibility := time.Duration(cfg.Redis.JobQueueVisibilitySeconds) * time.Second
	if visibility <= 0 {
		visibility = 5 * time.Minute
	}

	// 1. USERS_UPDATE_USERS_USAGE_QUEUE (concurrency: 5)
	if handlers.HandleUserUsage != nil {
		_ = processor.RegisterQueue(QueueOptions{
			Name:              UsersUpdateUsageQueue,
			Concurrency:       5,
			VisibilityTimeout: visibility,
		}, map[string]Handler{
			JobUpdateUserUsage: func(ctx context.Context, job Job) error {
				var payload UserUsagePayload
				if err := json.Unmarshal(job.Payload, &payload); err != nil {
					return err
				}
				return handlers.HandleUserUsage(ctx, payload)
			},
		})
	}

	// 2. USERS_WATCHDOG_QUEUE (concurrency: 1)
	if handlers.HandleWatchdogReview != nil {
		_ = processor.RegisterQueue(QueueOptions{
			Name:              UsersWatchdogQueue,
			Concurrency:       1,
			VisibilityTimeout: visibility,
		}, map[string]Handler{
			JobWatchdogReview: func(ctx context.Context, job Job) error {
				var payload WatchdogReviewPayload
				if err := json.Unmarshal(job.Payload, &payload); err != nil {
					return err
				}
				return handlers.HandleWatchdogReview(ctx, payload)
			},
		})
	}

	// 3. USERS_RESET_USER_TRAFFIC_QUEUE (concurrency: 1)
	if handlers.HandleResetUserTraffic != nil {
		_ = processor.RegisterQueue(QueueOptions{
			Name:              UsersResetTrafficQueue,
			Concurrency:       1,
			VisibilityTimeout: visibility,
		}, map[string]Handler{
			JobResetUserTraffic: func(ctx context.Context, job Job) error {
				var payload ResetTrafficPayload
				if err := json.Unmarshal(job.Payload, &payload); err != nil {
					return err
				}
				return handlers.HandleResetUserTraffic(ctx, payload)
			},
		})
	}

	// 4. SQUADS_ACTIONS_QUEUE (concurrency: 1)
	if handlers.HandleSquadAction != nil {
		_ = processor.RegisterQueue(QueueOptions{
			Name:              SquadsActionsQueue,
			Concurrency:       1,
			VisibilityTimeout: visibility,
		}, map[string]Handler{
			JobSquadAction: func(ctx context.Context, job Job) error {
				var payload SquadActionPayload
				if err := json.Unmarshal(job.Payload, &payload); err != nil {
					return err
				}
				return handlers.HandleSquadAction(ctx, payload)
			},
		})
	}

	// 5. USER_EVENTS_QUEUE (concurrency: 50)
	if handlers.HandleUserEvent != nil {
		_ = processor.RegisterQueue(QueueOptions{
			Name:              UserEventsQueue,
			Concurrency:       50,
			VisibilityTimeout: visibility,
		}, map[string]Handler{
			JobUserEvent: func(ctx context.Context, job Job) error {
				var payload UserEventPayload
				if err := json.Unmarshal(job.Payload, &payload); err != nil {
					return err
				}
				return handlers.HandleUserEvent(ctx, payload)
			},
		})
	}

	miscDispatcherMu.Lock()
	miscJobs = &MiscQueuesDispatcher{processor: processor}
	miscDispatcherMu.Unlock()

	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.RoleService(logger.RoleWorkers, logger.ServiceJobs).Info("6 misc named queues connected", "queues", "USAGE, WATCHDOG, RESET_TRAFFIC, SERIAL_OPS, SQUADS, USER_EVENTS")
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

func EnqueueUsersUsageUpdate(ctx context.Context, payload UserUsagePayload) (bool, error) {
	return enqueueMiscJob(ctx, UsersUpdateUsageQueue, JobUpdateUserUsage, payload, JobOptions{
		Attempts: 3,
	})
}

func EnqueueWatchdogReview(ctx context.Context, payload WatchdogReviewPayload) (bool, error) {
	return enqueueMiscJob(ctx, UsersWatchdogQueue, JobWatchdogReview, payload, JobOptions{
		ID:       fmt.Sprintf("watchdog:%s", payload.ReviewType),
		DedupeID: fmt.Sprintf("watchdog:%s", payload.ReviewType),
		Attempts: 1,
	})
}

func EnqueueResetUserTraffic(ctx context.Context, payload ResetTrafficPayload) (bool, error) {
	return enqueueMiscJob(ctx, UsersResetTrafficQueue, JobResetUserTraffic, payload, JobOptions{
		ID:       fmt.Sprintf("reset_traffic:%s", payload.UserUUID),
		DedupeID: fmt.Sprintf("reset_traffic:%s", payload.UserUUID),
		Attempts: 3,
	})
}

func EnqueueSquadAction(ctx context.Context, payload SquadActionPayload) (bool, error) {
	return enqueueMiscJob(ctx, SquadsActionsQueue, JobSquadAction, payload, JobOptions{
		Attempts: 3,
	})
}

func EnqueueUserEvent(ctx context.Context, payload UserEventPayload) (bool, error) {
	return enqueueMiscJob(ctx, UserEventsQueue, JobUserEvent, payload, JobOptions{
		Attempts: 5,
	})
}

func enqueueMiscJob(ctx context.Context, queue string, jobName string, payload any, options JobOptions) (bool, error) {
	miscDispatcherMu.RLock()
	dispatcher := miscJobs
	miscDispatcherMu.RUnlock()
	if dispatcher == nil || dispatcher.processor == nil {
		return false, nil
	}
	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}
	err = dispatcher.processor.Enqueue(ctx, queue, jobName, rawPayload, options)
	return err == nil, err
}
