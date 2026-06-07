package sdk

import (
	"context"
	"time"

	singboxapi "exodus-node/sdk/singboxapi"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type singboxStatsService struct {
	conn           *grpc.ClientConn
	client         singboxapi.StatsServiceClient
	requestTimeout time.Duration
}

func newSingboxStatsService(cfg Config) (*singboxStatsService, error) {
	// Do not block node startup on the local core API. sing-box can be down while
	// the node is still useful for receiving a fixed config from the panel. The
	// connection will be established lazily by gRPC, and per-request calls below
	// will report temporary core errors without killing the node process.
	conn, err := grpc.DialContext(
		context.Background(),
		cfg.target(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
