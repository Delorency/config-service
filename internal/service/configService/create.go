package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) CreateConfig(ctx context.Context, serviceID string, data []byte) (*models.Config, error) {
	_, err := s.repo.GetCurrent(ctx, serviceID)
	if err == nil {
		return nil, fmt.Errorf("config for service %s already exists", serviceID)
	}

	if s.validator != nil {
		if err := s.validator.ValidateByServiceID(serviceID, data); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	config := &models.Config{
		ServiceID: serviceID,
		Data:      data,
		Version:   1,
	}

	if err := s.repo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	updated, _ := s.repo.GetCurrent(ctx, serviceID)

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      updated.Data,
		Version:   updated.Version,
		Timestamp: time.Now().Unix(),
	})

	return config, nil
}
