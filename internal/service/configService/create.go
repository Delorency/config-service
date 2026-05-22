package configService

import (
	"context"
	"fmt"
	"main/internal/models"
	"main/internal/watcher"
	"time"
)

func (s *ConfigService) CreateConfig(ctx context.Context, serviceID string, data []byte) (*models.Config, error) {
	current, _ := s.repo.GetActualVersion(ctx, serviceID)

	nextVersion := 1
	if current != nil {
		nextVersion = current.Version + 1
	}

	if s.validator != nil {
		if err := s.validator.ValidateByServiceID(serviceID, data); err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}
	}

	config := &models.Config{
		ServiceID: serviceID,
		Data:      data,
		Version:   nextVersion,
	}

	if err := s.repo.Create(ctx, config); err != nil {
		return nil, err
	}

	s.watcher.Notify(serviceID, &watcher.ConfigUpdate{
		ServiceID: serviceID,
		Data:      data,
		Version:   nextVersion,
		Timestamp: time.Now().Unix(),
	})

	return config, nil
}
