package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cerberus-node/config"
)

// DeployConfigTaskPayload is JSON payload for SubmitTask(operation=deploy_config).
// It accepts either "config" or "singbox_config" with raw sing-box JSON.
type DeployConfigTaskPayload struct {
	Config        json.RawMessage       `json:"config"`
	SingboxConfig json.RawMessage       `json:"singbox_config"`
	Listen        string                `json:"listen"`
	Restart       *bool                 `json:"restart"`
	Stats         DeployConfigTaskStats `json:"stats"`
	SRSLists      []SRSListItem         `json:"srs_lists"`
	Modules       DeployModulesPayload  `json:"modules"`
}

type DeployConfigTaskStats struct {
	Enabled   *bool    `json:"enabled"`
	Inbounds  []string `json:"inbounds"`
	Outbounds []string `json:"outbounds"`
	Users     []string `json:"users"`
}

type DeployModulesPayload struct {
	HaproxyEnabled bool               `json:"haproxy_enabled"`
	HaproxyUsers   []HaproxyUserEntry `json:"haproxy_users"`
}

type HaproxyUserEntry struct {
	Username       string `json:"username"`
	VLESSUUID      string `json:"vless_uuid"`
	TrojanPassword string `json:"trojan_password"`
}

type DeploySummary struct {
	ConfigPath string
	Listen     string
	Inbounds   int
	Outbounds  int
	Users      int
	Restarted  bool
}

// DeployConfig applies sing-box config, injects experimental.v2ray_api stats block and restarts core if requested.
func (s *NodeServer) DeployConfig(task DeployConfigTaskPayload) (DeploySummary, error) {
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

	if len(task.SRSLists) > 0 {
		if summary, err := s.SyncSRSLists(task.SRSLists); err != nil {
			s.Cfg.Logger.Warn("Failed to sync SRS lists during deploy", "error", err)
		} else {
			s.Cfg.Logger.Info(
				"SRS lists synced during deploy",
				"total", summary.Total,
				"configured", summary.Configured,
				"downloaded", summary.Downloaded,
				"failed", summary.Failed,
			)
		}
	}

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
	if err := os.WriteFile(configPath, finalConfig, 0o644); err != nil {
		return DeploySummary{}, fmt.Errorf("write config: %w", err)
	}
	s.Cfg.Logger.Info("Sing-box config written", "path", configPath, "bytes", len(finalConfig))

	if err := applyHaproxyModule(task.Modules); err != nil {
		return DeploySummary{}, err
	}

	shouldRestart := task.Restart != nil && *task.Restart

	restarted := false
	if shouldRestart {
		s.Cfg.Logger.Info("Restart requested by deploy payload")
		if err := restartCoreProcess(s.Cfg); err != nil {
			return DeploySummary{}, err
		}
		if task.Modules.HaproxyEnabled {
			result := restartHaproxyContainer()
			switch {
			case result.Reloaded:
				s.Cfg.Logger.Info("HAProxy container reloaded", "container", "haproxy")
			case result.Restarted:
				s.Cfg.Logger.Warn("HAProxy soft reload failed, restarted container", "container", "haproxy", "warning", result.Warning)
			default:
				s.Cfg.Logger.Warn("HAProxy reload skipped", "container", "haproxy", "warning", result.Warning)
			}
		}
		restarted = true
	} else {
		s.Cfg.Logger.Info("Restart skipped by deploy payload")
	}

	return DeploySummary{
		ConfigPath: configPath,
		Listen:     listen,
		Inbounds:   len(summary.Inbounds),
		Outbounds:  len(summary.Outbounds),
		Users:      len(summary.Users),
		Restarted:  restarted,
	}, nil
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
	var cfg map[string]any
	if err := json.Unmarshal(rawConfig, &cfg); err != nil {
		return nil, buildSummary{}, fmt.Errorf("invalid sing-box JSON: %w", err)
	}
	normalizeLocalRuleSetPaths(cfg, filepath.Dir(config.FixedSingboxConfigPath))

	inboundTags := dedupeKeepOrder(opt.ExplicitInTags)
	if len(inboundTags) == 0 {
		inboundTags = extractTags(cfg["inbounds"])
	}

	outboundTags := dedupeKeepOrder(opt.ExplicitOutTags)
	if len(outboundTags) == 0 {
		outboundTags = extractTags(cfg["outbounds"])
	}

	users := dedupeKeepOrder(opt.ExplicitUsers)
	if len(users) == 0 {
		users = extractUsers(cfg["inbounds"])
	}

	experimental := mapAny(cfg["experimental"])
	v2rayAPI := map[string]any{
		"listen": opt.Listen,
		"stats": map[string]any{
			"enabled":   opt.Enabled,
			"inbounds":  nonNilSlice(inboundTags),
			"outbounds": nonNilSlice(outboundTags),
			"users":     nonNilSlice(users),
		},
	}
	experimental["v2ray_api"] = v2rayAPI
	cfg["experimental"] = experimental

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
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		tag, _ := m["tag"].(string)
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
		inboundMap, ok := inbound.(map[string]any)
		if !ok {
			continue
		}
		usersArr, ok := inboundMap["users"].([]any)
		if !ok {
			continue
		}
		for _, u := range usersArr {
			userMap, ok := u.(map[string]any)
			if !ok {
				continue
			}
			name, _ := userMap["name"].(string)
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

func normalizeLocalRuleSetPaths(cfg map[string]any, baseDir string) {
	route, ok := cfg["route"].(map[string]any)
	if !ok || strings.TrimSpace(baseDir) == "" {
		return
	}
	rawRuleSets, ok := route["rule_set"].([]any)
	if !ok || len(rawRuleSets) == 0 {
		return
	}

	for _, raw := range rawRuleSets {
		ruleSet, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		rsType, _ := ruleSet["type"].(string)
		if strings.ToLower(strings.TrimSpace(rsType)) != "local" {
			continue
		}
		pathValue, _ := ruleSet["path"].(string)
		pathValue = strings.TrimSpace(pathValue)
		if pathValue == "" || filepath.IsAbs(pathValue) {
			continue
		}
		ruleSet["path"] = filepath.Clean(filepath.Join(baseDir, pathValue))
	}
}

func applyHaproxyModule(modules DeployModulesPayload) error {
	const usersFilePath = "/app/haproxy/data/users.csv"

	if !modules.HaproxyEnabled {
		if err := os.Remove(usersFilePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove haproxy users file: %w", err)
		}
		return nil
	}

	lines := make([]string, 0, len(modules.HaproxyUsers)*2)
	for _, user := range modules.HaproxyUsers {
		username := strings.TrimSpace(user.Username)
		if username == "" {
			continue
		}
		if uuid := strings.TrimSpace(user.VLESSUUID); uuid != "" {
			lines = append(lines, fmt.Sprintf("1,%s,%s", username, uuid))
		}
		if trojan := normalizeTrojanHash(user.TrojanPassword); trojan != "" {
			lines = append(lines, fmt.Sprintf("1,%s,%s", username, trojan))
		}
	}

	if err := os.MkdirAll(filepath.Dir(usersFilePath), 0o755); err != nil {
		return fmt.Errorf("create haproxy data dir: %w", err)
	}

	content := ""
	if len(lines) > 0 {
		content = strings.Join(lines, "\n") + "\n"
	}
	if err := os.WriteFile(usersFilePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write haproxy users file: %w", err)
	}

	return nil
}

func normalizeTrojanHash(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	sum := sha256.Sum224([]byte(secret))
	return hex.EncodeToString(sum[:])
}
