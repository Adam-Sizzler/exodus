package notifications

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
)

const (
	ScopeUser            = "user"
	ScopeUserHWIDDevices = "user_hwid_devices"
	ScopeNode            = "node"
	ScopeService         = "service"
	ScopeErrors          = "errors"
	ScopeCRM             = "crm"

	EventUserCreated                = "user.created"
	EventUserModified               = "user.modified"
	EventUserDeleted                = "user.deleted"
	EventUserRevoked                = "user.revoked"
	EventUserDisabled               = "user.disabled"
	EventUserEnabled                = "user.enabled"
	EventUserLimited                = "user.limited"
	EventUserExpired                = "user.expired"
	EventUserTrafficReset           = "user.traffic_reset"
	EventUserFirstConnected         = "user.first_connected"
	EventUserExpiresIn72Hours       = "user.expires_in_72_hours"
	EventUserExpiresIn48Hours       = "user.expires_in_48_hours"
	EventUserExpiresIn24Hours       = "user.expires_in_24_hours"
	EventUserExpired24HoursAgo      = "user.expired_24_hours_ago"
	EventUserBandwidthThreshold     = "user.bandwidth_usage_threshold_reached"
	EventUserNotConnected           = "user.not_connected"
	EventUserHWIDDeviceAdded        = "user_hwid_devices.added"
	EventUserHWIDDeviceDeleted      = "user_hwid_devices.deleted"
	EventNodeCreated                = "node.created"
	EventNodeModified               = "node.modified"
	EventNodeDisabled               = "node.disabled"
	EventNodeEnabled                = "node.enabled"
	EventNodeDeleted                = "node.deleted"
	EventNodeConnectionLost         = "node.connection_lost"
	EventNodeConnectionRestored     = "node.connection_restored"
	EventNodeTrafficNotify          = "node.traffic_notify"
	EventServicePanelStarted        = "service.panel_started"
	EventLoginAttemptFailed         = "service.login_attempt_failed"
	EventLoginAttemptSuccess        = "service.login_attempt_success"
	EventServiceSubpageChanged      = "service.subpage_config_changed"
	EventBandwidthMaxNotification   = "errors.bandwidth_usage_threshold_reached_max_notifications"
	EventInfraBillingIn7Days        = "crm.infra_billing_node_payment_in_7_days"
	EventInfraBillingIn48Hours      = "crm.infra_billing_node_payment_in_48hrs"
	EventInfraBillingIn24Hours      = "crm.infra_billing_node_payment_in_24hrs"
	EventInfraBillingDueToday       = "crm.infra_billing_node_payment_due_today"
	EventInfraBillingOverdue24Hours = "crm.infra_billing_node_payment_overdue_24hrs"
	EventInfraBillingOverdue48Hours = "crm.infra_billing_node_payment_overdue_48hrs"
	EventInfraBillingOverdue7Days   = "crm.infra_billing_node_payment_overdue_7_days"
)

type Event struct {
	Scope     string         `json:"scope"`
	Event     string         `json:"event"`
	Timestamp string         `json:"timestamp"`
	Data      map[string]any `json:"data"`
	Meta      map[string]any `json:"meta,omitempty"`
}

type Notifier struct {
	cfg    *config.BackendConfig
	client *http.Client
}

