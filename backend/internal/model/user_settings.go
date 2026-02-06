package model

import (
	"time"

	"gorm.io/gorm"
)

// UserSettings represents user-specific settings for plugin control
type UserSettings struct {
	UserID               string `gorm:"primary_key;column:user_id;type:varchar(255)" json:"userId"`
	CollectionEnabled    bool   `gorm:"column:collection_enabled;not null;default:false" json:"collectionEnabled"`
	CollectionDailyLimit int    `gorm:"column:collection_daily_limit;not null;default:500" json:"collectionDailyLimit"`
	CollectionBatchLimit int    `gorm:"column:collection_batch_limit;not null;default:50" json:"collectionBatchLimit"`
	CreatedAt            time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
	UpdatedAt            time.Time `gorm:"column:updated_at;not null;default:now()" json:"updatedAt"`
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
