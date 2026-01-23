package config

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"v2ray-stat/logger"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// Config holds the configuration settings for the backend.
type Config struct {
	Log          LogConfig           `yaml:"log_config"`
	V2RSSub      V2RSSubConfig       `yaml:"subscription"`
	TZ           string              `yaml:"TZ"`
	NodeMetadata map[string]NodeMeta `yaml:"node_metadata"`
	Subscription SubscriptionConfig  `yaml:"sub_config"`
	Logger       *logger.Logger
	Ctx          context.Context
}

type LogConfig struct {
	LogLevel string `yaml:"level"`
	LogMode  string `yaml:"mode"`
}

// V2RSSubConfig holds v2rs-sub specific settings.
type V2RSSubConfig struct {
	Address  string `yaml:"listen_address"`
	Port     int    `yaml:"listen_port"`
	GrpcPort int    `yaml:"listen_grpc_port"`
}

type SubscriptionConfig struct {
	Defaults DefaultsConfig        `yaml:"defaults"`
	Groups   map[string]UserConfig `yaml:"groups"`
	Users    map[string]*UserGroup `yaml:"users"`
	UserMap  map[string]UserConfig
}

type NodeMeta struct {
	RemarkPlaceholder string `yaml:"remark_placeholder"`
	DomainPlaceholder string `yaml:"domain_placeholder"`
	IPPlaceholder     string `yaml:"ip_placeholder"`
	PortPlaceholder   string `yaml:"port_placeholder"`
}

type DefaultsConfig struct {
	Clients       []string                     `yaml:"clients"`
	IncludeNodes  []string                     `yaml:"nodes"`
	NodeTemplates map[string]map[string]string `yaml:"templates"`
	Headers       map[string]string            `yaml:"headers"`
}

type UserConfig struct {
	Group         string                       `yaml:"group"`
	Clients       []string                     `yaml:"clients"`
	IncludeNodes  []string                     `yaml:"nodes"`
	NodeTemplates map[string]map[string]string `yaml:"templates"`
	Headers       map[string]string            `yaml:"headers"`
}

type UserGroup struct {
	UserConfig
	Users []string `yaml:"users"`
}

var (
	config Config
	mu     sync.RWMutex
)

var defaultConfig = Config{
	Log: LogConfig{
		LogLevel: "none",
		LogMode:  "inclusive",
	},
	V2RSSub: V2RSSubConfig{
		Address:  "127.0.0.1",
		Port:     9964,
		GrpcPort: 9963,
	},
	TZ:           "",
	NodeMetadata: map[string]NodeMeta{},
	Subscription: SubscriptionConfig{
		Defaults: DefaultsConfig{
			Clients:      []string{},
			IncludeNodes: []string{},
			Headers:      map[string]string{},
		},
		Groups:  map[string]UserConfig{},
		Users:   map[string]*UserGroup{},
		UserMap: map[string]UserConfig{},
	},
}

