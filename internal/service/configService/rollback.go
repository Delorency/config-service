package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) Rollback(ctx context.Context, serviceID string, targetVersion int) (*models.Config, error) {
	target, err := s.repo.GetVersion(ctx, serviceID, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get target version: %w", err)
	}
	if target == nil {
		return nil, fmt.Errorf("version %d not found", targetVersion)
	}

	config, err := s.repo.Rollback(ctx, serviceID, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to rollback: %w", err)
	}

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      config.Data,
		Version:   config.Version,
		Timestamp: time.Now().Unix(),
	})

	return config, nil
}
