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

func (r *ConfigRepository) GetCurrent(ctx context.Context, serviceID string) (*models.Config, error) {
	var config models.Config
	err := r.db.WithContext(ctx).
		Where("service_id = ?", serviceID).
		First(&config).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &config, err
}

func (r *ConfigRepository) Create(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Create(config).Error
}

func (r *ConfigRepository) Update(ctx context.Context, config *models.Config) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var current models.Config
		if err := tx.Where("service_id = ?", config.ServiceID).First(&current).Error; err != nil {
			return err
		}

		history := models.ConfigHistory{
			ServiceID: current.ServiceID,
			Data:      current.Data,
			Version:   current.Version,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		current.Data = config.Data
		current.Version = current.Version + 1

		return tx.Save(&current).Error
	})
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

func (r *ConfigRepository) ListVersions(ctx context.Context, serviceID string, limit int) ([]models.ConfigHistory, error) {
	var history []models.ConfigHistory
	err := r.db.WithContext(ctx).
		Where("service_id = ?", serviceID).
		Order("version DESC").
		Limit(limit).
		Find(&history).Error

	return history, err
}

func (r *ConfigRepository) Rollback(ctx context.Context, serviceID string, targetVersion int) (*models.Config, error) {
	var result *models.Config

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var target models.ConfigHistory
		if err := tx.Where("service_id = ? AND version = ?", serviceID, targetVersion).
			First(&target).Error; err != nil {
			return fmt.Errorf("version %d not found", targetVersion)
		}

		var current models.Config
		if err := tx.Where("service_id = ?", serviceID).First(&current).Error; err != nil {
			return err
		}

		history := models.ConfigHistory{
			ServiceID: current.ServiceID,
			Data:      current.Data,
			Version:   current.Version,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		current.Data = target.Data
		current.Version = current.Version + 1

		if err := tx.Save(&current).Error; err != nil {
			return err
		}

		result = &current
		return nil
	})

	return result, err
}