// LoadConfig reads configuration from the specified YAML file.
func LoadConfig(configFile string) (Config, error) {
	cfg := defaultConfig
	cfg.Logger, _ = logger.NewLoggerWithValidation("debug", "inclusive", cfg.TZ, os.Stderr)
	cfg.Logger.Trace("Attempting to load config file", "file", configFile)

	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Logger, _ = logger.NewLoggerWithValidation("warn", "inclusive", cfg.TZ, os.Stderr)
			cfg.Logger.Warn("Configuration file not found, using default values", "file", configFile)
			mu.Lock()
			config = cfg
			mu.Unlock()
			return cfg, nil
		}
		cfg.Logger.Error("Error accessing configuration file", "file", configFile, "error", err)
		return cfg, fmt.Errorf("error accessing configuration file %s: %v", configFile, err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		cfg.Logger.Error("Error reading configuration file", "file", configFile, "error", err)
		return cfg, fmt.Errorf("error reading configuration file %s: %v", configFile, err)
	}
	cfg.Logger.Debug("Read config file successfully", "file", configFile)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		cfg.Logger.Error("Error parsing YAML configuration", "file", configFile, "error", err)
		return cfg, fmt.Errorf("error parsing YAML configuration from %s: %v", configFile, err)
	}

	cfg.Logger, err = logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, cfg.TZ, os.Stderr)
	if err != nil {
		cfg.Logger.Error("Failed to initialize logger", "error", err)
		return cfg, fmt.Errorf("failed to initialize logger: %v", err)
	}
	cfg.Logger.Info("Logger initialized with level and mode", "level", cfg.Log.LogLevel, "mode", cfg.Log.LogMode)

	// Normalize users into UserMap
	cfg.Subscription.UserMap = make(map[string]UserConfig)
	for groupName, userGroup := range cfg.Subscription.Users {
		if len(userGroup.Users) > 0 {
			for _, u := range userGroup.Users {
				userConfig := UserConfig{
					Group:         groupName,
					Clients:       slices.Clone(userGroup.Clients),
					IncludeNodes:  slices.Clone(userGroup.IncludeNodes),
					NodeTemplates: make(map[string]map[string]string),
					Headers:       make(map[string]string),
				}
				for mode, templates := range userGroup.NodeTemplates {
					userConfig.NodeTemplates[mode] = make(map[string]string)
					maps.Copy(userConfig.NodeTemplates[mode], templates)
				}
				maps.Copy(userConfig.Headers, userGroup.Headers)
				if override, exists := cfg.Subscription.Users[u]; exists && len(override.Users) == 0 {
					if override.Clients != nil {
						userConfig.Clients = slices.Clone(override.Clients)
					}
					if override.IncludeNodes != nil {
						userConfig.IncludeNodes = slices.Clone(override.IncludeNodes)
					}
					if override.NodeTemplates != nil {
						userConfig.NodeTemplates = make(map[string]map[string]string)
						for mode, templates := range override.NodeTemplates {
							userConfig.NodeTemplates[mode] = make(map[string]string)
							maps.Copy(userConfig.NodeTemplates[mode], templates)
						}
					}
					if override.Headers != nil {
						userConfig.Headers = make(map[string]string)
						maps.Copy(userConfig.Headers, override.Headers)
					}
				}
				cfg.Subscription.UserMap[u] = userConfig
			}
		} else {
			userConfig := UserConfig{
				Group:         userGroup.Group,
				Clients:       slices.Clone(userGroup.Clients),
				IncludeNodes:  slices.Clone(userGroup.IncludeNodes),
				NodeTemplates: make(map[string]map[string]string),
				Headers:       make(map[string]string),
			}
			for mode, templates := range userGroup.NodeTemplates {
				userConfig.NodeTemplates[mode] = make(map[string]string)
				maps.Copy(userConfig.NodeTemplates[mode], templates)
			}
			maps.Copy(userConfig.Headers, userGroup.Headers)
			cfg.Subscription.UserMap[groupName] = userConfig
		}
	}

	if cfg.V2RSSub.Port != "" {
		if cfg.V2RSSub.Port < 1 || cfg.V2RSSub.Port > 65535 {
			cfg.Logger.Warn("Invalid v2rs-sub.port, using default", "port", cfg.V2RSSub.Port, "default", defaultConfig.V2RSSub.Port)
			cfg.V2RSSub.Port = defaultConfig.V2RSSub.Port
		}
	}

	if cfg.V2RSSub.GrpcPort != "" {
		if cfg.V2RSSub.GrpcPort < 1 || cfg.V2RSSub.GrpcPort > 65535 {
			cfg.Logger.Warn("Invalid grpc_port, using default", "port", cfg.V2RSSub.GrpcPort, "default", defaultConfig.V2RSSub.GrpcPort)
			cfg.V2RSSub.GrpcPort = defaultConfig.V2RSSub.GrpcPort
		}
	}

	if cfg.TZ != "" {
		if _, err := time.LoadLocation(cfg.TZ); err != nil {
			cfg.Logger.Warn("Invalid timezone value, using default", "TZ", cfg.TZ)
			cfg.TZ = defaultConfig.TZ
		}
	}

	for node, meta := range cfg.NodeMetadata {
		if meta.DomainPlaceholder == "" {
			cfg.Logger.Warn("Empty domain_placeholder for node", "node", node)
		}
	}

	if len(cfg.Subscription.Defaults.Clients) == 0 {
		cfg.Logger.Warn("defaults.clients is empty, using empty list")
	}
	if len(cfg.Subscription.Defaults.IncludeNodes) == 0 {
		cfg.Logger.Warn("defaults.nodes is empty, using empty list")
	}
	for mode, templates := range cfg.Subscription.Defaults.NodeTemplates {
		if len(templates) == 0 {
			cfg.Logger.Warn("defaults.templates is empty for mode", "mode", mode)
		}
		for _, node := range cfg.Subscription.Defaults.IncludeNodes {
			if _, ok := templates[node]; !ok {
				cfg.Logger.Warn("defaults: no template specified for node in mode", "node", node, "mode", mode)
			}
			if _, ok := cfg.NodeMetadata[node]; !ok {
				cfg.Logger.Warn("defaults: no metadata specified for node", "node", node)
			}
		}
	}

	for groupName, group := range cfg.Subscription.Groups {
		if len(group.Clients) == 0 {
			cfg.Logger.Warn("group clients is empty", "group", groupName)
		}
		if len(group.IncludeNodes) == 0 {
			cfg.Logger.Warn("group nodes is empty", "group", groupName)
		}
		for mode, templates := range group.NodeTemplates {
			if len(templates) == 0 {
				cfg.Logger.Warn("group templates is empty for mode", "group", groupName, "mode", mode)
			}
			for _, node := range group.IncludeNodes {
				if _, ok := templates[node]; !ok {
					cfg.Logger.Warn("group: no template specified for node in mode", "group", groupName, "node", node, "mode", mode)
				}
				if _, ok := cfg.NodeMetadata[node]; !ok {
					cfg.Logger.Warn("group: no metadata specified for node", "group", groupName, "node", node)
				}
			}
		}
	}

	for userName, user := range cfg.Subscription.UserMap {
		if user.Group != "" {
			if _, ok := cfg.Subscription.Groups[user.Group]; !ok {
				cfg.Logger.Warn("user: group not found", "user", userName, "group", user.Group)
			}
		}
		if len(user.IncludeNodes) > 0 {
			for mode, templates := range user.NodeTemplates {
				for _, node := range user.IncludeNodes {
					if _, ok := templates[node]; !ok {
						if user.Group != "" {
							if group, ok := cfg.Subscription.Groups[user.Group]; ok {
								if _, ok := group.NodeTemplates[mode][node]; !ok {
									cfg.Logger.Warn("user: no template specified for node in group or user config for mode", "user", userName, "node", node, "group", user.Group, "mode", mode)
								}
							}
						} else if _, ok := cfg.Subscription.Defaults.NodeTemplates[mode][node]; !ok {
							cfg.Logger.Warn("user: no template specified for node in user config or defaults for mode", "user", userName, "node", node, "mode", mode)
						}
					}
					if _, ok := cfg.NodeMetadata[node]; !ok {
						cfg.Logger.Warn("user: no metadata specified for node", "user", userName, "node", node)
					}
				}
			}
		}
	}

	cfg.Logger.Debug("Configuration validated")
	cfg.Logger.Info("Configuration loaded successfully", "file", configFile)

	mu.Lock()
	config = cfg
	mu.Unlock()

	return cfg, nil
}

