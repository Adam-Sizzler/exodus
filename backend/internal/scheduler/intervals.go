package scheduler

import "time"

// Job intervals and durations across background services and monitors.
const (
	// Metrics & export
	MetricExportMetricsInterval   = 15 * time.Second
	MetricSyncMetricsInterval     = 6 * time.Hour
	ExportNodeConnectionsInterval = 5 * time.Minute
	RecordUserUsageInterval       = 15 * time.Second

	// Node monitoring & reconnection
	NodeHealthCheckInterval  = 10 * time.Second
	NodeReconnectInterval    = 10 * time.Second
	SubNodeReconnectInterval = 10 * time.Second
	RecordNodeUsageInterval  = 30 * time.Second
	ReviewNodesInterval      = 1 * time.Hour

	// User review watchdog
	FindExpiredUsersInterval              = 30 * time.Second
	FindExceededTrafficUsageUsersInterval = 45 * time.Second

	// User notifications
	ExpireNotificationsInterval            = 1 * time.Minute
	BandwidthUsageNotificationsInterval    = 5 * time.Minute
	NotConnectedUsersNotificationsInterval = 10 * time.Minute

	// Service & maintenance
	ResetNodeTrafficInterval = 24 * time.Hour

	// Node & subscription streaming heartbeats
	StreamIdleTimeout   = 45 * time.Second
	StreamWatchInterval = 10 * time.Second
)

// Standard cron expressions for scheduled calendar jobs.
const (
	CronReviewNodesInterval                        = "0 * * * *"   // every hour
	CronResetNodeTrafficDay1AM                     = "0 1 * * *"   // every day at 01:00
	CronResetUserTrafficDaily                      = "5 0 * * *"   // every day at 00:05
	CronResetUserTrafficMonthlyRolling             = "10 0 * * *"  // every day at 00:10
	CronResetUserTrafficWeekly                     = "15 0 * * 1"  // every Monday at 00:15
	CronResetUserTrafficMonthly                    = "20 0 1 * *"  // 1st of month at 00:20
	CronExpireNotifications                        = "* * * * *"   // every minute
	CronBandwidthUsageNotifications                = "*/5 * * * *" // every 5 minutes
	CronNotConnectedUsersNotifications             = "*/10 * * * *"// every 10 minutes
	CronServiceCleanOldUsageRecords                = "30 0 * * 1"  // every Monday at 00:30
	CronServiceVacuumTables                        = "45 0 * * 1"  // every Monday at 00:45
	CronCRMInfraBillingNodesNotifications          = "0 17 * * *"  // every day at 17:00
	CronSRSListsAvailabilityCheck                  = "0 */12 * * *"// every 12 hours
)