func New(cfg *config.BackendConfig) *Notifier {
	return &Notifier{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func Emit(ctx context.Context, cfg *config.BackendConfig, event Event) {
	if cfg == nil {
		return
	}
	copied := event
	if enqueueWithGlobalDispatcher(ctx, copied) {
		return
	}
	_ = ctx
	go func() {
		sendCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		New(cfg).Send(sendCtx, copied)
	}()
}

func EmitSync(ctx context.Context, cfg *config.BackendConfig, event Event) {
	_ = New(cfg).Send(ctx, event)
}

func (n *Notifier) Enabled() bool {
	if n == nil || n.cfg == nil {
		return false
	}
	return n.telegramEnabled() || n.webhookEnabled()
}

func (n *Notifier) Send(ctx context.Context, event Event) error {
	if n == nil || !n.Enabled() {
		return nil
	}
	if strings.TrimSpace(event.Timestamp) == "" {
		event.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if event.Data == nil {
		event.Data = map[string]any{}
	}

	var errs []error
	if n.webhookEnabled() && n.cfg.Notifications.EventChannelEnabled(event.Event, "webhook") {
		if err := n.sendWebhook(ctx, event); err != nil {
			if n.cfg.Logger != nil {
				n.cfg.Logger.Warn("Webhook notification failed", "event", event.Event, "error", err)
			}
			errs = append(errs, err)
		}
	}
	if n.telegramEnabled() && n.cfg.Notifications.EventChannelEnabled(event.Event, "telegram") && !boolValue(event.Meta, "skipTelegramNotification") {
		if err := n.sendTelegram(ctx, event); err != nil {
			if n.cfg.Logger != nil {
				n.cfg.Logger.Warn("Telegram notification failed", "event", event.Event, "error", err)
			}
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (n *Notifier) telegramEnabled() bool {
	if n == nil || n.cfg == nil {
		return false
	}
	settings := n.cfg.Notifications
	return settings.TelegramEnabled && strings.TrimSpace(settings.TelegramBotToken) != ""
}

func (n *Notifier) webhookEnabled() bool {
	if n == nil || n.cfg == nil {
		return false
	}
	settings := n.cfg.Notifications
	return settings.WebhookEnabled && len(settings.WebhookURLs) > 0 && strings.TrimSpace(settings.WebhookSecret) != ""
}

func (n *Notifier) sendWebhook(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	timestamp := event.Timestamp
	signature := hmacSHA256Hex(n.cfg.Notifications.WebhookSecret, payload)

	var lastErr error
	for _, rawURL := range n.cfg.Notifications.WebhookURLs {
		url := strings.TrimSpace(rawURL)
		if url == "" {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Exodus-Signature", signature)
		req.Header.Set("X-Exodus-Timestamp", timestamp)
		req.Header.Set("User-Agent", "Exodus")

		resp, err := n.client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		_ = resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("webhook %s returned %d: %s", url, resp.StatusCode, strings.TrimSpace(string(body)))
			continue
		}
	}
	return lastErr
}

func (n *Notifier) sendTelegram(ctx context.Context, event Event) error {
	chatID, threadID := n.telegramTarget(event.Scope)
	if strings.TrimSpace(chatID) == "" {
		return nil
	}
	text := formatTelegramMessage(event)
	if strings.TrimSpace(text) == "" {
		return nil
	}

	payload := map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}
	if id, ok := parseOptionalInt(threadID); ok {
		payload["message_thread_id"] = id
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := "https://api.telegram.org/bot" + strings.TrimSpace(n.cfg.Notifications.TelegramBotToken) + "/sendMessage"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		if resp.StatusCode == http.StatusTooManyRequests {
			if retryAfter := telegramRetryAfter(respBody); retryAfter > 0 {
				return RateLimitError{RetryAfter: retryAfter, Err: fmt.Errorf("telegram returned %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))}
			}
		}
		return fmt.Errorf("telegram returned %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}

type RateLimitError struct {
	RetryAfter time.Duration
	Err        error
}

func (e RateLimitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "notification rate limited"
}

func (e RateLimitError) Unwrap() error {
	return e.Err
}

func telegramRetryAfter(body []byte) time.Duration {
	var payload struct {
		Parameters struct {
			RetryAfter int `json:"retry_after"`
		} `json:"parameters"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0
	}
	if payload.Parameters.RetryAfter <= 0 {
		return 0
	}
	return time.Duration(payload.Parameters.RetryAfter) * time.Second
}

func (n *Notifier) telegramTarget(scope string) (string, string) {
	settings := n.cfg.Notifications
	switch scope {
	case ScopeUser:
		return settings.TelegramUsersChatID, settings.TelegramUsersThreadID
	case ScopeUserHWIDDevices:
		return settings.TelegramUsersChatID, settings.TelegramUsersThreadID
	case ScopeNode:
		return settings.TelegramNodesChatID, settings.TelegramNodesThreadID
	case ScopeService, ScopeErrors:
		return settings.TelegramNodesChatID, settings.TelegramNodesThreadID
	case ScopeCRM:
		return settings.TelegramCRMChatID, settings.TelegramCRMThreadID
	default:
		return "", ""
	}
}

func hmacSHA256Hex(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func parseOptionalInt(value string) (int64, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(trimmed, 10, 64)
	return parsed, err == nil
}

func formatTelegramMessage(event Event) string {
	switch event.Scope {
	case ScopeUser:
		return formatUserMessage(event)
	case ScopeUserHWIDDevices:
		return formatUserHWIDDeviceMessage(event)
	case ScopeNode:
		return formatNodeMessage(event)
	case ScopeService, ScopeErrors:
		return formatGenericMessage(event)
	case ScopeCRM:
		return formatCRMMessage(event)
	default:
		return ""
	}
}

func formatUserMessage(event Event) string {
	username := html.EscapeString(stringValue(event.Data, "username"))
	if username == "" {
		username = html.EscapeString(stringValue(event.Data, "uuid"))
	}
	header := strings.TrimPrefix(event.Event, "user.")
	fullInfo := fmt.Sprintf(
		"<b>Username:</b> <code>%s</code>\n<b>Traffic limit:</b> <code>%s</code>\n<b>Valid until:</b> <code>%s</code>\n<b>Sub:</b> <code>%s</code>",
		username,
		html.EscapeString(formatBytes(int64Value(event.Data, "trafficLimitBytes"))),
		html.EscapeString(stringValue(event.Data, "expireAt")),
		html.EscapeString(stringValue(event.Data, "shortUuid")),
	)
	switch event.Event {
	case EventUserBandwidthThreshold:
		traffic := int64Value(event.Data, "usedTrafficBytes")
		if traffic == 0 {
			traffic = int64Value(nestedMap(event.Data, "userTraffic"), "usedTrafficBytes")
		}
		return fmt.Sprintf(
			"<b>#%s</b>\n---------\n<b>Username:</b> <code>%s</code>\n<b>Traffic:</b> <code>%s</code>\n<b>Limit:</b> <code>%s</code>\n<b>Threshold:</b> <code>%d%%</code>",
			html.EscapeString(header),
			username,
			html.EscapeString(formatBytes(traffic)),
			html.EscapeString(formatBytes(int64Value(event.Data, "trafficLimitBytes"))),
			intValue(event.Data, "lastTriggeredThreshold"),
		)
	case EventUserNotConnected:
		return fmt.Sprintf(
			"<b>#not_connected_after_%d_hours</b>\n---------\n<b>Username:</b> <code>%s</code>",
			intValue(event.Meta, "notConnectedAfterHours"),
			username,
		)
	case EventUserCreated, EventUserModified, EventUserRevoked:
		return fmt.Sprintf("<b>#%s</b>\n---------\n%s", html.EscapeString(header), fullInfo)
	case EventUserTrafficReset:
		traffic := int64Value(event.Data, "usedTrafficBytes")
		if traffic == 0 {
			traffic = int64Value(nestedMap(event.Data, "userTraffic"), "usedTrafficBytes")
		}
		return fmt.Sprintf(
			"<b>#traffic_reset</b>\n---------\n<b>Username:</b> <code>%s</code>\n<b>Traffic:</b> <code>%s</code>",
			username,
			html.EscapeString(formatBytes(traffic)),
		)
	default:
		return fmt.Sprintf("<b>#%s</b>\n---------\n<b>Username:</b> <code>%s</code>", html.EscapeString(header), username)
	}
}

func formatNodeMessage(event Event) string {
	if event.Event != EventNodeTrafficNotify {
		return fmt.Sprintf(
			"<b>#%s</b>\n---------\n<b>Name:</b> <code>%s</code>\n<b>Address:</b> <code>%s:%d</code>\n<b>Reason:</b> <code>%s</code>\n<b>Last status change:</b> <code>%s</code>",
			html.EscapeString(strings.TrimPrefix(event.Event, "node.")),
			html.EscapeString(firstNonEmptyString(stringValue(event.Data, "name"), stringValue(event.Data, "nodeName"), stringValue(event.Data, "uuid"))),
			html.EscapeString(stringValue(event.Data, "address")),
			intValue(event.Data, "port"),
			html.EscapeString(stringValue(event.Data, "lastStatusMessage")),
			html.EscapeString(stringValue(event.Data, "lastStatusChange")),
		)
	}
	return fmt.Sprintf(
		"<b>#nodeTrafficNotify</b>\n---------\n<b>Name:</b> <code>%s</code>\n<b>Address:</b> <code>%s:%d</code>\n<b>Traffic:</b> <code>%s</code> / <code>%s</code>\n<b>Traffic reset day:</b> <code>%d</code>\n<b>Percent:</b> <code>%d%%</code>",
		html.EscapeString(stringValue(event.Data, "name")),
		html.EscapeString(stringValue(event.Data, "address")),
		intValue(event.Data, "port"),
		html.EscapeString(formatBytes(int64Value(event.Data, "trafficUsedBytes"))),
		html.EscapeString(formatBytes(int64Value(event.Data, "trafficLimitBytes"))),
		intValue(event.Data, "trafficResetDay"),
		intValue(event.Data, "notifyPercent"),
	)
}

func formatUserHWIDDeviceMessage(event Event) string {
	username := html.EscapeString(stringValue(event.Data, "username"))
	if username == "" {
		username = html.EscapeString(stringValue(event.Data, "userUuid"))
	}
	return fmt.Sprintf(
		"<b>#%s</b>\n---------\n<b>User:</b> <code>%s</code>\n<b>HWID:</b> <code>%s</code>\n<b>Platform:</b> <code>%s</code>\n<b>Model:</b> <code>%s</code>",
		html.EscapeString(strings.TrimPrefix(event.Event, "user_hwid_devices.")),
		username,
		html.EscapeString(stringValue(event.Data, "hwid")),
		html.EscapeString(stringValue(event.Data, "platform")),
		html.EscapeString(stringValue(event.Data, "deviceModel")),
	)
}

func formatGenericMessage(event Event) string {
	switch event.Event {
	case EventServicePanelStarted:
		return fmt.Sprintf("<b>#panel_started</b>\n---------\nExodus <code>%s</code> is up and running.", html.EscapeString(firstNonEmptyString(stringValue(event.Data, "version"), stringValue(event.Data, "panelVersion"))))
	case EventLoginAttemptFailed, EventLoginAttemptSuccess:
		loginAttempt := nestedMap(event.Data, "loginAttempt")
		if len(loginAttempt) == 0 {
			loginAttempt = event.Data
		}
		return fmt.Sprintf(
			"<b>#%s</b>\n---------\n<b>User:</b> <code>%s</code>\n<b>IP:</b> <code>%s</code>\n<b>User agent:</b> <code>%s</code>\n<b>Description:</b> <code>%s</code>",
			html.EscapeString(strings.TrimPrefix(event.Event, "service.")),
			html.EscapeString(stringValue(loginAttempt, "username")),
			html.EscapeString(stringValue(loginAttempt, "ip")),
			html.EscapeString(stringValue(loginAttempt, "userAgent")),
			html.EscapeString(stringValue(loginAttempt, "description")),
		)
	case EventServiceSubpageChanged:
		subpageConfig := nestedMap(event.Data, "subpageConfig")
		return fmt.Sprintf(
			"<b>#subpage_config_changed</b>\n---------\n<b>Action:</b> <code>%s</code>\n<b>UUID:</b> <code>%s</code>",
			html.EscapeString(stringValue(subpageConfig, "action")),
			html.EscapeString(stringValue(subpageConfig, "uuid")),
		)
	case EventBandwidthMaxNotification:
		return fmt.Sprintf(
			"<b>#bandwidth_usage_threshold_reached_max_notifications</b>\n---------\n<b>Batch:</b> <code>%d</code>\n<b>Total processed:</b> <code>%d</code>",
			intValue(event.Data, "batchSize"),
			intValue(event.Data, "totalProcessed"),
		)
	}
	title := event.Event
	if title == "" {
		title = event.Scope
	}
	main := firstNonEmptyString(
		stringValue(event.Data, "username"),
		stringValue(event.Data, "name"),
		stringValue(event.Data, "nodeName"),
		stringValue(event.Data, "uuid"),
		stringValue(event.Data, "userUuid"),
	)
	if main == "" {
		return fmt.Sprintf("<b>#%s</b>", html.EscapeString(title))
	}
	return fmt.Sprintf("<b>#%s</b>\n---------\n<b>Target:</b> <code>%s</code>", html.EscapeString(title), html.EscapeString(main))
}

func formatCRMMessage(event Event) string {
	return fmt.Sprintf(
		"<b>#%s</b>\n---------\n<b>Provider:</b> <code>%s</code>\n<b>Node:</b> <code>%s</code>\n<b>Due date:</b> <code>%s</code>\n<a href=\"%s\">Open Provider Panel</a>",
		html.EscapeString(strings.TrimPrefix(event.Event, "crm.")),
		html.EscapeString(stringValue(event.Data, "providerName")),
		html.EscapeString(stringValue(event.Data, "nodeName")),
		html.EscapeString(stringValue(event.Data, "nextBillingAt")),
		html.EscapeString(stringValue(event.Data, "loginUrl")),
	)
}

func nestedMap(data map[string]any, key string) map[string]any {
	if data == nil {
		return nil
	}
	raw, exists := data[key]
	if !exists || raw == nil {
		return nil
	}
	if typed, ok := raw.(map[string]any); ok {
		return typed
	}
	return nil
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func stringValue(data map[string]any, key string) string {
	if data == nil {
		return ""
	}
	raw, exists := data[key]
	if !exists || raw == nil {
		return ""
	}
	switch value := raw.(type) {
	case string:
		return value
	case fmt.Stringer:
		return value.String()
	default:
		return fmt.Sprint(value)
	}
}

func intValue(data map[string]any, key string) int {
	return int(int64Value(data, key))
}

func int64Value(data map[string]any, key string) int64 {
	if data == nil {
		return 0
	}
	switch value := data[key].(type) {
	case int:
		return int64(value)
	case int64:
		return value
	case float64:
		return int64(value)
	case json.Number:
		parsed, _ := value.Int64()
		return parsed
	default:
		return 0
	}
}

func boolValue(data map[string]any, key string) bool {
	if data == nil {
		return false
	}
	value, _ := data[key].(bool)
	return value
}

func formatBytes(value int64) string {
	if value < 0 {
		value = 0
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	floatValue := float64(value)
	unit := 0
	for floatValue >= 1024 && unit < len(units)-1 {
		floatValue /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d %s", value, units[unit])
	}
	return fmt.Sprintf("%.2f %s", floatValue, units[unit])
}
