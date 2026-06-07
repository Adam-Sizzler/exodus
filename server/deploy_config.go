package server

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"exodus-node/config"

	"github.com/iancoleman/orderedmap"
)

// DeployConfigTaskPayload is JSON payload for SubmitTask(operation=deploy_config).
// It accepts either "config" or "singbox_config" with raw sing-box JSON.
type DeployConfigTaskPayload struct {
	Config            json.RawMessage       `json:"config"`
	SingboxConfig     json.RawMessage       `json:"singbox_config"`
	Listen            string                `json:"listen"`
	Restart           *bool                 `json:"restart"`
	ForceRestart      *bool                 `json:"force_restart"`
	ForceRestartCamel *bool                 `json:"forceRestart"`
	Stats             DeployConfigTaskStats `json:"stats"`
	SRSLists          []SRSListItem         `json:"srs_lists"`
	Modules           DeployModulesPayload  `json:"modules"`
}

type DeployConfigTaskStats struct {
	Enabled   *bool    `json:"enabled"`
	Inbounds  []string `json:"inbounds"`
	Outbounds []string `json:"outbounds"`
	Users     []string `json:"users"`
}

type DeployModulesPayload struct {
	HaproxyEnabled bool                    `json:"haproxy_enabled"`
	HaproxyUsers   []HaproxyUserEntry      `json:"haproxy_users"`
	IngressFilter  NftIngressFilterPayload `json:"ingress_filter"`
	EgressFilter   NftEgressFilterPayload  `json:"egress_filter"`
}

type HaproxyUserEntry struct {
	Username       string `json:"username"`
	VLESSUUID      string `json:"vless_uuid"`
	TrojanPassword string `json:"trojan_password"`
}

type DeploySummary struct {
	ConfigPath            string
	Listen                string
	Inbounds              int
	Outbounds             int
	Users                 int
	Restarted             bool
	ForceRestart          bool
	ConfigChanged         bool
	HaproxyUsersChanged   bool
	SRSDownloadedOnDeploy bool
	ReloadError           string
	CoreProcessBefore     string
	CoreProcessAfter      string
	CoreReady             bool
	CoreConfigValid       bool
}

