package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"cerberus-node/api"
	"cerberus-node/config"
	"cerberus-node/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NodeServer implements the NodeService gRPC service.
type NodeServer struct {
	proto.UnimplementedNodeServiceServer
	Cfg        *config.NodeConfig
	apiService *api.Service
}

// NewNodeServer creates a new NodeServer instance.
func NewNodeServer(cfg *config.NodeConfig) (*NodeServer, error) {
	apiService, err := api.NewService(cfg)
	if err != nil {
		cfg.Logger.Error("Failed to initialize core API service", "error", err)
		return nil, fmt.Errorf("initialize core API service: %w", err)
	}

	nodeServer := &NodeServer{
		Cfg:        cfg,
		apiService: apiService,
	}
	nodeServer.startSRSAutoUpdater()

	return nodeServer, nil
}

// GetApiStats retrieves API statistics from the node.
func (s *NodeServer) GetApiStats(ctx context.Context, req *proto.GetApiStatsRequest) (*proto.GetApiStatsResponse, error) {
	s.Cfg.Logger.Debug("Received GetApiStats request")
	apiData, err := s.apiService.GetApiResponse(ctx)
	if err != nil {
		s.Cfg.Logger.Error("Failed to get API response", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get API response: %v", err)
	}

	response := &proto.GetApiStatsResponse{}
	for _, stat := range apiData.Stat {
		response.Stats = append(response.Stats, &proto.Stat{
			Name:  stat.Name,
			Value: stat.Value,
		})
	}

	s.Cfg.Logger.Debug("Returning API stats", "stats_count", len(response.Stats))
	return response, nil
}

// GetLogData returns log data. Log processor is removed in this simplified node mode.
func (s *NodeServer) GetLogData(ctx context.Context, req *proto.GetLogDataRequest) (*proto.GetLogDataResponse, error) {
	_ = ctx
	_ = req
	return &proto.GetLogDataResponse{UserLogData: map[string]*proto.UserLogData{}}, nil
}

// StreamNodeData handles bidirectional streaming for node data.
func (s *NodeServer) StreamNodeData(stream proto.NodeService_StreamNodeDataServer) error {
	const defaultIntervalSeconds = 10

	reqCh := make(chan *proto.NodeDataRequest)
	recvErrCh := make(chan error, 1)
	sendErrCh := make(chan error, 1)
	updateIntervalCh := make(chan time.Duration, 1)

	// Receiver goroutine: reads client messages (interval updates / keepalive semantics).
	go func() {
		defer close(reqCh)
		for {
			req, err := stream.Recv()
			if err != nil {
				recvErrCh <- err
				return
			}
			reqCh <- req
		}
	}()

	// Sender goroutine: starts pushing stats immediately after stream establishment.
	go func() {
		ticker := time.NewTicker(defaultIntervalSeconds * time.Second)
		defer ticker.Stop()
		consecutiveStatsErrors := 0

		for {
			select {
			case <-stream.Context().Done():
				return
			case next := <-updateIntervalCh:
				ticker.Stop()
				ticker = time.NewTicker(next)
			case <-ticker.C:
				stats, err := s.GetApiStats(stream.Context(), &proto.GetApiStatsRequest{})
				if err != nil {
					consecutiveStatsErrors++
					s.Cfg.Logger.Error("Failed to collect stats for stream", "error", err)
					if consecutiveStatsErrors >= 3 {
						diagnostic := s.diagnoseCoreFailure(stream.Context(), err)
						sendErrCh <- status.Error(codes.Internal, diagnostic)
						return
					}
					continue
				}
				consecutiveStatsErrors = 0
				if err := stream.Send(&proto.NodeDataResponse{
					Response: &proto.NodeDataResponse_Stats{Stats: stats},
				}); err != nil {
					sendErrCh <- err
					return
				}
			}
		}
	}()

	s.Cfg.Logger.Info("Stream started, auto-push enabled", "interval_seconds", defaultIntervalSeconds)

	for {
		select {
		case <-stream.Context().Done():
			s.Cfg.Logger.Info("Stream context canceled", "error", stream.Context().Err())
			return stream.Context().Err()
		case err := <-sendErrCh:
			s.Cfg.Logger.Error("Failed to send stream data", "error", err)
			return err
		case err := <-recvErrCh:
			if err == io.EOF {
				s.Cfg.Logger.Info("Stream closed by client")
				return nil
			}
			s.Cfg.Logger.Error("Failed to receive stream request", "error", err)
			return err
		case req, ok := <-reqCh:
			if !ok || req == nil {
				continue
			}

			switch req.Request.(type) {
			case *proto.NodeDataRequest_Config:
				interval := req.GetConfig().IntervalSeconds
				if interval <= 0 {
					return status.Errorf(codes.InvalidArgument, "invalid interval: %d", interval)
				}
				select {
				case updateIntervalCh <- time.Duration(interval) * time.Second:
				default:
					<-updateIntervalCh
					updateIntervalCh <- time.Duration(interval) * time.Second
				}
				s.Cfg.Logger.Info("Stream interval updated", "interval_seconds", interval)
			case *proto.NodeDataRequest_ListUsers:
				return status.Error(codes.Unimplemented, statsOnlyMessage)
			default:
				return status.Errorf(codes.InvalidArgument, "unknown request type")
			}
		}
	}
}
