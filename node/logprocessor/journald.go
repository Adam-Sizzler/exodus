//go:build journald
// +build journald

package logprocessor

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"v2ray-stat/logger"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

type JournalLogSource struct {
	journal *sdjournal.Journal
	logger  *logger.Logger
	buffer  []string
	mu      sync.Mutex
	done    chan struct{}
	wg      sync.WaitGroup
}

func NewJournalLogSource(serviceName string, logger *logger.Logger) (LogSource, error) {
	j, err := sdjournal.NewJournal()
	if err != nil {
		logger.Error("Failed to open journald", "error", err)
		return nil, fmt.Errorf("failed to open journal: %w", err)
	}

	targetUnit := serviceName
	if !strings.HasSuffix(targetUnit, ".service") {
		targetUnit += ".service"
	}

	if err := j.AddMatch("_SYSTEMD_UNIT=" + targetUnit); err != nil {
		j.Close()
		logger.Error("Failed to add match filter for unit", "unit", targetUnit, "error", err)
		return nil, fmt.Errorf("failed to add match: %w", err)
	}

	if err := j.SeekTail(); err != nil {
		j.Close()
		logger.Error("Failed to seek tail in journal", "error", err)
		return nil, fmt.Errorf("seek tail failed: %w", err)
	}

	_, _ = j.Next()

	js := &JournalLogSource{
		journal: j,
		logger:  logger,
		buffer:  make([]string, 0, 128),
		done:    make(chan struct{}),
	}

	go js.readLoop()
	logger.Info("Journald log source started", "unit", targetUnit)

	return js, nil
}

func (js *JournalLogSource) readLoop() {
	defer close(js.done)
	defer js.journal.Close()

	for {
		select {
		case <-js.done:
			js.logger.Debug("Journal read loop stopped")
			return
		default:
			n, err := js.journal.Next()
			if err != nil {
				js.logger.Error("Error advancing journal iterator", "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if n == 0 {
				js.journal.Wait(1500 * time.Millisecond)
				js.logger.Trace("Waiting for new journal entries")
				continue
			}

			msg, err := js.journal.GetDataValue("MESSAGE")
			if err != nil || msg == "" {
				continue
			}

			js.mu.Lock()
			js.buffer = append(js.buffer, msg)
			js.mu.Unlock()

			js.logger.Trace("Added journal line to buffer", "buffer_size", len(js.buffer))
		}
	}
}

func (js *JournalLogSource) FetchNewLines() ([]string, error) {
	js.mu.Lock()
	defer js.mu.Unlock()

	select {
	case <-js.done:
		js.logger.Warn("Attempt to read from closed journal source")
		return nil, fmt.Errorf("journal source is closed")
	default:
	}

	if len(js.buffer) == 0 {
		js.logger.Trace("Journal buffer is empty")
		return []string{}, nil
	}

	lines := make([]string, len(js.buffer))
	copy(lines, js.buffer)

	js.logger.Debug("Extracted lines from journal buffer", "count", len(lines))
	js.buffer = js.buffer[:0]

	return lines, nil
}

func (js *JournalLogSource) Close() error {
	close(js.done)
	js.wg.Wait()

	js.logger.Info("Journald source closed")
	return nil
}
