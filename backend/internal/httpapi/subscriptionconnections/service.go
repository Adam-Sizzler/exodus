package subscriptionconnections

import (
	"context"
	"fmt"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/security"
	monitor "exodus/internal/subscriptionnodes"

	"github.com/google/uuid"
)

type SubscriptionConnectionService struct {
	repo *SubscriptionConnectionRepository
	cfg  *config.BackendConfig
}

func NewSubscriptionConnectionService(repo *SubscriptionConnectionRepository, cfg *config.BackendConfig) *SubscriptionConnectionService {
	return &SubscriptionConnectionService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *SubscriptionConnectionService) CreateNode(ctx context.Context, req createNodeRequest) (nodeRecord, error) {
	if req.ProviderUUID != nil && strings.TrimSpace(*req.ProviderUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID)); err != nil {
			return nodeRecord{}, fmt.Errorf("invalid providerUuid")
		}
	}
	schema := normalizeAPISchema(req.APISchema)
	grpcAuthToken := ""
	if req.GRPCAuthToken != nil {
		grpcAuthToken = *req.GRPCAuthToken
	}
	grpcAuthToken, err := security.ResolveGRPCAuthToken(grpcAuthToken)
	if err != nil {
		return nodeRecord{}, err
	}
	subpageConfigUUID := normalizeNullableUUID(req.SubpageConfigUUID)
	if subpageConfigUUID != nil {
		if _, err := uuid.Parse(*subpageConfigUUID); err != nil {
			return nodeRecord{}, fmt.Errorf("invalid subpageConfigUuid")
		}
		exists, err := s.repo.subpageConfigExists(ctx, *subpageConfigUUID)
		if err != nil {
			return nodeRecord{}, err
		}
		if !exists {
			return nodeRecord{}, fmt.Errorf("subpageConfigUuid not found")
		}
	}

	nodeUUID := uuid.NewString()
	now := time.Now().UTC()
	err = s.repo.createNode(ctx, nodeUUID, req, schema, grpcAuthToken, subpageConfigUUID, now)
	if err != nil {
		return nodeRecord{}, err
	}

	monitor.RequestSubNodeSync()
	if subpageConfigUUID != nil {
		subpageConfig, err := s.repo.fetchSubpageConfigRaw(ctx, *subpageConfigUUID)
		if err != nil {
			s.cfg.Logger.Warn("Failed to preload subpage config for created subscription node", "node_uuid", nodeUUID, "subpage_config_uuid", *subpageConfigUUID, "error", err)
		} else {
			monitor.RequestSubNodeSubpageConfigPush(*subpageConfigUUID, subpageConfig, nodeUUID)
		}
	}

	return s.repo.getNodeByUUID(ctx, nodeUUID)
}

func (s *SubscriptionConnectionService) UpdateNode(ctx context.Context, req updateNodeRequest) (nodeRecord, error) {
	if req.SubpageConfigUUID.Set && req.SubpageConfigUUID.Value != nil && strings.TrimSpace(*req.SubpageConfigUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.SubpageConfigUUID.Value)); err != nil {
			return nodeRecord{}, fmt.Errorf("invalid subpageConfigUuid")
		}
	}
	if req.ProviderUUID.Set && req.ProviderUUID.Value != nil && strings.TrimSpace(*req.ProviderUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID.Value)); err != nil {
			return nodeRecord{}, fmt.Errorf("invalid providerUuid")
		}
	}
	var grpcAuthToken *string
	if req.GRPCAuthToken != nil {
		token, err := security.ResolveGRPCAuthToken(*req.GRPCAuthToken)
		if err != nil {
			return nodeRecord{}, err
		}
		grpcAuthToken = &token
	}

	currentNode, err := s.repo.getNodeByUUID(ctx, req.UUID)
	if err != nil {
		return nodeRecord{}, err
	}

	finalSchema := normalizeSubNodeSchema(currentNode.APISchema)
	if req.APISchema != nil {
		finalSchema = normalizeAPISchema(req.APISchema)
	}

	var finalSubpageConfigUUID *string
	if currentNode.SubpageConfigUUID != nil {
		trimmed := strings.TrimSpace(*currentNode.SubpageConfigUUID)
		if trimmed != "" {
			finalSubpageConfigUUID = &trimmed
		}
	}
	if req.SubpageConfigUUID.Set {
		finalSubpageConfigUUID = normalizeNullableUUID(req.SubpageConfigUUID.Value)
	}
	if finalSubpageConfigUUID != nil {
		exists, err := s.repo.subpageConfigExists(ctx, *finalSubpageConfigUUID)
		if err != nil {
			return nodeRecord{}, err
		}
		if !exists {
			return nodeRecord{}, fmt.Errorf("subpageConfigUuid not found")
		}
	}

	clauses := make([]string, 0)
	args := make([]any, 0)
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	if req.Name != nil {
		add("name", strings.TrimSpace(*req.Name))
	}
	if req.Address != nil {
		add("address", strings.TrimSpace(*req.Address))
	}
	if req.PublicDomain.Set {
		if req.PublicDomain.Value == nil || strings.TrimSpace(*req.PublicDomain.Value) == "" {
			clauses = append(clauses, "public_domain = NULL")
		} else {
			add("public_domain", normalizePublicDomain(req.PublicDomain.Value))
		}
	}
	if req.Port != nil {
		add("port", *req.Port)
	}
	if req.APISchema != nil {
		add("api_schema", finalSchema)
	}
	if req.APIPath != nil {
		add("api_path", normalizeAPIPath(req.APIPath))
	}
	if grpcAuthToken != nil {
		add("grpc_auth_token", *grpcAuthToken)
	}
	if req.Tags != nil {
		add("tags", normalizeTags(*req.Tags))
	}
	if req.ProviderUUID.Set {
		if req.ProviderUUID.Value == nil || strings.TrimSpace(*req.ProviderUUID.Value) == "" {
			clauses = append(clauses, "provider_uuid = NULL")
		} else {
			add("provider_uuid", strings.TrimSpace(*req.ProviderUUID.Value))
		}
	}

	err = s.repo.updateNode(ctx, req.UUID, clauses, args, req.SubpageConfigUUID.Set, finalSubpageConfigUUID)
	if err != nil {
		return nodeRecord{}, err
	}

	monitor.RequestSubNodeSync()
	monitor.RequestSubNodeDeploy(req.UUID)
	previousSubpageConfigUUID := normalizeNullableUUID(currentNode.SubpageConfigUUID)
	if previousSubpageConfigUUID != nil {
		if finalSubpageConfigUUID == nil || *finalSubpageConfigUUID != *previousSubpageConfigUUID {
			monitor.RequestSubNodeSubpageConfigPush(*previousSubpageConfigUUID, nil, req.UUID)
		}
	}
	if finalSubpageConfigUUID != nil {
		subpageConfig, err := s.repo.fetchSubpageConfigRaw(ctx, *finalSubpageConfigUUID)
		if err != nil {
			s.cfg.Logger.Warn("Failed to push selected subpage config to subscription node", "node_uuid", req.UUID, "subpage_config_uuid", *finalSubpageConfigUUID, "error", err)
		} else {
			monitor.RequestSubNodeSubpageConfigPush(*finalSubpageConfigUUID, subpageConfig, req.UUID)
		}
	}

	return s.repo.getNodeByUUID(ctx, req.UUID)
}

