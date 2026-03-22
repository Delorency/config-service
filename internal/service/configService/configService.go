package configService

import (
	"context"
	"main/internal/core/validator"
	"main/internal/models"
	"main/internal/repo"
	"main/internal/watcher"
)

type ConfigService struct {
	repo      *repo.ConfigRepository
	watcher   *watcher.WatcherManager
	validator *validator.Validator
}

type ConfigServiceI interface {
	GetConfig(context.Context, string) (*models.Config, error)
	GetVersion(context.Context, string, int) (*models.ConfigHistory, error)
	CreateConfig(context.Context, string, []byte) (*models.Config, error)
	UpdateConfig(context.Context, string, []byte) (*models.Config, error)
	Rollback(context.Context, string, int) (*models.Config, error)
}

func NewConfigService(
	repo *repo.ConfigRepository,
	watcher *watcher.WatcherManager,
	validator *validator.Validator,
) ConfigServiceI {
	return &ConfigService{
		repo:      repo,
		watcher:   watcher,
		validator: validator,
	}
}
