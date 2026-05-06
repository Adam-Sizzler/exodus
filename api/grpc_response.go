package api

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"exodus-node/config"
	"exodus-node/constant"
	"exodus-node/sdk"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stat represents a single statistic entry.
type Stat struct {
	Name  string
	Value string
}

// ApiResponse contains the collected statistics.
type ApiResponse struct {
	Stat []Stat
}

// Service is a thin API facade over core SDK for node gRPC handlers.
type Service struct {
	cfg *config.NodeConfig
	api *sdk.API

	singboxVersion string
	nodeVersion    string
	cpuCount       int
	cpuModel       string
	totalRAMBytes  uint64
}

func NewService(cfg *config.NodeConfig) (*Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil node config")
	}
	coreAPI, err := sdk.New(sdk.Config{
		CoreType: config.FixedCoreType,
		Address:  config.FixedCoreAPIAddress,
		Port:     config.FixedCoreAPIGRPCPort,
	})
	if err != nil {
		return nil, fmt.Errorf("initialize core SDK: %w", err)
	}
	return &Service{
		cfg:            cfg,
		api:            coreAPI,
		singboxVersion: detectSingboxVersion(),
		nodeVersion:    constant.Version,
		cpuCount:       runtime.NumCPU(),
		cpuModel:       detectCPUModel(),
		totalRAMBytes:  detectTotalRAMBytes(),
	}, nil
}

func (s *Service) Close() error {
	if s == nil || s.api == nil {
		return nil
	}
	return s.api.Close()
}

// GetApiResponse retrieves statistics from the configured core Stats API.
func (s *Service) GetApiResponse(ctx context.Context) (*ApiResponse, error) {
	if s == nil || s.api == nil || s.api.Stats == nil {
		return nil, status.Error(codes.FailedPrecondition, "core SDK is not initialized")
	}

	stats, err := s.api.Stats.QueryStats(ctx, sdk.QueryOptions{
		Patterns: []string{
			`^inbound>>>.*>>>traffic>>>(?:uplink|downlink)$`,
			`^outbound>>>.*>>>traffic>>>(?:uplink|downlink)$`,
			`^user>>>.*>>>traffic>>>(?:uplink|downlink)$`,
		},
		Regexp: true,
		Reset:  true,
	})
	if err != nil {
		s.cfg.Logger.Error("Failed to execute core stats query", "error", err, "core_type", config.FixedCoreType)
		return nil, fmt.Errorf("query core stats: %w", err)
	}

	result := &ApiResponse{
		Stat: make([]Stat, 0, len(stats)),
	}
	for _, item := range stats {
		s.cfg.Logger.Trace("Processing core stat", "name", item.Name, "value", item.Value, "core_type", config.FixedCoreType)
		result.Stat = append(result.Stat, Stat{
			Name:  item.Name,
			Value: strconv.FormatInt(item.Value, 10),
		})
	}

	singboxUptimeSeconds := int64(0)
	sysStats, sysErr := s.api.Stats.GetSysStats(ctx)
	if sysErr != nil {
		s.cfg.Logger.Warn("Failed to read core sys stats", "error", sysErr, "core_type", config.FixedCoreType)
	} else if sysStats != nil {
		singboxUptimeSeconds = int64(sysStats.Uptime)
	}

	// Runtime metadata consumed by backend node monitor.
	result.Stat = append(result.Stat,
		Stat{Name: "singbox_version", Value: s.singboxVersion},
		Stat{Name: "node_version", Value: s.nodeVersion},
		Stat{Name: "singbox_uptime", Value: strconv.FormatInt(singboxUptimeSeconds, 10)},
		Stat{Name: "cpu_count", Value: strconv.Itoa(s.cpuCount)},
		Stat{Name: "cpu_model", Value: s.cpuModel},
		Stat{Name: "total_ram", Value: formatIECBytes(s.totalRAMBytes)},
	)
	s.cfg.Logger.Debug("Retrieved core stats", "count", len(result.Stat), "core_type", config.FixedCoreType)

	return result, nil
}

func detectSingboxVersion() string {
	out, err := exec.Command("/usr/local/bin/sing-box", "version").Output()
	if err != nil {
		return "unknown"
	}
	line := strings.TrimSpace(strings.SplitN(string(out), "\n", 2)[0])
	// Example: "sing-box version 1.13.x"
	if idx := strings.LastIndex(line, " "); idx > 0 && idx+1 < len(line) {
		return strings.TrimSpace(line[idx+1:])
	}
	return line
}

func detectCPUModel() string {
	content, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "unknown"
}

func detectTotalRAMBytes() uint64 {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				break
			}
			kib, parseErr := strconv.ParseUint(fields[1], 10, 64)
			if parseErr != nil {
				break
			}
			return kib * 1024
		}
	}
	return 0
}

func formatIECBytes(bytes uint64) string {
	if bytes == 0 {
		return "0 B"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	size := float64(bytes)
	idx := 0
	for size >= 1024 && idx < len(units)-1 {
		size /= 1024
		idx++
	}
	if idx == 0 {
		return fmt.Sprintf("%.0f %s", size, units[idx])
	}
	value := math.Round(size*100) / 100
	return fmt.Sprintf("%.2f %s", value, units[idx])
}