func GetConfig() Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

// WatchConfig watches the config file for changes and reloads it.
func WatchConfig(ctx context.Context, cfg *Config, wg *sync.WaitGroup) {
	defer wg.Done()
	if cfg.Logger == nil {
		return
	}
	cfg.Logger.Trace("Starting to watch config file", "file", "config.yml")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		cfg.Logger.Error("Failed to create config watcher", "error", err)
		return
	}
	defer watcher.Close()

	err = watcher.Add("config.yml")
	if err != nil {
		cfg.Logger.Error("Failed to add watch for config.yml", "error", err)
		return
	}
	cfg.Logger.Debug("Added watch for config.yml")

	for {
		select {
		case <-ctx.Done():
			cfg.Logger.Debug("Stopping config watcher due to context cancellation")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				cfg.Logger.Warn("Config watcher closed unexpectedly")
				return
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				cfg.Logger.Info("config.yml changed, reloading...")
				if newCfg, err := LoadConfig("config.yml"); err != nil {
					cfg.Logger.Error("Failed to reload config", "error", err)
				} else {
					mu.Lock()
					config = newCfg
					mu.Unlock()
					cfg.Logger.Info("Config reloaded successfully")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				cfg.Logger.Warn("Config watcher error channel closed")
				return
			}
			cfg.Logger.Error("Watcher error", "error", err)
		}
	}
}
