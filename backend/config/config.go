package config

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"v2ray-stat/logger"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration settings for the backend.
type BackendConfig struct {
	Logger       *logger.Logger
	Log          LogConfig               `yaml:"log_config"`
	APIToken     string                  `yaml:"api_token"`
	TZ           string                  `yaml:"TZ"`
	Monitor      MonitorConfig           `yaml:"monitor"`
	V2RS         V2RSConfig              `yaml:"backend"`
	Paths        PathsConfig             `yaml:"paths"`
	Subscription *SubscriptionConnConfig `yaml:"subscription"`
	StatsColumns StatsColumns            `yaml:"stats_columns"`
	Nodes        map[string]NodeConfig   `yaml:"node_metadata"`
}

type LogConfig struct {
	LogLevel string `yaml:"level"`
	LogMode  string `yaml:"mode"`
}

// V2RSConfig holds v2ray-stat specific settings.
type V2RSConfig struct {
	Address string `yaml:"listen_address"`
	Port    string `yaml:"listen_port"`
}

// NodeConfig holds configuration for a single node.
type NodeConfig struct {
	NodeName   string
	Schema     string      `yaml:"schema"`
	Address    string      `yaml:"address"`
	Port       string      `yaml:"port"`
	Path       string      `yaml:"path"`
	MTLSConfig *MTLSConfig `yaml:"mtls"`
}

// SubscriptionConfig holds configuration for a subscription.
type SubscriptionConnConfig struct {
	Schema     string      `yaml:"schema"`
	Address    string      `yaml:"address"`
	Port       string      `yaml:"port"`
	Path       string      `yaml:"path"`
	MTLSConfig *MTLSConfig `yaml:"mtls"`
}

// MTLSConfig holds mTLS configuration.
type MTLSConfig struct {
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	CACert string `yaml:"ca_cert"`
}

// MonitorConfig holds monitoring-related settings.
type MonitorConfig struct {
	TickerInterval      int `yaml:"ticker_interval"`
	OnlineRateThreshold int `yaml:"online_rate_threshold"`
}

// PathsConfig holds paths settings.
type PathsConfig struct {
	Database string `yaml:"database"`
}

// StatsColumns holds column configuration for stats display.
type StatsColumns struct {
	Server StatsSection `yaml:"server"`
	Client StatsSection `yaml:"client"`
}

// StatsSection holds columns and sort configuration for a section.
type StatsSection struct {
	Sort      string `yaml:"sort"`
	SortBy    string
	SortOrder string
	Columns   []string `yaml:"columns"`
}

var defaultConfig = BackendConfig{
	Log: LogConfig{
		LogLevel: "none",
		LogMode:  "inclusive",
	},
	APIToken: "",
	TZ:       "UTC",
	V2RS: V2RSConfig{
		Address: "127.0.0.1",
		Port:    "9243",
	},
	Paths: PathsConfig{
		Database: "./data.db",
	},
	Monitor: MonitorConfig{
		TickerInterval:      10,
		OnlineRateThreshold: 0,
	},
	Nodes: make(map[string]NodeConfig),
	Subscription: &SubscriptionConnConfig{
		Schema:     "",
		Address:    "",
		Port:       "",
		Path:       "",
		MTLSConfig: nil,
	},
	StatsColumns: StatsColumns{
		Server: StatsSection{Sort: "source asc", Columns: []string{}},
		Client: StatsSection{Sort: "last_seen desc", Columns: []string{}},
	},
}

