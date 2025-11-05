package server

import (
	"context"
	"fmt"
	"io"
	"time"

	"v2ray-stat/node/api"
	"v2ray-stat/node/config"
	"v2ray-stat/node/logprocessor"
	"v2ray-stat/proto"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NodeServer implements the NodeService gRPC service.
type NodeServer struct {
	proto.UnimplementedNodeServiceServer
	Cfg          *config.NodeConfig
	logProcessor *logprocessor.LogProcessor
	taskManager  *TaskManager
}

// NewNodeServer creates a new NodeServer instance.
func NewNodeServer(cfg *config.NodeConfig) (*NodeServer, error) {
	processor, err := logprocessor.NewLogProcessor(cfg)
	if err != nil {
		cfg.Logger.Error("Failed to create log processor", "error", err)
		return nil, fmt.Errorf("create log processor: %w", err)
	}
	server := &NodeServer{
		Cfg:          cfg,
		logProcessor: processor,
		taskManager:  NewTaskManager(cfg),
	}
	return server, nil
}

// GetApiStats retrieves API statistics from the node.
func (s *NodeServer) GetApiStats(ctx context.Context, req *proto.GetApiStatsRequest) (*proto.GetApiStatsResponse, error) {
	s.Cfg.Logger.Debug("Received GetApiStats request")
	apiData, err := api.GetApiResponse(s.Cfg)
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

// GetLogData retrieves processed log data from the node.
func (s *NodeServer) GetLogData(ctx context.Context, req *proto.GetLogDataRequest) (*proto.GetLogDataResponse, error) {
	response, err := s.logProcessor.ReadNewLines()
	if err != nil {
		s.Cfg.Logger.Error("Failed to read log data", "error", err)
		return nil, status.Errorf(codes.Internal, "read log data: %v", err)
	}
	s.Cfg.Logger.Debug("Returning log data via gRPC", "users", len(response.UserLogData))
	return response, nil
}

// StreamNodeData handles bidirectional streaming for node data.
func (s *NodeServer) StreamNodeData(stream proto.NodeService_StreamNodeDataServer) error {
	var ticker *time.Ticker
	defer func() {
		if ticker != nil {
			ticker.Stop()
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			s.Cfg.Logger.Info("Stream closed by client or context canceled", "error", stream.Context().Err())
			return stream.Context().Err()
		default:
			req, err := stream.Recv()
			if err == io.EOF {
				s.Cfg.Logger.Info("Stream closed by client")
				return nil
			}
			if err != nil {
				s.Cfg.Logger.Error("Failed to receive request", "error", err)
				return err
			}

			switch req.Request.(type) {
			case *proto.NodeDataRequest_Config:
				interval := req.GetConfig().IntervalSeconds
				if interval <= 0 {
					s.Cfg.Logger.Error("Invalid interval received", "interval", interval)
					return status.Errorf(codes.InvalidArgument, "invalid interval: %d", interval)
				}
				if ticker != nil {
					ticker.Stop()
				}
				ticker = time.NewTicker(time.Duration(interval) * time.Second)
				s.Cfg.Logger.Info("Received stream config", "interval_seconds", interval)

				// Start sending stats and log data periodically
				go func() {
					for {
						select {
						case <-ticker.C:
							// Send stats
							stats, err := s.GetApiStats(stream.Context(), &proto.GetApiStatsRequest{})
							if err != nil {
								s.Cfg.Logger.Error("Failed to get stats", "error", err)
								continue
							}
							if err := stream.Send(&proto.NodeDataResponse{
								Response: &proto.NodeDataResponse_Stats{Stats: stats},
							}); err != nil {
								s.Cfg.Logger.Error("Failed to send stats", "error", err)
								return
							}

							// Send log data
							logData, err := s.GetLogData(stream.Context(), &proto.GetLogDataRequest{})
							if err != nil {
								s.Cfg.Logger.Error("Failed to get log data", "error", err)
								continue
							}
							if err := stream.Send(&proto.NodeDataResponse{
								Response: &proto.NodeDataResponse_LogData{LogData: logData},
							}); err != nil {
								s.Cfg.Logger.Error("Failed to send log data", "error", err)
								return
							}
						case <-stream.Context().Done():
							s.Cfg.Logger.Debug("Stream context done, stopping ticker")
							return
						}
					}
				}()

			case *proto.NodeDataRequest_ListUsers:
				users, err := s.ListUsers(stream.Context(), req.GetListUsers())
				if err != nil {
					s.Cfg.Logger.Error("Failed to get users", "error", err)
					return err
				}
				if err := stream.Send(&proto.NodeDataResponse{
					Response: &proto.NodeDataResponse_Users{Users: users},
				}); err != nil {
					s.Cfg.Logger.Error("Failed to send users", "error", err)
					return err
				}
				s.Cfg.Logger.Debug("Sent users response", "user_count", len(users.Users))

			default:
				s.Cfg.Logger.Error("Unknown request type")
				return status.Errorf(codes.InvalidArgument, "unknown request type")
			}
		}
	}
}

// SubmitTask submits a task for async processing.
func (s *NodeServer) SubmitTask(ctx context.Context, task *proto.NodeTask) (*rpcstatus.Status, error) {
	s.Cfg.Logger.Info("Received task submission", "task_id", task.TaskId, "operation", task.Operation)

	// Store task as pending
	s.taskManager.StoreTask(task.TaskId, TaskStatusPending, "", nil, nil)

	// Process task asynchronously
	go s.taskManager.ProcessTask(ctx, task, s)

	s.Cfg.Logger.Debug("Task submitted successfully", "task_id", task.TaskId)
	return &rpcstatus.Status{Code: int32(codes.OK), Message: "task submitted"}, nil
}

// GetTaskStatus retrieves the status of a task.
func (s *NodeServer) GetTaskStatus(ctx context.Context, req *proto.TaskStatusRequest) (*proto.TaskStatusResponse, error) {
	s.Cfg.Logger.Debug("Received task status request", "task_id", req.TaskId)

	result, exists := s.taskManager.GetTask(req.TaskId)
	if !exists {
		s.Cfg.Logger.Warn("Task not found", "task_id", req.TaskId)
		return &proto.TaskStatusResponse{
			TaskId: req.TaskId,
			Status: string(TaskStatusError),
			ErrorMessage: "task not found",
		}, nil
	}

	s.Cfg.Logger.Debug("Returning task status", "task_id", req.TaskId, "status", result.Status)
	return &proto.TaskStatusResponse{
		TaskId:       result.TaskID,
		Status:       string(result.Status),
		ErrorMessage: result.ErrorMessage,
		Users:        result.Users,
		Credentials:  result.Credentials,
	}, nil
}
