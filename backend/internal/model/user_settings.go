package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserSettings represents user-specific settings for plugin control
type UserSettings struct {
	UserID               uuid.UUID `gorm:"type:uuid;primary_key" json:"userId"`
	CollectionEnabled    bool      `gorm:"not null;default:false" json:"collectionEnabled"`
	CollectionDailyLimit int       `gorm:"not null;default:500" json:"collectionDailyLimit"`
	CollectionBatchLimit int       `gorm:"not null;default:50" json:"collectionBatchLimit"`
	CreatedAt            time.Time `gorm:"not null;default:now()" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"not null;default:now()" json:"updatedAt"`
}

// TableName specifies the table name for UserSettings
func (UserSettings) TableName() string {
	return "user_settings"
}

// BeforeUpdate hook to update UpdatedAt
func (us *UserSettings) BeforeUpdate(tx *gorm.DB) error {
	us.UpdatedAt = time.Now()
	return nil
}
