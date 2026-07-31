package squads

import (
	"context"
	"fmt"
	"strings"

	"exodus/internal/config"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

type SquadService struct {
	repo *SquadRepository
	cfg  *config.BackendConfig
}

func NewSquadService(repo *SquadRepository, cfg *config.BackendConfig) *SquadService {
	return &SquadService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *SquadService) CreateSquad(ctx context.Context, req InternalSquadCreateRequest) (InternalSquadAPI, error) {
	squadUUID := uuid.NewString()
	err := s.repo.createSquad(ctx, squadUUID, req)
	if err != nil {
		return InternalSquadAPI{}, err
	}

	squad, err := s.repo.getSquadByUUID(ctx, squadUUID)
	if err != nil {
		return InternalSquadAPI{}, err
	}

	monitor.RequestNodeDeploy(true)

	membersCount, _ := s.repo.getSquadMembersCount(ctx, squadUUID)
	inbounds, _ := s.repo.getSquadInbounds(ctx, squadUUID)
	return buildInternalSquadResponse(squad, membersCount, inbounds), nil
}

func (s *SquadService) UpdateSquad(ctx context.Context, req InternalSquadUpdateRequest) (InternalSquadAPI, error) {
	clauses := make([]string, 0)
	args := make([]any, 0)
	idx := 1
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column, idx))
		args = append(args, value)
		idx++
	}

	if req.Name != nil {
		add("name", strings.TrimSpace(*req.Name))
	}
	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}

	err := s.repo.updateSquad(ctx, req.UUID, clauses, args, req.Inbounds)
	if err != nil {
		return InternalSquadAPI{}, err
	}

	squad, err := s.repo.getSquadByUUID(ctx, req.UUID)
	if err != nil {
		return InternalSquadAPI{}, err
	}

	monitor.RequestNodeDeploy(true)

	membersCount, _ := s.repo.getSquadMembersCount(ctx, req.UUID)
	inbounds, _ := s.repo.getSquadInbounds(ctx, req.UUID)
	return buildInternalSquadResponse(squad, membersCount, inbounds), nil
}

func (s *SquadService) DeleteSquad(ctx context.Context, squadUUID string) (string, error) {
	name, err := s.repo.deleteSquad(ctx, squadUUID)
	if err != nil {
		return "", err
	}
	monitor.RequestNodeDeploy(true)
	return name, nil
}

func (s *SquadService) ReorderSquads(ctx context.Context, req reorderSquadsRequest) error {
	return s.repo.reorderSquads(ctx, req.Squads)
}

func (s *SquadService) SetInboundAssignments(ctx context.Context, req InboundAssignmentRequest) ([]ConfigProfileInboundToNode, error) {
	err := s.repo.setInboundAssignments(ctx, req.NodeUUID, req.InboundUUIDs)
	if err != nil {
		return nil, err
	}

	monitor.RequestNodeDeploy(true, req.NodeUUID)

	return s.repo.getInboundAssignments(ctx, req.NodeUUID)
}

func (s *SquadService) SetSquadInbounds(ctx context.Context, req SquadInboundsRequest) error {
	err := s.repo.setSquadInbounds(ctx, req.SquadUUID, req.InboundUUIDs)
	if err != nil {
		return err
	}
	monitor.RequestNodeDeploy(true)
	return nil
}

func (s *SquadService) SetSquadMembers(ctx context.Context, req SquadMembersRequest) error {
	err := s.repo.setSquadMembers(ctx, req.SquadUUID, req.UserIDs)
	if err != nil {
		return err
	}
	monitor.RequestNodeDeploy(true)
	return nil
}
