package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) CreateConfig(ctx context.Context, serviceID string, data []byte) (*models.Config, error) {
	if s.validator != nil {
		if err := s.validator.ValidateByServiceID(serviceID, data); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	existing, _ := s.repo.GetActualVersion(ctx, serviceID)

	nextVersion := 1
	if existing != nil {
		nextVersion = existing.Version + 1
	}

	config := &models.Config{
		ServiceID: serviceID,
		Data:      data,
		Version:   nextVersion,
	}

	if err := s.repo.CreateOrUpdate(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      data,
		Version:   nextVersion,
		Timestamp: time.Now().Unix(),
	})

	return config, nil
}
