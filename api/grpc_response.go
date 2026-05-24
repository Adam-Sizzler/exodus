package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

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

	networkMu            sync.Mutex
	defaultInterface     string
	previousNetworkStats map[string]networkCounter
}

type systemInfo struct {
	Arch              string   `json:"arch"`
	CPUs              int      `json:"cpus"`
	CPUModel          string   `json:"cpuModel"`
	MemoryTotal       uint64   `json:"memoryTotal"`
	Hostname          string   `json:"hostname"`
	Platform          string   `json:"platform"`
	Release           string   `json:"release"`
	Type              string   `json:"type"`
	Version           string   `json:"version"`
	NetworkInterfaces []string `json:"networkInterfaces"`
}

type systemInterfaceStats struct {
	Interface     string  `json:"interface"`
	RXBytesPerSec float64 `json:"rxBytesPerSec"`
	TXBytesPerSec float64 `json:"txBytesPerSec"`
	RXTotal       uint64  `json:"rxTotal"`
	TXTotal       uint64  `json:"txTotal"`
}

type systemStats struct {
	MemoryFree uint64                `json:"memoryFree"`
	MemoryUsed uint64                `json:"memoryUsed"`
	Uptime     float64               `json:"uptime"`
	LoadAvg    []float64             `json:"loadAvg"`
	Interface  *systemInterfaceStats `json:"interface"`
}

type networkCounter struct {
	rxBytes   uint64
	txBytes   uint64
	timestamp time.Time
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
		cfg:                  cfg,
		api:                  coreAPI,
		singboxVersion:       detectSingboxVersion(),
		nodeVersion:          constant.Version,
		cpuCount:             runtime.NumCPU(),
		cpuModel:             detectCPUModel(),
		totalRAMBytes:        detectTotalRAMBytes(),
		defaultInterface:     detectDefaultInterface(),
		previousNetworkStats: readNetworkCounters(),
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
	if infoJSON, statsJSON, err := s.collectSystemStatsJSON(); err != nil {
		s.cfg.Logger.Warn("Failed to collect system stats", "error", err)
	} else {
		result.Stat = append(result.Stat,
			Stat{Name: "system_info", Value: infoJSON},
			Stat{Name: "system_stats", Value: statsJSON},
		)
	}
	s.cfg.Logger.Debug("Retrieved core stats", "count", len(result.Stat), "core_type", config.FixedCoreType)

	return result, nil
}

func (s *Service) collectSystemStatsJSON() (string, string, error) {
	info := systemInfo{
		Arch:              runtime.GOARCH,
		CPUs:              s.cpuCount,
		CPUModel:          s.cpuModel,
		MemoryTotal:       s.totalRAMBytes,
		Hostname:          detectHostname(),
		Platform:          runtime.GOOS,
		Release:           detectKernelRelease(),
		Type:              detectOSType(),
		Version:           detectOSVersion(),
		NetworkInterfaces: detectNetworkInterfaces(),
	}

	memoryFree := detectAvailableRAMBytes()
	stats := systemStats{
		MemoryFree: memoryFree,
		MemoryUsed: saturatingSub(s.totalRAMBytes, memoryFree),
		Uptime:     detectSystemUptime(),
		LoadAvg:    detectLoadAverage(),
		Interface:  s.collectInterfaceStats(),
	}

	infoBytes, err := json.Marshal(info)
	if err != nil {
		return "", "", err
	}
	statsBytes, err := json.Marshal(stats)
	if err != nil {
		return "", "", err
	}
	return string(infoBytes), string(statsBytes), nil
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

func detectAvailableRAMBytes() uint64 {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, "MemAvailable:") {
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

func detectHostname() string {
	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		return "unknown"
	}
	return hostname
}

func detectKernelRelease() string {
	out, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return "unknown"
	}
	if release := strings.TrimSpace(string(out)); release != "" {
		return release
	}
	return "unknown"
}

func detectOSType() string {
	if runtime.GOOS == "linux" {
		return "Linux"
	}
	return runtime.GOOS
}

func detectOSVersion() string {
	if content, err := os.ReadFile("/proc/version"); err == nil {
		if version := strings.TrimSpace(string(content)); version != "" {
			return version
		}
	}
	return detectKernelRelease()
}

func detectSystemUptime() float64 {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	fields := strings.Fields(string(content))
	if len(fields) == 0 {
		return 0
	}
	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return uptime
}

func detectLoadAverage() []float64 {
	content, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return []float64{0, 0, 0}
	}
	fields := strings.Fields(string(content))
	load := []float64{0, 0, 0}
	for i := 0; i < len(load) && i < len(fields); i++ {
		value, err := strconv.ParseFloat(fields[i], 64)
		if err == nil {
			load[i] = value
		}
	}
	return load
}

func detectNetworkInterfaces() []string {
	seen := make(map[string]struct{})

	ifaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range ifaces {
			if strings.TrimSpace(iface.Name) != "" {
				seen[iface.Name] = struct{}{}
			}
		}
	}

	counters := readNetworkCounters()
	for iface := range counters {
		seen[iface] = struct{}{}
	}

	interfaces := make([]string, 0, len(seen))
	for iface := range seen {
		interfaces = append(interfaces, iface)
	}
	sort.Strings(interfaces)
	return interfaces
}

func detectDefaultInterface() string {
	content, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(content), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[1] == "00000000" {
			return strings.TrimSpace(fields[0])
		}
	}
	return ""
}

func readNetworkCounters() map[string]networkCounter {
	content, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return map[string]networkCounter{}
	}
	now := time.Now()
	counters := make(map[string]networkCounter)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 10 {
			continue
		}
		iface := strings.TrimSuffix(parts[0], ":")
		rxBytes, rxErr := strconv.ParseUint(parts[1], 10, 64)
		txBytes, txErr := strconv.ParseUint(parts[9], 10, 64)
		if iface == "" || rxErr != nil || txErr != nil {
			continue
		}
		counters[iface] = networkCounter{
			rxBytes:   rxBytes,
			txBytes:   txBytes,
			timestamp: now,
		}
	}
	return counters
}

func (s *Service) collectInterfaceStats() *systemInterfaceStats {
	current := readNetworkCounters()
	if len(current) == 0 {
		return nil
	}

	s.networkMu.Lock()
	defer s.networkMu.Unlock()

	if s.defaultInterface == "" {
		s.defaultInterface = detectDefaultInterface()
	}
	if s.defaultInterface == "" {
		interfaces := make([]string, 0, len(current))
		for iface := range current {
			interfaces = append(interfaces, iface)
		}
		sort.Strings(interfaces)
		for _, iface := range interfaces {
			if iface != "lo" {
				s.defaultInterface = iface
				break
			}
		}
	}
	if s.defaultInterface == "" {
		return nil
	}

	counter, ok := current[s.defaultInterface]
	if !ok {
		return nil
	}

	result := &systemInterfaceStats{
		Interface: s.defaultInterface,
		RXTotal:   counter.rxBytes,
		TXTotal:   counter.txBytes,
	}
	if previous, ok := s.previousNetworkStats[s.defaultInterface]; ok {
		elapsed := counter.timestamp.Sub(previous.timestamp).Seconds()
		if elapsed > 0 {
			result.RXBytesPerSec = float64(saturatingSub(counter.rxBytes, previous.rxBytes)) / elapsed
			result.TXBytesPerSec = float64(saturatingSub(counter.txBytes, previous.txBytes)) / elapsed
		}
	}
	s.previousNetworkStats = current

	return result
}

func saturatingSub(current uint64, previous uint64) uint64 {
	if current <= previous {
		return 0
	}
	return current - previous
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
