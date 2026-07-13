package configprofiles

import (
	"context"
	"fmt"
	"strings"

	"exodus/internal/config"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

type ConfigProfileService struct {
	repo *ConfigProfileRepository
	cfg  *config.BackendConfig
}

func NewConfigProfileService(repo *ConfigProfileRepository, cfg *config.BackendConfig) *ConfigProfileService {
	return &ConfigProfileService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *ConfigProfileService) CreateConfigProfile(ctx context.Context, req createConfigProfileRequest) (configProfileRecord, error) {
	profileUUID := uuid.NewString()
	err := s.repo.createConfigProfile(ctx, profileUUID, req)
	if err != nil {
		return configProfileRecord{}, err
	}

	monitor.RequestNodeDeploy(true)
	return s.repo.getConfigProfileRecordByUUID(ctx, profileUUID)
}

func (s *ConfigProfileService) UpdateConfigProfile(ctx context.Context, req updateConfigProfileRequest) (configProfileRecord, error) {
	clauses := make([]string, 0)
	args := make([]any, 0)
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	if req.Name != nil {
		add("name", strings.TrimSpace(*req.Name))
	}
	if req.Config != nil {
		add("config", *req.Config)
	}

	err := s.repo.updateConfigProfile(ctx, req.UUID, clauses, args, req.Config)
	if err != nil {
		return configProfileRecord{}, err
	}

	monitor.RequestNodeDeploy(true)
	return s.repo.getConfigProfileRecordByUUID(ctx, req.UUID)
}

func (s *ConfigProfileService) DeleteConfigProfile(ctx context.Context, profileUUID string) error {
	err := s.repo.deleteConfigProfile(ctx, profileUUID)
	if err != nil {
		return err
	}
	monitor.RequestNodeDeploy(true)
	return nil
}

func (s *ConfigProfileService) ReorderConfigProfiles(ctx context.Context, req reorderConfigProfilesRequest) error {
	return s.repo.reorderConfigProfiles(ctx, req.Items)
}

func (s *ConfigProfileService) CreateSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	return s.repo.createSnippet(ctx, req)
}

func (s *ConfigProfileService) UpdateSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	return s.repo.updateSnippet(ctx, req)
}

func (s *ConfigProfileService) DeleteSnippet(ctx context.Context, name string) error {
	return s.repo.deleteSnippet(ctx, name)
}
