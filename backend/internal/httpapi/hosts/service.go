package hosts

import (
	"context"

	"exodus/internal/config"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

type HostService struct {
	repo *HostRepository
	cfg  *config.BackendConfig
}

func NewHostService(repo *HostRepository, cfg *config.BackendConfig) *HostService {
	return &HostService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *HostService) CreateHost(ctx context.Context, req HostCreateRequestAPI) (hostRecord, error) {
	if err := s.repo.ensureConfigProfileInbound(ctx, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
		return hostRecord{}, err
	}
	if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID != "" {
		if err := s.repo.ensureXrayJSONTemplate(ctx, *req.XrayJSONTemplateUUID); err != nil {
			return hostRecord{}, err
		}
	}

	xhttpExtra, err := normalizeJSONValue(req.XHTTPExtraParams, true)
	if err != nil {
		return hostRecord{}, err
	}
	mux, err := normalizeJSONValue(req.MuxParams, true)
	if err != nil {
		return hostRecord{}, err
	}
	sockopt, err := normalizeJSONValue(req.SockoptParams, true)
	if err != nil {
		return hostRecord{}, err
	}
	finalMask, err := normalizeJSONValue(req.FinalMask, true)
	if err != nil {
		return hostRecord{}, err
	}

	hostUUID := uuid.NewString()
	err = s.repo.createHost(ctx, hostUUID, req, xhttpExtra, mux, sockopt, finalMask)
	if err != nil {
		return hostRecord{}, err
	}

	if len(req.Nodes) > 0 {
		monitor.RequestNodeDeploy(true, req.Nodes...)
	}

	return s.repo.getHostByUUID(ctx, hostUUID)
}

func (s *HostService) UpdateHost(ctx context.Context, req HostUpdateRequestAPI) (hostRecord, error) {
	if req.Inbound != nil {
		if err := s.repo.ensureConfigProfileInbound(ctx, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
			return hostRecord{}, err
		}
	}
	if req.XrayJSONTemplateUUID.Set && req.XrayJSONTemplateUUID.Value != nil && *req.XrayJSONTemplateUUID.Value != "" {
		if err := s.repo.ensureXrayJSONTemplate(ctx, *req.XrayJSONTemplateUUID.Value); err != nil {
			return hostRecord{}, err
		}
	}

	clauses, args, err := buildHostUpdateClauses(req.hostUpdateFields)
	if err != nil {
		return hostRecord{}, err
	}

	// Fetch old nodes to determine who to deploy to (merge old & new)
	oldNodesMap, _ := s.repo.getHostNodes(ctx, []string{req.UUID})
	oldNodes := oldNodesMap[req.UUID]

	squads := req.ExcludedInternalSquads
	if req.InternalSquads != nil {
		squads = req.InternalSquads.Squads
	}

	err = s.repo.updateHost(ctx, req.UUID, clauses, args, req.Nodes, squads)
	if err != nil {
		return hostRecord{}, err
	}

	// Trigger node deploys
	deploys := dedupeStrings(append(oldNodes, req.Nodes...))
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}

	return s.repo.getHostByUUID(ctx, req.UUID)
}

func (s *HostService) DeleteHost(ctx context.Context, hostUUID string) error {
	// Fetch old nodes to redeploy
	oldNodesMap, _ := s.repo.getHostNodes(ctx, []string{hostUUID})
	oldNodes := oldNodesMap[hostUUID]

	err := s.repo.deleteHost(ctx, hostUUID)
	if err != nil {
		return err
	}

	if len(oldNodes) > 0 {
		monitor.RequestNodeDeploy(true, oldNodes...)
	}
	return nil
}

func (s *HostService) ReorderHosts(ctx context.Context, req reorderHostsRequest) error {
	return s.repo.reorderHosts(ctx, req.Hosts)
}

func (s *HostService) BulkEnableHosts(ctx context.Context, req bulkUUIDsRequest) error {
	// Fetch nodes linked to these hosts to deploy updates
	nodesMap, _ := s.repo.getHostNodes(ctx, req.UUIDs)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}

	err := s.repo.bulkUpdateHostsEnabled(ctx, req.UUIDs, true)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}

func (s *HostService) BulkDisableHosts(ctx context.Context, req bulkUUIDsRequest) error {
	// Fetch nodes linked to these hosts to deploy updates
	nodesMap, _ := s.repo.getHostNodes(ctx, req.UUIDs)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}

	err := s.repo.bulkUpdateHostsEnabled(ctx, req.UUIDs, false)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}

func (s *HostService) BulkDeleteHosts(ctx context.Context, req bulkUUIDsRequest) error {
	// Fetch nodes linked to these hosts to deploy updates
	nodesMap, _ := s.repo.getHostNodes(ctx, req.UUIDs)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}

	err := s.repo.bulkDeleteHosts(ctx, req.UUIDs)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}

func (s *HostService) BulkSetInbound(ctx context.Context, req setInboundRequest) error {
	if err := s.repo.ensureConfigProfileInbound(ctx, req.ConfigProfileUUID, req.ConfigProfileInboundUUID); err != nil {
		return err
	}

	// Fetch nodes linked to these hosts to deploy updates
	nodesMap, _ := s.repo.getHostNodes(ctx, req.UUIDs)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}

	err := s.repo.bulkSetInbound(ctx, req.UUIDs, req.ConfigProfileUUID, req.ConfigProfileInboundUUID)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}

func (s *HostService) BulkSetPort(ctx context.Context, req setPortRequest) error {
	// Fetch nodes linked to these hosts to deploy updates
	nodesMap, _ := s.repo.getHostNodes(ctx, req.UUIDs)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}

	err := s.repo.bulkSetPort(ctx, req.UUIDs, req.Port)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}

func (s *HostService) BulkUpdateHosts(ctx context.Context, req HostBulkUpdateRequestAPI) error {
	if req.Inbound != nil {
		if err := s.repo.ensureConfigProfileInbound(ctx, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
			return err
		}
	}
	if req.XrayJSONTemplateUUID.Set && req.XrayJSONTemplateUUID.Value != nil && *req.XrayJSONTemplateUUID.Value != "" {
		if err := s.repo.ensureXrayJSONTemplate(ctx, *req.XrayJSONTemplateUUID.Value); err != nil {
			return err
		}
	}

	clauses, args, err := buildHostUpdateClauses(req.hostUpdateFields)
	if err != nil {
		return err
	}

	// Fetch nodes linked to these hosts to deploy updates (merge old & new)
	nodesMap, _ := s.repo.getHostNodes(ctx, req.Uuids)
	var deploys []string
	for _, hostNodes := range nodesMap {
		deploys = append(deploys, hostNodes...)
	}
	if req.Nodes != nil {
		deploys = append(deploys, req.Nodes...)
	}

	squads := req.ExcludedInternalSquads
	if req.InternalSquads != nil {
		squads = req.InternalSquads.Squads
	}

	err = s.repo.bulkUpdateHosts(ctx, req.Uuids, clauses, args, req.Nodes, squads)
	if err != nil {
		return err
	}

	deploys = dedupeStrings(deploys)
	if len(deploys) > 0 {
		monitor.RequestNodeDeploy(true, deploys...)
	}
	return nil
}