// DeployConfig applies sing-box config, injects experimental.v2ray_api stats block and starts/restarts the managed core if requested.
func (s *NodeServer) DeployConfig(ctx context.Context, task DeployConfigTaskPayload) (DeploySummary, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	rawConfig := task.Config
	if len(rawConfig) == 0 {
		rawConfig = task.SingboxConfig
	}
	if len(rawConfig) == 0 {
		return DeploySummary{}, fmt.Errorf("deploy payload does not contain sing-box config")
	}
	s.Cfg.Logger.Info("DeployConfig started", "config_bytes", len(rawConfig))

	configPath := config.FixedSingboxConfigPath

	listen := task.Listen
	if listen == "" {
		listen = fmt.Sprintf("%s:%d", config.FixedCoreAPIAddress, config.FixedCoreAPIGRPCPort)
	}

	enabled := true
	if task.Stats.Enabled != nil {
		enabled = *task.Stats.Enabled
	}

	// SRS files are not synced as part of deploy_config.
	// They are refreshed only on node startup and via explicit sync_srs_lists task.
	srsDownloadedOnDeploy := false

	finalConfig, summary, err := BuildSingboxConfigWithV2RayAPI(rawConfig, BuildOptions{
		Listen:          listen,
		Enabled:         enabled,
		ExplicitUsers:   task.Stats.Users,
		ExplicitInTags:  task.Stats.Inbounds,
		ExplicitOutTags: task.Stats.Outbounds,
	})
	if err != nil {
		return DeploySummary{}, fmt.Errorf("build sing-box config: %w", err)
	}
	s.Cfg.Logger.Debug("Built sing-box config with v2ray_api", "listen", listen, "inbounds", len(summary.Inbounds), "outbounds", len(summary.Outbounds), "users", len(summary.Users))

	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		return DeploySummary{}, fmt.Errorf("create config dir: %w", err)
	}

	finalConfigHash := sha256Hex(finalConfig)
	configChanged := true
	if currentConfig, err := os.ReadFile(configPath); err == nil {
		configChanged = sha256Hex(currentConfig) != finalConfigHash
	} else if !os.IsNotExist(err) {
		return DeploySummary{}, fmt.Errorf("read current config: %w", err)
	}

	if configChanged {
		if err := os.WriteFile(configPath, finalConfig, 0o644); err != nil {
			return DeploySummary{}, fmt.Errorf("write config: %w", err)
		}
		s.Cfg.Logger.Info("Sing-box config written", "path", configPath, "bytes", len(finalConfig))
	} else {
		s.Cfg.Logger.Debug("Sing-box config unchanged, write skipped", "path", configPath)
	}

	haproxyUsersChanged, err := applyHaproxyModule(task.Modules)
	if err != nil {
		return DeploySummary{}, err
	}
	if err := applyNftablesModule(task.Modules); err != nil {
		return DeploySummary{}, err
	}

	if haproxyUsersChanged {
		reloadResult := reloadHaproxyUsers()
		switch {
		case reloadResult.Reloaded:
			s.Cfg.Logger.Info("HAProxy users cache reloaded", "socket", haproxyRuntimeSocketPath, "result", reloadResult.Output)
		case reloadResult.Skipped:
			s.Cfg.Logger.Debug("HAProxy users reload skipped", "socket", haproxyRuntimeSocketPath, "warning", reloadResult.Warning)
		default:
			s.Cfg.Logger.Warn("HAProxy users reload failed", "socket", haproxyRuntimeSocketPath, "warning", reloadResult.Warning)
		}
	} else {
		s.Cfg.Logger.Debug("HAProxy users unchanged, runtime reload skipped")
	}

	shouldRestart := task.restartRequested()
	forceRestart := task.forceRestartRequested()

	restarted := false
	reloadError := ""
	coreProcessBefore := ""
	coreProcessAfter := ""
	coreReady := false
	coreConfigValid := true
	if shouldRestart {
		if s.shouldRunCoreLifecycle(ctx, task, configChanged) {
			if forceRestart && !configChanged {
				s.Cfg.Logger.Warn("Core force restart requested by deploy payload")
			} else {
				s.Cfg.Logger.Info("Core managed lifecycle requested by deploy payload")
			}

			lifecycle := restartCoreProcessLifecycle(ctx, s.Cfg, s.apiService)
			coreProcessBefore = lifecycle.ProcessBefore
			coreProcessAfter = lifecycle.ProcessAfter
			coreReady = lifecycle.Ready
			coreConfigValid = lifecycle.ConfigValid
			if lifecycle.failed() {
				reloadError = lifecycle.Error
				s.Cfg.Logger.Error(
					"Core managed lifecycle failed after config deploy; keeping node control plane alive",
					"error", lifecycle.Error,
					"process_before", lifecycle.ProcessBefore,
					"process_after", lifecycle.ProcessAfter,
					"config_valid", lifecycle.ConfigValid,
				)
			} else {
				restarted = lifecycle.Started
			}
		} else {
			s.Cfg.Logger.Info("Core lifecycle skipped: config unchanged and core API healthy")
			coreReady = true
		}
	} else {
		s.Cfg.Logger.Info("Core lifecycle skipped by deploy payload")
	}

	return DeploySummary{
		ConfigPath:            configPath,
		Listen:                listen,
		Inbounds:              len(summary.Inbounds),
		Outbounds:             len(summary.Outbounds),
		Users:                 len(summary.Users),
		Restarted:             restarted,
		ForceRestart:          forceRestart,
		ConfigChanged:         configChanged,
		HaproxyUsersChanged:   haproxyUsersChanged,
		SRSDownloadedOnDeploy: srsDownloadedOnDeploy,
		ReloadError:           reloadError,
		CoreProcessBefore:     coreProcessBefore,
		CoreProcessAfter:      coreProcessAfter,
		CoreReady:             coreReady,
		CoreConfigValid:       coreConfigValid,
	}, nil
}

func (task DeployConfigTaskPayload) restartRequested() bool {
	if task.forceRestartRequested() {
		return true
	}
	return task.Restart != nil && *task.Restart
}

func (task DeployConfigTaskPayload) forceRestartRequested() bool {
	if task.ForceRestart != nil {
		return *task.ForceRestart
	}
	return task.ForceRestartCamel != nil && *task.ForceRestartCamel
}

func shouldReloadCore(task DeployConfigTaskPayload, configChanged bool) bool {
	if !task.restartRequested() {
		return false
	}
	return configChanged || task.forceRestartRequested()
}

