package migrations

import (
	"main/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20260321_create_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&models.Config{}, &models.ConfigHistory{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(models.Config{}.TableName(), models.ConfigHistory{}.TableName())
		},
	}
}
