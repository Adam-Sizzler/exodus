package server

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"exodus-node/api"
	"exodus-node/config"
)

const (
	coreProcessName         = "singbox"
	coreHealthcheckAttempts = 10
	coreHealthcheckInterval = 2 * time.Second
	supervisorCallTimeout   = 25 * time.Second
)

type supervisorProcessInfo struct {
	Name        string
	StateName   string
	Description string
}

type supervisorClient struct {
	socketPath string
	username   string
	password   string
	httpClient *http.Client
	logger     interface {
		Trace(string, ...any)
		Debug(string, ...any)
		Warn(string, ...any)
	}
}

type coreLifecycleResult struct {
	ProcessBefore string
	ProcessAfter  string
	Started       bool
	Ready         bool
	Error         string
}

func (r coreLifecycleResult) failed() bool {
	return strings.TrimSpace(r.Error) != ""
}

func restartCoreProcessLifecycle(ctx context.Context, cfg *config.NodeConfig, apiService *api.Service) coreLifecycleResult {
	result := coreLifecycleResult{}
	if cfg == nil || cfg.Logger == nil {
		result.Error = "node config/logger is nil"
		return result
	}
	log := cfg.LoggerFor("SingboxService")
	if apiService != nil {
		apiService.MarkCoreOffline()
	}

	log.Debug("Starting managed core lifecycle", "process", coreProcessName, "config", config.FixedSingboxConfigPath)
	// lifecycle: the freshly received config is already the current
	// managed core config. Do not pre-gate supervisor start with a local check here.
	// If the config is invalid, supervisor/core start returns the business error and
	// the node control-plane remains alive for the next corrective deploy.
	log.Debug("Starting core through supervisor without pre-validation gate", "config", config.FixedSingboxConfigPath)

	supervisor, err := newSupervisorClient(cfg)
	if err != nil {
		result.Error = fmt.Sprintf("create supervisor client: %v", err)
		log.Error("Failed to start Sing-box: " + err.Error())
		return result
	}

	before, err := supervisor.GetProcessInfo(ctx, coreProcessName)
	if err != nil {
		result.Error = fmt.Sprintf("get supervisor process info: %v", err)
		log.Error("Failed to start Sing-box: " + err.Error())
		return result
	}
	result.ProcessBefore = before.StateName
	log.Debug("Core supervisor state before start", "state", before.StateName, "description", before.Description)

	if shouldStopBeforeStart(before.StateName) {
		log.Debug("Stopping existing core process before managed start", "state", before.StateName)
		if err := supervisor.StopProcess(ctx, coreProcessName, true); err != nil {
			if isSupervisorAlreadyStoppedError(err) {
				log.Debug("Core process was already stopped", "error", err)
			} else {
				result.Error = fmt.Sprintf("stop %s: %v", coreProcessName, err)
				log.Error("Failed to stop Sing-box Process: " + err.Error())
				return result
			}
		} else {
			log.Debug("Core process stopped")
		}
	} else {
		log.Debug("Core stop skipped before start", "state", before.StateName)
	}

	log.Debug("Starting core process through supervisor", "process", coreProcessName)
	if err := supervisor.StartProcess(ctx, coreProcessName, true); err != nil {
		if isSupervisorAlreadyStartedError(err) {
			log.Debug("Core process was already started", "error", err)
		} else {
			result.Error = fmt.Sprintf("start %s: %v", coreProcessName, err)
			log.Error("Failed to start Sing-box: " + err.Error())
			if after, infoErr := supervisor.GetProcessInfo(ctx, coreProcessName); infoErr == nil {
				result.ProcessAfter = after.StateName
			}
			return result
		}
	}
	result.Started = true
	log.Debug("Core process start command accepted")

	after, err := supervisor.GetProcessInfo(ctx, coreProcessName)
	if err != nil {
		result.Error = fmt.Sprintf("get supervisor process info after start: %v", err)
		log.Error("Failed to start Sing-box: " + err.Error())
		return result
	}
	result.ProcessAfter = after.StateName
	log.Debug("Core supervisor state after start", "state", after.StateName, "description", after.Description)

	if !isSupervisorRunningState(after.StateName) && !strings.EqualFold(after.StateName, "STARTING") {
		result.Error = fmt.Sprintf("core supervisor state is %s: %s", after.StateName, after.Description)
		log.Error(renderCoreFailedMessage(after.StateName, result.Error))
		return result
	}

	if err := waitForCoreAPIReady(ctx, cfg, apiService); err != nil {
		result.Error = fmt.Sprintf("core API healthcheck failed: %v", err)
		if apiService != nil {
			apiService.MarkCoreOffline()
		}
		log.Error("Failed to start Sing-box: "+err.Error(), "state", after.StateName)
		return result
	}

	result.Ready = true
	if apiService != nil {
		apiService.MarkCoreOnline()
	}
	log.Log(renderCoreStartedMessage(result.ProcessAfter))
	return result
}

func newSupervisorClient(cfg *config.NodeConfig) (*supervisorClient, error) {
	socketPath := strings.TrimSpace(os.Getenv("SUPERVISORD_SOCKET_PATH"))
	if socketPath == "" {
		return nil, fmt.Errorf("SUPERVISORD_SOCKET_PATH is empty")
	}
	client := &supervisorClient{
		socketPath: socketPath,
		username:   strings.TrimSpace(os.Getenv("SUPERVISORD_USER")),
		password:   strings.TrimSpace(os.Getenv("SUPERVISORD_PASSWORD")),
		logger:     cfg.LoggerFor("SingboxService"),
	}
	client.httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", client.socketPath)
			},
		},
		Timeout: supervisorCallTimeout,
	}
	return client, nil
}

