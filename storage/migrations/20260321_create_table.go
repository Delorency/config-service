package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

func CreateTable() *gormigrate.Migration {
	return &gormigrate.Migration{
		// ID: "20260321_create_table",
		// Migrate: func(tx *gorm.DB) error {
		// 	return tx.AutoMigrate(&models.RefreshToken{}, &models.User{})
		// },
		// Rollback: func(tx *gorm.DB) error {
		// 	return tx.Migrator().DropTable(models.RefreshToken{}.TableName(), models.User{}.TableName())
		// },
	}
}
