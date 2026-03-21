package migrations

import (
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		CreateTable(),
	})

	if err := m.Migrate(); err != nil {
		log.Fatalln("Migration error")
	}
}