func (s *SubscriptionConnectionService) DeleteNode(ctx context.Context, nodeUUID string) error {
	err := s.repo.deleteNode(ctx, nodeUUID)
	if err != nil {
		return err
	}
	monitor.RequestSubNodeSync()
	return nil
}

func (s *SubscriptionConnectionService) EnableNode(ctx context.Context, nodeUUID string) (nodeRecord, error) {
	err := s.repo.setNodeDisabled(ctx, nodeUUID, false)
	if err != nil {
		return nodeRecord{}, err
	}
	monitor.RequestSubNodeSync()
	return s.repo.getNodeByUUID(ctx, nodeUUID)
}

func (s *SubscriptionConnectionService) DisableNode(ctx context.Context, nodeUUID string) (nodeRecord, error) {
	err := s.repo.setNodeDisabled(ctx, nodeUUID, true)
	if err != nil {
		return nodeRecord{}, err
	}
	monitor.RequestSubNodeSync()
	return s.repo.getNodeByUUID(ctx, nodeUUID)
}

func (s *SubscriptionConnectionService) RestartNode(ctx context.Context, nodeUUID string) error {
	node, err := s.repo.getNodeByUUID(ctx, nodeUUID)
	if err != nil {
		return err
	}
	if node.IsDisabled {
		return fmt.Errorf("node is disabled")
	}
	monitor.RequestSubNodeDeploy(nodeUUID)
	return nil
}

func (s *SubscriptionConnectionService) ResetNodeTraffic(ctx context.Context, nodeUUID string) error {
	_, err := s.repo.getNodeByUUID(ctx, nodeUUID)
	return err
}

func (s *SubscriptionConnectionService) RestartAllNodes(ctx context.Context) error {
	enabledCount, err := s.repo.countEnabledNodes(ctx)
	if err != nil {
		return err
	}
	if enabledCount == 0 {
		return errNoEnabledNodes
	}
	monitor.RequestSubNodeDeploy()
	return nil
}

func (s *SubscriptionConnectionService) ReorderNodes(ctx context.Context, req reorderNodesRequest) error {
	err := s.repo.reorderNodes(ctx, req.Nodes)
	if err != nil {
		return err
	}
	monitor.RequestSubNodeSync()
	return nil
}

func (s *SubscriptionConnectionService) BulkNodesActions(ctx context.Context, req bulkNodesActionsRequest) error {
	for _, nodeUUID := range req.UUIDs {
		switch req.Action {
		case "ENABLE":
			if err := s.repo.setNodeDisabled(ctx, nodeUUID, false); err != nil {
				return err
			}
		case "DISABLE":
			if err := s.repo.setNodeDisabled(ctx, nodeUUID, true); err != nil {
				return err
			}
		case "RESTART":
			if _, err := s.repo.getNodeByUUID(ctx, nodeUUID); err != nil {
				return err
			}
		case "RESET_TRAFFIC":
			if _, err := s.repo.getNodeByUUID(ctx, nodeUUID); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid bulk action")
		}
	}
	switch req.Action {
	case "RESTART":
		monitor.RequestSubNodeDeploy(req.UUIDs...)
	default:
		monitor.RequestSubNodeSync()
	}
	return nil
}

func (s *SubscriptionConnectionService) BulkProfileModification(ctx context.Context, req bulkProfileModificationRequest) error {
	// Original code only performed check validateUUIDs, no DB mutations
	return nil
}
