package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey represents an API key for plugin authentication
type APIKey struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_api_keys_user_id" json:"userId"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Key       string    `gorm:"type:varchar(255);not null;uniqueIndex" json:"key"` // hashed value
	IsActive  bool      `gorm:"type:boolean;not null;default:true" json:"isActive"`
	LastUsed  *time.Time `json:"lastUsed"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for APIKey
func (APIKey) TableName() string {
	return "api_keys"
}
