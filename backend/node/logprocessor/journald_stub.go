//go:build !journald

package logprocessor

import (
	"fmt"
	"v2ray-stat/logger"
)

type JournalLogSource struct{}

func (js *JournalLogSource) FetchNewLines() ([]string, error) {
	return nil, fmt.Errorf("journald support disabled in this build")
}

func (js *JournalLogSource) Close() error {
	return nil
}

func NewJournalLogSource(serviceName string, logger *logger.Logger) (LogSource, error) {
	logger.Error("Journald log source is not supported in this build")
	return nil, fmt.Errorf("build tag 'journald' is missing")
}
