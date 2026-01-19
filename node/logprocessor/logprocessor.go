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
	"v2ray-stat/node/config"
	"v2ray-stat/proto"

	"github.com/coreos/go-systemd/v22/sdjournal"
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
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	offset, err := f.Seek(0, 2)
	if err != nil {
		f.Close()
		return nil, err
	}

	fs := &FileLogSource{
		path:   path,
		file:   f,
		offset: offset,
		logger: logger,
	}

	go fs.runDailyCleanup()

	return fs, nil
}

func (fs *FileLogSource) FetchNewLines() ([]string, error) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	_, err := fs.file.Seek(fs.offset, 0)
	if err != nil {
		return nil, err
	}

	var lines []string
	scanner := bufio.NewScanner(fs.file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	pos, _ := fs.file.Seek(0, 1)
	fs.offset = pos
	return lines, nil
}

func (fs *FileLogSource) Close() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.file.Close()
}

func (fs *FileLogSource) runDailyCleanup() {
	dailyTicker := time.NewTicker(24 * time.Hour)
	defer dailyTicker.Stop()

	for range dailyTicker.C {
		fs.mu.Lock()
		if err := fs.file.Close(); err != nil {
			fs.logger.Error("Error closing log file", "file", fs.path, "error", err)
		}
		f, err := os.OpenFile(fs.path, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			fs.logger.Error("Error reopening log file after truncation", "file", fs.path, "error", err)
			fs.mu.Unlock()
			return
		}
		fs.file = f
		fs.offset = 0
		fs.mu.Unlock()
		fs.logger.Info("Log file successfully truncated", "file", fs.path)
	}
}

type JournalLogSource struct {
	journal *sdjournal.Journal
	logger  *logger.Logger
	buffer  []string
	mu      sync.Mutex
	done    chan struct{}
}

func NewJournalLogSource(serviceName string, logger *logger.Logger) (*JournalLogSource, error) {
	j, err := sdjournal.NewJournal()
	if err != nil {
		return nil, fmt.Errorf("failed to create journal reader: %v", err)
	}

	targetUnit := serviceName
	if !strings.HasSuffix(targetUnit, ".service") {
		targetUnit += ".service"
	}

	if err := j.AddMatch("_SYSTEMD_UNIT=" + targetUnit); err != nil {
		j.Close()
		return nil, fmt.Errorf("failed to add journal match: %v", err)
	}

	if err := j.SeekTail(); err != nil {
		j.Close()
		return nil, fmt.Errorf("failed to seek tail: %v", err)
	}

	_, err = j.Next()

	js := &JournalLogSource{
		journal: j,
		logger:  logger,
		buffer:  make([]string, 0),
		done:    make(chan struct{}),
	}

	go js.readLoop()

	logger.Info("Started native journal log source", "service", targetUnit)
	return js, nil
}

func (js *JournalLogSource) readLoop() {
	defer close(js.done)
	defer js.journal.Close()

	for {
		c, err := js.journal.Next()
		if err != nil {
			js.logger.Error("Error iterating journal", "error", err)
			return
		}

		if c == 0 {
			js.journal.Wait(2 * time.Second)
			continue
		}

		msg, err := js.journal.GetDataValue("MESSAGE")
		if err != nil {
			continue
		}

		js.mu.Lock()
		js.buffer = append(js.buffer, msg)
		js.mu.Unlock()
	}
}

func (js *JournalLogSource) FetchNewLines() ([]string, error) {
	js.mu.Lock()
	defer js.mu.Unlock()

	select {
	case <-js.done:
		return nil, fmt.Errorf("journal reader closed unexpected")
	default:
	}

	if len(js.buffer) == 0 {
		return []string{}, nil
	}

	lines := make([]string, len(js.buffer))
	copy(lines, js.buffer)
	js.buffer = js.buffer[:0]

	return lines, nil
}

func (js *JournalLogSource) Close() error {
	return nil
}

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
		source, err = NewJournalLogSource(cfg.Core.LogServiceName, cfg.Logger)
	} else {
		source, err = NewFileLogSource(cfg.Core.AccessLog, cfg.Logger)
	}

	if err != nil {
		return nil, err
	}

	compiledRegex, err := regexp.Compile(cfg.Core.AccessLogRegex)
	if err != nil {
		return nil, fmt.Errorf("invalid regex: %v", err)
	}

	return &LogProcessor{
		cfg:    cfg,
		source: source,
		logger: cfg.Logger,
		regex:  compiledRegex,
	}, nil
}

func ProcessLogLine(line string, dnsStats map[string]map[string]int, cfg *config.NodeConfig, re *regexp.Regexp) (string, []string, bool) {
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return "", nil, false
	}

	getValue := func(index int) string {
		if index > 0 && index < len(matches) {
			return strings.TrimSpace(matches[index])
		}
		return ""
	}

	user := getValue(cfg.Core.LogMap.User)
	ip := getValue(cfg.Core.LogMap.IP)
	domain := getValue(cfg.Core.LogMap.Domain)

	if user == "" {
		return "", nil, false
	}

	var validIPs []string
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
		lp.logger.Error("Error fetching log lines", "error", err)
		return nil, err
	}

	dnsStats := make(map[string]map[string]int)
	ipUpdates := make(map[string]map[string]struct{})

	for _, line := range lines {
		user, validIPs, ok := ProcessLogLine(line, dnsStats, lp.cfg, lp.regex)
		if ok {
			if ipUpdates[user] == nil {
				ipUpdates[user] = make(map[string]struct{})
			}
			for _, ip := range validIPs {
				ipUpdates[user][ip] = struct{}{}
			}
		}
	}

	response := &proto.GetLogDataResponse{UserLogData: make(map[string]*proto.UserLogData)}
	for user, domains := range dnsStats {
		response.UserLogData[user] = &proto.UserLogData{
			ValidIps: make([]string, 0, len(ipUpdates[user])),
			DnsStats: make(map[string]int32),
		}
		for ip := range ipUpdates[user] {
			response.UserLogData[user].ValidIps = append(response.UserLogData[user].ValidIps, ip)
		}
		for domain, count := range domains {
			response.UserLogData[user].DnsStats[domain] = int32(count)
		}
	}
	return response, nil
}