func (c *supervisorClient) GetProcessInfo(ctx context.Context, name string) (supervisorProcessInfo, error) {
	body, err := c.call(ctx, "supervisor.getProcessInfo", xmlRPCStringParam(name))
	if err != nil {
		return supervisorProcessInfo{}, err
	}
	return supervisorProcessInfo{
		Name:        firstNonEmptyXMLMember(body, "name"),
		StateName:   strings.ToUpper(firstNonEmptyXMLMember(body, "statename")),
		Description: firstNonEmptyXMLMember(body, "description"),
	}, nil
}

func (c *supervisorClient) StartProcess(ctx context.Context, name string, wait bool) error {
	_, err := c.call(ctx, "supervisor.startProcess", xmlRPCStringParam(name), xmlRPCBoolParam(wait))
	return err
}

func (c *supervisorClient) StopProcess(ctx context.Context, name string, wait bool) error {
	_, err := c.call(ctx, "supervisor.stopProcess", xmlRPCStringParam(name), xmlRPCBoolParam(wait))
	return err
}

func (c *supervisorClient) call(ctx context.Context, method string, params ...string) (string, error) {
	if c.logger != nil {
		c.logger.Trace("Calling supervisor XML-RPC", "method", method, "socket", c.socketPath)
	}

	ctx, cancel := context.WithTimeout(ctx, supervisorCallTimeout)
	defer cancel()

	var payload strings.Builder
	payload.WriteString(`<?xml version="1.0"?>`)
	payload.WriteString(`<methodCall><methodName>`)
	payload.WriteString(html.EscapeString(method))
	payload.WriteString(`</methodName><params>`)
	for _, param := range params {
		payload.WriteString(`<param><value>`)
		payload.WriteString(param)
		payload.WriteString(`</value></param>`)
	}
	payload.WriteString(`</params></methodCall>`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://supervisor/RPC2", bytes.NewBufferString(payload.String()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/xml")
	if c.username != "" || c.password != "" {
		token := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		req.Header.Set("Authorization", "Basic "+token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	body := string(data)
	if readErr != nil {
		return body, readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return body, fmt.Errorf("supervisor XML-RPC HTTP %d: %s", resp.StatusCode, strings.TrimSpace(body))
	}
	if fault := parseXMLRPCFault(body); fault != "" {
		return body, fmt.Errorf("supervisor XML-RPC fault: %s", fault)
	}
	return body, nil
}

func xmlRPCStringParam(value string) string {
	return `<string>` + html.EscapeString(value) + `</string>`
}

func xmlRPCBoolParam(value bool) string {
	if value {
		return `<boolean>1</boolean>`
	}
	return `<boolean>0</boolean>`
}

func parseXMLRPCFault(body string) string {
	if !strings.Contains(body, "<fault>") {
		return ""
	}
	if fault := firstNonEmptyXMLMember(body, "faultString"); fault != "" {
		return fault
	}
	return strings.TrimSpace(body)
}

func firstNonEmptyXMLMember(body, key string) string {
	patterns := []string{
		`(?s)<member>\s*<name>` + regexp.QuoteMeta(key) + `</name>\s*<value>\s*<string>(.*?)</string>`,
		`(?s)<member>\s*<name>` + regexp.QuoteMeta(key) + `</name>\s*<value>\s*<int>(.*?)</int>`,
		`(?s)<member>\s*<name>` + regexp.QuoteMeta(key) + `</name>\s*<value>\s*<i4>(.*?)</i4>`,
		`(?s)<member>\s*<name>` + regexp.QuoteMeta(key) + `</name>\s*<value>\s*<boolean>(.*?)</boolean>`,
	}
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(body)
		if len(match) == 2 {
			return strings.TrimSpace(html.UnescapeString(match[1]))
		}
	}
	return ""
}

func shouldStopBeforeStart(stateName string) bool {
	switch strings.ToUpper(strings.TrimSpace(stateName)) {
	case "RUNNING", "STARTING", "BACKOFF", "STOPPING":
		return true
	default:
		return false
	}
}

func isSupervisorRunningState(stateName string) bool {
	return strings.EqualFold(strings.TrimSpace(stateName), "RUNNING")
}

func isSupervisorAlreadyStoppedError(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "not running") || strings.Contains(text, "already stopped")
}

func isSupervisorAlreadyStartedError(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(err.Error())
	return strings.Contains(text, "already started") || strings.Contains(text, "already running")
}

func waitForCoreAPIReady(ctx context.Context, cfg *config.NodeConfig, apiService *api.Service) error {
	log := cfg.LoggerFor("SingboxService")
	if apiService == nil {
		return fmt.Errorf("core API service is nil")
	}

	var lastErr error
	for attempt := 1; attempt <= coreHealthcheckAttempts; attempt++ {
		checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := apiService.CheckCoreReady(checkCtx)
		cancel()
		if err == nil {
			log.Debug("Core API healthcheck passed", "attempt", attempt)
			return nil
		}
		lastErr = err
		log.Debug("Get Sing-box internal status attempt failed", "attempt", attempt, "retries_left", coreHealthcheckAttempts-attempt, "error", err)

		if attempt == coreHealthcheckAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(coreHealthcheckInterval):
		}
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("unknown core API healthcheck error")
	}
	return lastErr
}
