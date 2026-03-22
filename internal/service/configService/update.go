package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) UpdateConfig(ctx context.Context, serviceID string, data []byte) (*models.Config, error) {
	existing, err := s.repo.GetCurrent(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current config: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("config not found for service %s, use create instead", serviceID)
	}

	if s.validator != nil {
		if err := s.validator.ValidateByServiceID(serviceID, data); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	config := &models.Config{
		ServiceID: serviceID,
		Data:      data,
	}

	if err := s.repo.Update(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to update config: %w", err)
	}

	updated, _ := s.repo.GetCurrent(ctx, serviceID)

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      updated.Data,
		Version:   updated.Version,
		Timestamp: time.Now().Unix(),
	})

	return updated, nil
}
