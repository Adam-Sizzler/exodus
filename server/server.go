package server

import (
	"context"
	"fmt"

	"exodus-node/api"
	"exodus-node/config"
	"exodus-node/proto"

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
	if cfg == nil {
		return nil, fmt.Errorf("node config is nil")
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("node logger is nil")
	}

	apiService, err := api.NewService(cfg)
	if err != nil {
		cfg.LoggerFor("StatsService").Error("Failed to initialize core API service", "error", err)
		return nil, fmt.Errorf("initialize core API service: %w", err)
	}

	nodeServer := &NodeServer{
		Cfg:        cfg,
		apiService: apiService,
	}
	cfg.LoggerFor("Supervisor").Info("[OK] Supervisord initialized")
	cfg.LoggerFor("NetworkStatsService").Info("Network stats polling started (interval: 1000ms, default: " + detectDefaultNetworkInterfaceForLogs() + ")")
	cfg.LoggerFor("NftService").Info("[PLUGIN] Ingress Filter: available")
	cfg.LoggerFor("NftService").Info("[PLUGIN] Egress Filter: available")
	cfg.LoggerFor("HAProxyService").Info("[PLUGIN] HAProxy Runtime API: available")
	logNftablesAvailability(cfg.LoggerFor("HandlerService"))
	nodeServer.startSRSAutoUpdater()

	return nodeServer, nil
}

func (s *NodeServer) Close() error {
	if s == nil || s.apiService == nil {
		return nil
	}
	return s.apiService.Close()
}

// GetApiStats retrieves API statistics from the node.
func (s *NodeServer) GetApiStats(ctx context.Context, req *proto.GetApiStatsRequest) (*proto.GetApiStatsResponse, error) {
	s.Cfg.LoggerFor("StatsService").Debug("Received GetApiStats request")
	apiData, err := s.apiService.GetApiResponse(ctx)
	if err != nil {
		s.Cfg.LoggerFor("StatsService").Error("Failed to get API response", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get API response: %v", err)
	}

	response := &proto.GetApiStatsResponse{}
	for _, stat := range apiData.Stat {
		response.Stats = append(response.Stats, &proto.Stat{
			Name:  stat.Name,
			Value: stat.Value,
		})
	}

	s.Cfg.LoggerFor("StatsService").Debug("Returning API stats", "stats_count", len(response.Stats))
	return response, nil
}

// GetLogData returns log data. Log processor is removed in this simplified node mode.
func (s *NodeServer) GetLogData(ctx context.Context, req *proto.GetLogDataRequest) (*proto.GetLogDataResponse, error) {
	_ = ctx
	_ = req
	return &proto.GetLogDataResponse{UserLogData: map[string]*proto.UserLogData{}}, nil
}
