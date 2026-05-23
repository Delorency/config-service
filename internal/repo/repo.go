package repo

import (
	"context"
	"errors"
	"fmt"
	"main/internal/models"
	"time"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

func (r *ConfigRepository) GetActualVersion(ctx context.Context, serviceID string) (*models.Config, error) {
	var config models.Config
	err := r.db.WithContext(ctx).
		Where("service_id = ?", serviceID).
		Order("created_at DESC").
		First(&config).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &config, err
}

func (r *ConfigRepository) GetVersion(ctx context.Context, serviceID string, version int) (*models.ConfigHistory, error) {
	var history models.ConfigHistory
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND version = ?", serviceID, version).
		First(&history).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &history, err
}

func (r *ConfigRepository) ListVersions(ctx context.Context, serviceID string, limit, offset int) ([]models.ConfigHistory, int64, error) {
	var history []models.ConfigHistory
	var total int64

	query := r.db.WithContext(ctx).Where("service_id = ?", serviceID)

	if err := query.Model(&models.ConfigHistory{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("version DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error

	return history, total, err
}

func (r *ConfigRepository) Update(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Config{}).
			Where("service_id = ?", config.ServiceID).
			Updates(map[string]any{
				"data":       config.Data,
				"version":    config.Version,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}

		history := models.ConfigHistory{
			ServiceID: config.ServiceID,
			Data:      config.Data,
			Version:   config.Version,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to save history: %w", err)
		}

		return nil
	})
}

func (r *ConfigRepository) Create(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(config).Error; err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}

		history := models.ConfigHistory{
			ServiceID: config.ServiceID,
			Data:      config.Data,
			Version:   config.Version,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to save history: %w", err)
		}

		return nil
	})
}

func (r *ConfigRepository) CreateOrUpdate(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.Config
		err := tx.Where("service_id = ?", config.ServiceID).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(config).Error; err != nil {
				return fmt.Errorf("failed to create config: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check existing config: %w", err)
		} else {
			if err := tx.Model(&models.Config{}).
				Where("service_id = ?", config.ServiceID).
				Updates(map[string]interface{}{
					"data":       config.Data,
					"version":    config.Version,
					"updated_at": time.Now(),
				}).Error; err != nil {
				return fmt.Errorf("failed to update config: %w", err)
			}
		}

		history := models.ConfigHistory{
			ServiceID: config.ServiceID,
			Data:      config.Data,
			Version:   config.Version,
		}

		if err := tx.Create(&history).Error; err != nil {
			return fmt.Errorf("failed to save history: %w", err)
		}

		return nil
	})
}

func (r *ConfigRepository) Rollback(ctx context.Context, serviceID string, targetVersion int) (*models.Config, error) {
	var result *models.Config

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var target models.ConfigHistory
		if err := tx.Where("service_id = ? AND version = ?", serviceID, targetVersion).
			First(&target).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("version %d not found", targetVersion)
			}
			return fmt.Errorf("failed to get target version: %w", err)
		}

		var current models.Config
		if err := tx.Where("service_id = ?", serviceID).First(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("current config not found")
			}
			return fmt.Errorf("failed to get current config: %w", err)
		}

		if current.Version == targetVersion {
			return fmt.Errorf("already at version %d", targetVersion)
		}

		nextVersion := current.Version + 1

		if err := tx.Model(&models.Config{}).
			Where("service_id = ?", serviceID).
			Updates(map[string]any{
				"data":       target.Data,
				"version":    nextVersion,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return fmt.Errorf("failed to update config: %w", err)
		}

		rollbackHistory := models.ConfigHistory{
			ServiceID: serviceID,
			Data:      target.Data,
			Version:   nextVersion,
		}

		if err := tx.Create(&rollbackHistory).Error; err != nil {
			return fmt.Errorf("failed to save rollback history: %w", err)
		}

		current.Data = target.Data
		current.Version = nextVersion
		result = &current

		return nil
	})

	return result, err
}
