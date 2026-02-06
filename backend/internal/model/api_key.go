package model

import (
	"time"

	"gorm.io/gorm"
)

// APIKey represents an API key for plugin authentication
type APIKey struct {
	ID        string    `gorm:"primaryKey;column:id;type:varchar(255)" json:"id"`
	UserID    string    `gorm:"column:user_id;type:varchar(255);not null;index" json:"userId"`
	Name      string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Key       string    `gorm:"column:key;type:varchar(255);not null;uniqueIndex" json:"key"` // hashed value
	IsActive  bool      `gorm:"column:is_active;type:boolean;not null;default:true" json:"isActive"`
	LastUsed  *time.Time `gorm:"column:last_used" json:"lastUsed"`
	ExpiresAt *time.Time `gorm:"column:expires_at" json:"expiresAt"`
	CreatedAt time.Time `gorm:"column:created_at;not null;default:now()" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for APIKey
func (APIKey) TableName() string {
	return "api_keys"
}
