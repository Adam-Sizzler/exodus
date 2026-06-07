package server

import (
	"context"
	"encoding/json"
	"fmt"

	rpcstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"exodus-node/proto"
)

const statsOnlyMessage = "stats-only node: user management and task APIs are disabled for sing-box"

const (
	taskOperationDeployConfig       = "deploy_config"
	taskOperationSyncSRSLists       = "sync_srs_lists"
	taskOperationNodePluginExecutor = "node_plugin_executor"
)

func (s *NodeServer) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	_ = ctx
	_ = req
	return nil, status.Error(codes.Unimplemented, statsOnlyMessage)
}

func (s *NodeServer) AddUsers(ctx context.Context, req *proto.AddUsersRequest) (*proto.OperationResponse, error) {
	_ = ctx
	_ = req
	return nil, status.Error(codes.Unimplemented, statsOnlyMessage)
}

func (s *NodeServer) DeleteUsers(ctx context.Context, req *proto.DeleteUsersRequest) (*proto.OperationResponse, error) {
	_ = ctx
	_ = req
	return nil, status.Error(codes.Unimplemented, statsOnlyMessage)
}

func (s *NodeServer) SetUserEnabled(ctx context.Context, req *proto.SetUserEnabledRequest) (*proto.OperationResponse, error) {
	_ = ctx
	_ = req
	return nil, status.Error(codes.Unimplemented, statsOnlyMessage)
}

func (s *NodeServer) SubmitTask(ctx context.Context, task *proto.NodeTask) (*rpcstatus.Status, error) {
	_ = ctx
	if task == nil {
		return &rpcstatus.Status{
			Code:    int32(codes.InvalidArgument),
			Message: "task is nil",
		}, nil
	}
	s.Cfg.Logger.Info("SubmitTask received", "task_id", task.TaskId, "operation", task.Operation, "payload_bytes", len(task.Payload))

	switch task.Operation {
	case taskOperationDeployConfig:
		var payload DeployConfigTaskPayload
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			s.Cfg.Logger.Warn("Invalid deploy_config payload", "task_id", task.TaskId, "error", err)
			return &rpcstatus.Status{
				Code:    int32(codes.InvalidArgument),
				Message: fmt.Sprintf("invalid deploy_config payload: %v", err),
			}, nil
		}
		if len(payload.Config) == 0 && len(payload.SingboxConfig) == 0 {
			payload.Config = append(json.RawMessage(nil), task.Payload...)
		}

		summary, err := s.DeployConfig(ctx, payload)
		if err != nil {
			s.Cfg.Logger.Error("Failed to deploy sing-box config", "error", err)
			return &rpcstatus.Status{
				Code:    int32(codes.FailedPrecondition),
				Message: err.Error(),
			}, nil
		}

		message := fmt.Sprintf(
			"success: config_path=%s listen=%s inbounds=%d outbounds=%d users=%d restarted=%t force_restart=%t config_changed=%t haproxy_users_changed=%t srs_downloaded_on_deploy=%t core_ready=%t core_config_valid=%t core_process_before=%s core_process_after=%s",
			summary.ConfigPath,
			summary.Listen,
			summary.Inbounds,
			summary.Outbounds,
			summary.Users,
			summary.Restarted,
			summary.ForceRestart,
			summary.ConfigChanged,
			summary.HaproxyUsersChanged,
			summary.SRSDownloadedOnDeploy,
			summary.CoreReady,
			summary.CoreConfigValid,
			summary.CoreProcessBefore,
			summary.CoreProcessAfter,
		)
		if summary.ReloadError != "" {
			message += fmt.Sprintf(" reload_error=%q", summary.ReloadError)
		}

		return &rpcstatus.Status{
			Code:    int32(codes.OK),
			Message: message,
		}, nil
	case taskOperationSyncSRSLists:
		var payload struct {
			SRSLists []SRSListItem `json:"srs_lists"`
		}
		if err := json.Unmarshal(task.Payload, &payload); err != nil {
			s.Cfg.Logger.Warn("Invalid sync_srs_lists payload", "task_id", task.TaskId, "error", err)
			return &rpcstatus.Status{
				Code:    int32(codes.InvalidArgument),
				Message: fmt.Sprintf("invalid sync_srs_lists payload: %v", err),
			}, nil
		}
		summary, err := s.SyncSRSLists(payload.SRSLists)
		if err != nil {
			return &rpcstatus.Status{
				Code:    int32(codes.FailedPrecondition),
				Message: err.Error(),
			}, nil
		}
		return &rpcstatus.Status{
			Code: int32(codes.OK),
			Message: fmt.Sprintf(
				"success: total=%d configured=%d downloaded=%d failed=%d",
				summary.Total,
				summary.Configured,
				summary.Downloaded,
				summary.Failed,
			),
		}, nil
	case taskOperationNodePluginExecutor:
		accepted, err := ExecuteNodePluginCommand(task.Payload)
		if err != nil {
			s.Cfg.Logger.Warn("Node plugin executor command failed", "task_id", task.TaskId, "error", err)
			return &rpcstatus.Status{
				Code:    int32(codes.FailedPrecondition),
				Message: err.Error(),
			}, nil
		}
		return &rpcstatus.Status{
			Code:    int32(codes.OK),
			Message: fmt.Sprintf("success: accepted=%t", accepted),
		}, nil

	default:
		s.Cfg.Logger.Warn("Unsupported task operation", "task_id", task.TaskId, "operation", task.Operation)
		return &rpcstatus.Status{
			Code:    int32(codes.Unimplemented),
			Message: fmt.Sprintf("%s (operation=%s)", statsOnlyMessage, task.Operation),
		}, nil
	}
}

func (s *NodeServer) GetTaskStatus(ctx context.Context, req *proto.TaskStatusRequest) (*proto.TaskStatusResponse, error) {
	_ = ctx
	taskID := ""
	if req != nil {
		taskID = req.TaskId
	}
	return &proto.TaskStatusResponse{
		TaskId:       taskID,
		Status:       "unsupported",
		ErrorMessage: statsOnlyMessage,
	}, nil
}
