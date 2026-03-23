package sdk

import (
	"context"
	"time"

	singboxapi "cerberus-node/sdk/singboxapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type singboxStatsService struct {
	conn           *grpc.ClientConn
	client         singboxapi.StatsServiceClient
	requestTimeout time.Duration
}

func newSingboxStatsService(cfg Config) (*singboxStatsService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		cfg.target(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &singboxStatsService{
		conn:           conn,
		client:         singboxapi.NewStatsServiceClient(conn),
		requestTimeout: cfg.RequestTimeout,
	}, nil
}

func (s *singboxStatsService) QueryStats(ctx context.Context, options QueryOptions) ([]Stat, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	resp, err := s.client.QueryStats(ctx, &singboxapi.QueryStatsRequest{
		Pattern:  options.Pattern,
		Reset_:   options.Reset,
		Patterns: append([]string(nil), options.Patterns...),
		Regexp:   options.Regexp,
	})
	if err != nil {
		return nil, err
	}

	result := make([]Stat, 0, len(resp.GetStat()))
	for _, stat := range resp.GetStat() {
		result = append(result, Stat{
			Name:  stat.GetName(),
			Value: stat.GetValue(),
		})
	}
	return result, nil
}

func (s *singboxStatsService) GetSysStats(ctx context.Context) (*SysStats, error) {
	ctx, cancel := s.withTimeout(ctx)
	defer cancel()

	resp, err := s.client.GetSysStats(ctx, &singboxapi.SysStatsRequest{})
	if err != nil {
		return nil, err
	}

	return &SysStats{
		NumGoroutine: resp.GetNumGoroutine(),
		NumGC:        resp.GetNumGC(),
		Alloc:        resp.GetAlloc(),
		TotalAlloc:   resp.GetTotalAlloc(),
		Sys:          resp.GetSys(),
		Mallocs:      resp.GetMallocs(),
		Frees:        resp.GetFrees(),
		LiveObjects:  resp.GetLiveObjects(),
		PauseTotalNs: resp.GetPauseTotalNs(),
		Uptime:       resp.GetUptime(),
	}, nil
}

func (s *singboxStatsService) Close() error {
	if s == nil || s.conn == nil {
		return nil
	}
	return s.conn.Close()
}

func (s *singboxStatsService) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.requestTimeout <= 0 {
		return ctx, func() {}
	}
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, s.requestTimeout)
}
