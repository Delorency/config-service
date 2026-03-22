package configService

import (
	"context"
	"fmt"
	"main/internal/models"
)

func (s *ConfigService) GetConfig(ctx context.Context, serviceID string) (*models.Config, error) {
	config, err := s.repo.GetCurrent(ctx, serviceID)
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
