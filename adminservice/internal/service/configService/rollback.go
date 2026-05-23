package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) Rollback(ctx context.Context, serviceID string, targetVersion int) (*models.Config, error) {
	if s.validator != nil {
		if !s.validator.SchemaExists(serviceID) {
			return nil, fmt.Errorf("schema not found for service '%s'. Cannot rollback without schema", serviceID)
		}
	}

	target, err := s.repo.GetVersion(ctx, serviceID, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to get target version %d: %w", targetVersion, err)
	}

	if s.validator != nil {
		if err := s.validator.ValidateByServiceID(serviceID, target.Data); err != nil {
			return nil, fmt.Errorf("target version %d validation failed: %w", targetVersion, err)
		}
	}

	config, err := s.repo.Rollback(ctx, serviceID, targetVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to rollback to version %d: %w", targetVersion, err)
	}

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      config.Data,
		Version:   config.Version,
		Timestamp: time.Now().Unix(),
	})

	return config, nil
}