func (s *NodeServer) shouldRunCoreLifecycle(ctx context.Context, task DeployConfigTaskPayload, configChanged bool) bool {
	if !task.restartRequested() {
		return false
	}
	if configChanged || task.forceRestartRequested() {
		return true
	}
	if s == nil || s.apiService == nil {
		return true
	}
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.apiService.CheckCoreReady(checkCtx); err != nil {
		s.Cfg.Logger.Warn("Core lifecycle required because core API is not healthy even though config is unchanged", "error", err)
		return true
	}
	return false
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

type BuildOptions struct {
	Listen          string
	Enabled         bool
	ExplicitInTags  []string
	ExplicitOutTags []string
	ExplicitUsers   []string
}

type buildSummary struct {
	Inbounds  []string
	Outbounds []string
	Users     []string
}

func BuildSingboxConfigWithV2RayAPI(rawConfig json.RawMessage, opt BuildOptions) ([]byte, buildSummary, error) {
	cfg := orderedmap.New()
	if err := json.Unmarshal(rawConfig, cfg); err != nil {
		return nil, buildSummary{}, fmt.Errorf("invalid sing-box JSON: %w", err)
	}
	normalizeLocalRuleSetPathsOrdered(cfg, filepath.Dir(config.FixedSingboxConfigPath))

	inboundsValue, _ := cfg.Get("inbounds")
	outboundsValue, _ := cfg.Get("outbounds")

	inboundTags := dedupeKeepOrder(opt.ExplicitInTags)
	if len(inboundTags) == 0 {
		inboundTags = extractTags(inboundsValue)
	}

	outboundTags := dedupeKeepOrder(opt.ExplicitOutTags)
	if len(outboundTags) == 0 {
		outboundTags = extractTags(outboundsValue)
	}

	users := dedupeKeepOrder(opt.ExplicitUsers)
	if len(users) == 0 {
		users = extractUsers(inboundsValue)
	}

	experimental, _ := orderedMapByKey(cfg, "experimental")
	cacheFile, _ := orderedMapByKey(&experimental, "cache_file")
	cacheFile.Set("enabled", true)
	experimental.Set("cache_file", cacheFile)

	stats := orderedmap.New()
	stats.Set("enabled", opt.Enabled)
	stats.Set("inbounds", nonNilSlice(inboundTags))
	stats.Set("outbounds", nonNilSlice(outboundTags))
	stats.Set("users", nonNilSlice(users))

	v2rayAPI := orderedmap.New()
	v2rayAPI.Set("listen", opt.Listen)
	v2rayAPI.Set("stats", stats)

	experimental.Set("v2ray_api", v2rayAPI)
	cfg.Set("experimental", experimental)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, buildSummary{}, fmt.Errorf("marshal sing-box JSON: %w", err)
	}
	return data, buildSummary{
		Inbounds:  inboundTags,
		Outbounds: outboundTags,
		Users:     users,
	}, nil
}

func orderedMapByKey(container any, key string) (orderedmap.OrderedMap, bool) {
	raw, ok := getField(container, key)
	if !ok {
		return *orderedmap.New(), false
	}
	m, ok := toOrderedMap(raw)
	if !ok {
		return *orderedmap.New(), false
	}
	return m, true
}

func mapAny(v any) map[string]any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func extractTags(v any) []string {
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		tag := getFieldString(item, "tag")
		if tag != "" {
			out = append(out, tag)
		}
	}
	return dedupeKeepOrder(out)
}

func extractUsers(v any) []string {
	inbounds, ok := v.([]any)
	if !ok {
		return nil
	}
	users := make([]string, 0)
	for _, inbound := range inbounds {
		usersRaw, ok := getField(inbound, "users")
		if !ok {
			continue
		}
		usersArr, ok := usersRaw.([]any)
		if !ok {
			continue
		}
		for _, u := range usersArr {
			name := getFieldString(u, "name")
			if name != "" {
				users = append(users, name)
			}
		}
	}
	return dedupeKeepOrder(users)
}

func dedupeKeepOrder(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	out := make([]string, 0, len(in))
	seen := make(map[string]struct{}, len(in))
	for _, item := range in {
		if item == "" {
			continue
		}
		if _, exists := seen[item]; exists {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func nonNilSlice(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

func normalizeLocalRuleSetPathsOrdered(cfg *orderedmap.OrderedMap, baseDir string) {
	if cfg == nil || strings.TrimSpace(baseDir) == "" {
		return
	}
	routeRaw, ok := cfg.Get("route")
	if !ok {
		return
	}
	route, ok := toOrderedMap(routeRaw)
	if !ok {
		return
	}
	rawRuleSetsRaw, ok := route.Get("rule_set")
	if !ok {
		return
	}
	rawRuleSets, ok := rawRuleSetsRaw.([]any)
	if !ok || len(rawRuleSets) == 0 {
		return
	}

	changed := false
	for i, raw := range rawRuleSets {
		ruleSet, ok := toOrderedMap(raw)
		if !ok {
			continue
		}
		rsType := getFieldString(ruleSet, "type")
		if strings.ToLower(strings.TrimSpace(rsType)) != "local" {
			continue
		}
		pathValue := getFieldString(ruleSet, "path")
		pathValue = strings.TrimSpace(pathValue)
		if pathValue == "" || filepath.IsAbs(pathValue) {
			continue
		}
		ruleSet.Set("path", filepath.Clean(filepath.Join(baseDir, pathValue)))
		rawRuleSets[i] = ruleSet
		changed = true
	}

	if changed {
		route.Set("rule_set", rawRuleSets)
		cfg.Set("route", route)
	}
}

func toOrderedMap(v any) (orderedmap.OrderedMap, bool) {
	switch m := v.(type) {
	case orderedmap.OrderedMap:
		return m, true
	case *orderedmap.OrderedMap:
		if m == nil {
			return orderedmap.OrderedMap{}, false
		}
		return *m, true
	default:
		return orderedmap.OrderedMap{}, false
	}
}

func getField(v any, key string) (any, bool) {
	switch m := v.(type) {
	case map[string]any:
		value, ok := m[key]
		return value, ok
	case orderedmap.OrderedMap:
		return m.Get(key)
	case *orderedmap.OrderedMap:
		if m == nil {
			return nil, false
		}
		return m.Get(key)
	default:
		return nil, false
	}
}

func getFieldString(v any, key string) string {
	value, ok := getField(v, key)
	if !ok {
		return ""
	}
	s, _ := value.(string)
	return s
}
