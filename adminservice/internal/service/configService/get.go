package configService

import (
	"context"
	"fmt"
	"main/internal/models"
)

func (s *ConfigService) GetConfig(ctx context.Context, serviceID string) (*models.Config, error) {
	config, err := s.repo.GetActualVersion(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("config not found for service %s", serviceID)
	}
	return config, nil
}

func (s *ConfigService) GetVersion(ctx context.Context, serviceID string, version int) (*models.ConfigHistory, error) {
	history, err := s.repo.GetVersion(ctx, serviceID, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	if history == nil {
		return nil, fmt.Errorf("version %d not found for service %s", version, serviceID)
	}
	return history, nil
}

func (s *ConfigService) ListVersions(ctx context.Context, serviceID string, limit, offset int) ([]models.ConfigHistory, int64, error) {
	current, err := s.repo.GetActualVersion(ctx, serviceID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check config: %w", err)
	}
	if current == nil {
		return nil, 0, fmt.Errorf("config for service %s not found", serviceID)
	}

	// Получаем историю версий
	versions, total, err := s.repo.ListVersions(ctx, serviceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list versions: %w", err)
	}

	return versions, total, nil
}