// LoadConfig reads configuration from the specified YAML file.
func LoadConfig(configFile string) (BackendConfig, error) {
	cfg := defaultConfig

	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Logger, _ = logger.NewLoggerWithValidation("warn", "inclusive", cfg.TZ, os.Stderr)
			cfg.Logger.Warn("Configuration file not found, using default values", "file", configFile)
			return cfg, nil
		}
		return cfg, fmt.Errorf("error accessing configuration file %s: %v", configFile, err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, fmt.Errorf("error reading configuration file %s: %v", configFile, err)
	}
	cfg.Logger, _ = logger.NewLoggerWithValidation("debug", "inclusive", cfg.TZ, os.Stderr)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing YAML configuration from %s: %v", configFile, err)
	}

	cfg.Logger, err = logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, cfg.TZ, os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Validate configuration
	if cfg.V2RS.Port != "" {
		portNum, err := strconv.Atoi(cfg.V2RS.Port)
		if err != nil || portNum < 1 || portNum > 65535 {
			cfg.Logger.Warn("Invalid v2ray-stat.port, using default", "port", cfg.V2RS.Port, "default", defaultConfig.V2RS.Port)
			cfg.V2RS.Port = defaultConfig.V2RS.Port
		}
	}

	if cfg.Monitor.TickerInterval < 1 {
		cfg.Logger.Warn("Invalid v2ray-stat.monitor.ticker_interval, using default", "value", cfg.Monitor.TickerInterval, "default", defaultConfig.Monitor.TickerInterval)
		cfg.Monitor.TickerInterval = defaultConfig.Monitor.TickerInterval
	}

	if cfg.TZ != "" {
		if _, err := time.LoadLocation(cfg.TZ); err != nil {
			cfg.Logger.Warn("Invalid timezone value, using default", "timezone", cfg.TZ)
			cfg.TZ = defaultConfig.TZ
		}
	}

	// Validate nodes
	valideNodes := make(map[string]NodeConfig)
	for name, node := range cfg.Nodes {
		if name == "" {
			cfg.Logger.Warn("Invalid node configuration, skipping", "node_name", name, "address", node.Address, "port", node.Port)
			continue
		}
		if portNum, err := strconv.Atoi(node.Port); err != nil || portNum < 1 || portNum > 65535 {
			cfg.Logger.Warn("Invalid node port, skipping", "node_name", name, "port", node.Port)
			continue
		}

		validSchemas := []string{"", "http", "https", "grpc"}
		if !slices.Contains(validSchemas, node.Schema) {
			cfg.Logger.Warn("Invalid schema, fallback to insecure", "node_name", name, "port", node.Port)
			node.Schema = ""
		}
		if node.MTLSConfig != nil {
			if node.MTLSConfig.CACert == "" || node.MTLSConfig.Cert == "" || node.MTLSConfig.Key == "" {
				cfg.Logger.Warn("Incomplete mTLS configuration for node, disabling mTLS", "node_name", name)
				node.MTLSConfig = nil
			} else {
				mtlsValid := true
				for _, file := range []string{node.MTLSConfig.CACert, node.MTLSConfig.Cert, node.MTLSConfig.Key} {
					if _, err := os.Stat(file); os.IsNotExist(err) {
						cfg.Logger.Warn("mTLS certificate file not found for node, disabling mTLS", "node_name", name, "file", file)
						mtlsValid = false
						break
					}
				}
				if !mtlsValid {
					node.MTLSConfig = nil
				}
			}
		}
		node.NodeName = name
		valideNodes[name] = node
	}
	cfg.Nodes = valideNodes

	// Validate subscription
	if cfg.Subscription != nil {
		if cfg.Subscription.Address == "" || cfg.Subscription.Port == "" {
			cfg.Logger.Warn("Invalid subscription configuration, disabling", "address", cfg.Subscription.Address, "port", cfg.Subscription.Port)
			cfg.Subscription = nil
		} else {
			if portNum, err := strconv.Atoi(cfg.Subscription.Port); err != nil || portNum < 1 || portNum > 65535 {
				cfg.Logger.Warn("Invalid subscription port, disabling", "port", cfg.Subscription.Port)
				cfg.Subscription = nil
			} else {
				if cfg.Subscription.MTLSConfig != nil {
					if cfg.Subscription.MTLSConfig.Cert == "" || cfg.Subscription.MTLSConfig.Key == "" || cfg.Subscription.MTLSConfig.CACert == "" {
						cfg.Logger.Warn("Incomplete mTLS configuration for subscription, disabling mTLS")
						cfg.Subscription.MTLSConfig = nil
					} else {
						mtlsValid := true
						for _, file := range []string{cfg.Subscription.MTLSConfig.Cert, cfg.Subscription.MTLSConfig.Key, cfg.Subscription.MTLSConfig.CACert} {
							if _, err := os.Stat(file); os.IsNotExist(err) {
								cfg.Logger.Warn("mTLS certificate file not found for subscription, disabling mTLS", "file", file)
								mtlsValid = false
								break
							}
						}
						if !mtlsValid {
							cfg.Subscription.MTLSConfig = nil
						}
					}
				}
			}
		}
	}

	if len(cfg.StatsColumns.Server.Columns) == 0 {
		cfg.StatsColumns.Server.Columns = []string{} // Or set defaults if desired
	}
	if len(cfg.StatsColumns.Client.Columns) == 0 {
		cfg.StatsColumns.Client.Columns = []string{}
	}

	// Validate columns
	validServerColumns := []string{"node_name", "source", "rate", "uplink", "downlink", "sess_uplink", "sess_downlink"}
	validClientColumns := []string{"node_name", "user", "last_seen", "rate", "uplink", "downlink", "sess_uplink", "sess_downlink", "enabled", "sub_end", "renew", "lim_ip", "ips", "created", "id", "inbound_tag", "traffic_cap"}
	cfg.StatsColumns.Server.Columns = filterValidColumns(cfg.StatsColumns.Server.Columns, validServerColumns, cfg.Logger, "server")
	cfg.StatsColumns.Client.Columns = filterValidColumns(cfg.StatsColumns.Client.Columns, validClientColumns, cfg.Logger, "client")

	// Validate sort configuration
	cfg.StatsColumns.Server.SortBy, cfg.StatsColumns.Server.SortOrder = validateSort("Server", cfg.StatsColumns.Server.Sort, validServerColumns, cfg.Logger)
	cfg.StatsColumns.Client.SortBy, cfg.StatsColumns.Client.SortOrder = validateSort("Client", cfg.StatsColumns.Client.Sort, validClientColumns, cfg.Logger)

	cfg.Logger.Debug("Configuration validated", "nodes_count", len(cfg.Nodes), "subscription_defined", cfg.Subscription != nil)
	return cfg, nil
}

func filterValidColumns(columns, valid []string, log *logger.Logger, section string) []string {
	var filtered []string
	for _, col := range columns {
		if contains(valid, col) {
			filtered = append(filtered, col)
		} else {
			log.Warn("Invalid custom column, ignoring", "section", section, "column", col)
		}
	}
	return filtered
}

func validateSort(section, sortStr string, validColumns []string, log *logger.Logger) (string, string) {
	defaultBy, defaultOrder := "source", "asc"
	if section == "Client" {
		defaultBy = "user"
	}
	if sortStr == "" {
		return defaultBy, defaultOrder
	}
	parts := strings.Fields(sortStr)
	if len(parts) != 2 {
		log.Warn("Invalid sort format, using default", "section", section, "sort", sortStr)
		return defaultBy, defaultOrder
	}
	column, order := parts[0], strings.ToLower(parts[1])
	if !contains(validColumns, column) {
		log.Warn("Invalid sort column, using default", "section", section, "column", column)
		return defaultBy, defaultOrder
	}
	if order != "asc" && order != "desc" {
		log.Warn("Invalid sort order, using asc", "section", section, "order", order)
		order = "asc"
	}
	return column, order
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
