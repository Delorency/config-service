package repo

import (
	"context"
	"errors"
	"fmt"
	"main/internal/models"

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

func (r *ConfigRepository) Create(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(config).Error; err != nil {
			return err
		}

		history := models.ConfigHistory{
			ServiceID: config.ServiceID,
			Data:      config.Data,
			Version:   config.Version,
		}

		return tx.Create(&history).Error
	})
}

func (r *ConfigRepository) Rollback(ctx context.Context, serviceID string, targetVersion int) (*models.Config, error) {
	target, err := r.GetVersion(ctx, serviceID, targetVersion)
	if err != nil {
		return nil, err
	}
	if target == nil {
		return nil, fmt.Errorf("version %d not found", targetVersion)
	}

	current, err := r.GetActualVersion(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("current config not found")
	}

	newConfig := &models.Config{
		ServiceID: serviceID,
		Data:      target.Data,
		Version:   current.Version + 1,
	}

	err = r.Create(ctx, newConfig)
	if err != nil {
		return nil, fmt.Errorf("rollback failed: %w", err)
	}

	return newConfig, nil
}
