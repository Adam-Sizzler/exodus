package sdk

import (
	"context"
	"fmt"
	"time"
)

// Stat is a normalized traffic counter returned by a core Stats API.
type Stat struct {
	Name  string
	Value int64
}

// SysStats is a normalized runtime snapshot returned by core Stats API.
type SysStats struct {
	NumGoroutine uint32
	NumGC        uint32
	Alloc        uint64
	TotalAlloc   uint64
	Sys          uint64
	Mallocs      uint64
	Frees        uint64
	LiveObjects  uint64
	PauseTotalNs uint64
	Uptime       uint32
}

// QueryOptions is a transport-agnostic query request for core stats APIs.
type QueryOptions struct {
	Pattern  string
	Patterns []string
	Regexp   bool
	Reset    bool
}

// StatsService is a minimal SDK contract used by exodus-node.
type StatsService interface {
	QueryStats(ctx context.Context, options QueryOptions) ([]Stat, error)
	GetSysStats(ctx context.Context) (*SysStats, error)
	Close() error
}

// Config defines SDK connection settings.
type Config struct {
	CoreType       string
	Address        string
	Port           int
	DialTimeout    time.Duration
	RequestTimeout time.Duration
}

func (c Config) target() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}
