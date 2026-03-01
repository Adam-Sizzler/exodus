package logprocessor

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"v2ray-stat/logger"
	"v2ray-stat/backend/node/config"
	"v2ray-stat/proto"
)

type LogSource interface {
	FetchNewLines() ([]string, error)
	Close() error
}

type FileLogSource struct {
	path   string
	file   *os.File
	offset int64
	logger *logger.Logger
	mu     sync.Mutex
}

func NewFileLogSource(path string, logger *logger.Logger) (*FileLogSource, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		logger.Error("Failed to open log file", "path", path, "error", err)
		return nil, err
	}

	offset, err := f.Seek(0, 2)
	if err != nil {
		f.Close()
		logger.Error("Failed to seek to end of file", "path", path, "error", err)
		return nil, err
	}

	fs := &FileLogSource{
		path:   path,
		file:   f,
		offset: offset,
		logger: logger,
	}

	logger.Info("File log source initialized", "path", path, "initial_offset", offset)
	go fs.runDailyCleanup()

	return fs, nil
}

func (fs *FileLogSource) FetchNewLines() ([]string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	_, err := fs.file.Seek(fs.offset, 0)
	if err != nil {
		fs.logger.Error("Failed to seek in log file", "path", fs.path, "offset", fs.offset, "error", err)
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(fs.file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fs.logger.Error("Error scanning log file", "path", fs.path, "error", err)
		return nil, err
	}

	pos, _ := fs.file.Seek(0, 1)
	fs.offset = pos

	if len(lines) > 0 {
		fs.logger.Debug("Read new lines from file", "count", len(lines), "new_offset", pos)
	} else {
		fs.logger.Trace("No new lines found in file")
	}

	return lines, nil
}

func (fs *FileLogSource) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := fs.file.Close(); err != nil {
		fs.logger.Error("Error closing log file", "path", fs.path, "error", err)
		return err
	}

	fs.logger.Info("Log file closed", "path", fs.path)
	return nil
}

func (fs *FileLogSource) runDailyCleanup() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		fs.mu.Lock()

		if err := fs.file.Close(); err != nil {
			fs.logger.Error("Error closing file before truncation", "path", fs.path, "error", err)
		}

		f, err := os.OpenFile(fs.path, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			fs.logger.Error("Failed to recreate file after truncation", "path", fs.path, "error", err)
			fs.mu.Unlock()
			continue
		}

		fs.file = f
		fs.offset = 0
		fs.mu.Unlock()

		fs.logger.Info("Log file successfully truncated (daily rotation)", "path", fs.path)
	}
}

// ------------------- Log Processor -------------------

type LogProcessor struct {
	cfg    *config.NodeConfig
	source LogSource
	logger *logger.Logger
	regex  *regexp.Regexp
}

func NewLogProcessor(cfg *config.NodeConfig) (*LogProcessor, error) {
	var source LogSource
	var err error

	if cfg.Core.LogSource == "journal" && cfg.Core.LogServiceName != "" {
		cfg.Logger.Info("Log source selected: journald", "service", cfg.Core.LogServiceName)
		// Этот метод будет определен либо в journald.go, либо в journald_stub.go
		source, err = NewJournalLogSource(cfg.Core.LogServiceName, cfg.Logger)
	} else {
		cfg.Logger.Info("Log source selected: file", "path", cfg.Core.AccessLog)
		source, err = NewFileLogSource(cfg.Core.AccessLog, cfg.Logger)
	}

	if err != nil {
		cfg.Logger.Error("Failed to initialize log source", "type", cfg.Core.LogSource, "error", err)
		return nil, err
	}

	re, err := regexp.Compile(cfg.Core.AccessLogRegex)
	if err != nil {
		source.Close()
		cfg.Logger.Error("Invalid access log regex", "regex", cfg.Core.AccessLogRegex, "error", err)
		return nil, fmt.Errorf("invalid regex: %w", err)
	}

	cfg.Logger.Debug("Access log regex compiled successfully")
	return &LogProcessor{
		cfg:    cfg,
		source: source,
		logger: cfg.Logger,
		regex:  re,
	}, nil
}

func ProcessLogLine(line string, dnsStats map[string]map[string]int, cfg *config.NodeConfig, re *regexp.Regexp) (string, []string, bool) {
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return "", nil, false
	}
	uIdx, iIdx, dIdx := cfg.Core.LogMap.User, cfg.Core.LogMap.IP, cfg.Core.LogMap.Domain
	var user, ip, domain string
	if uIdx < len(matches) {
		user = strings.TrimSpace(matches[uIdx])
	}
	if user == "" {
		return "", nil, false
	}
	if iIdx < len(matches) {
		ip = strings.TrimSpace(matches[iIdx])
	}
	if dIdx < len(matches) {
		domain = strings.TrimSpace(matches[dIdx])
	}
	validIPs := []string{}
	if ip != "" {
		validIPs = append(validIPs, ip)
	}
	if dnsStats[user] == nil {
		dnsStats[user] = make(map[string]int)
	}
	if domain != "" {
		dnsStats[user][domain]++
	}
	return user, validIPs, true
}

func (lp *LogProcessor) ReadNewLines() (*proto.GetLogDataResponse, error) {
	lines, err := lp.source.FetchNewLines()
	if err != nil {
		lp.logger.Error("Failed to fetch new log lines", "error", err)
		return nil, err
	}
	lp.logger.Debug("Received lines for processing", "count", len(lines))
	dnsStats := make(map[string]map[string]int)
	ipUpdates := make(map[string]map[string]struct{})
	for _, line := range lines {
		lp.logger.Trace("Processing log line", "line", line)
		user, validIPs, ok := ProcessLogLine(line, dnsStats, lp.cfg, lp.regex)
		if ok {
			if ipUpdates[user] == nil {
				ipUpdates[user] = make(map[string]struct{})
			}
			for _, ip := range validIPs {
				ipUpdates[user][ip] = struct{}{}
			}
			lp.logger.Trace("Extracted user data",
				"user", user,
				"ip_count", len(validIPs),
				"domains_count", len(dnsStats[user]),
			)
		} else {
			lp.logger.Debug("Line did not match regex", "line", line)
		}
	}
	if len(dnsStats) > 0 {
		lp.logger.Debug("DNS stats collected", "users_with_domains", len(dnsStats))
		for user, domains := range dnsStats {
			lp.logger.Trace("DNS domains for user",
				"user", user,
				"domains", mapKeysToString(domains),
				"total_hits", sumMapValues(domains),
			)
		}
	} else {
		lp.logger.Trace("No DNS stats collected in this batch")
	}
	response := &proto.GetLogDataResponse{UserLogData: make(map[string]*proto.UserLogData)}
	for user, domains := range dnsStats {
		u := &proto.UserLogData{
			ValidIps: make([]string, 0, len(ipUpdates[user])),
			DnsStats: make(map[string]int32, len(domains)),
		}
		for ip := range ipUpdates[user] {
			u.ValidIps = append(u.ValidIps, ip)
		}
		for domain, count := range domains {
			u.DnsStats[domain] = int32(count)
		}
		response.UserLogData[user] = u
	}
	lp.logger.Debug("Prepared log data response",
		"users_count", len(response.UserLogData),
		"total_ips", countTotalIPs(ipUpdates),
	)
	return response, nil
}

func countTotalIPs(m map[string]map[string]struct{}) int {
	total := 0
	for _, ips := range m {
		total += len(ips)
	}
	return total
}

func mapKeysToString(m map[string]int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func sumMapValues(m map[string]int) int {
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}
