package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Config struct {
	gorm.Model
	ServiceID string         `gorm:"uniqueIndex;size:100;not null;column:service_id"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null"`
	Version   int            `gorm:"not null;default:1"`
}

func (Config) TableName() string {
	return "config"
}

type ConfigHistory struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	ServiceID string         `gorm:"size:100;not null;index;column:service_id"`
	Data      datatypes.JSON `gorm:"type:jsonb;not null"`
	Version   int            `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

func (ConfigHistory) TableName() string {
	return "config_history"
}
