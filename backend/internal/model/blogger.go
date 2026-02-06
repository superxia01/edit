package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Blogger 博主信息模型
type Blogger struct {
	ID               string      `gorm:"primaryKey;column:id;type:varchar(255)" json:"id"`
	UserID           string      `gorm:"column:user_id;type:varchar(255);not null;index" json:"userId"`
	XhsID            string      `gorm:"uniqueIndex;column:xhs_id;type:varchar(50)" json:"xhsId"`
	BloggerName      string      `gorm:"column:blogger_name;type:varchar(100)" json:"bloggerName"`
	AvatarURL        string      `gorm:"column:avatar_url;type:varchar(500)" json:"avatarUrl"`
	Description      string      `gorm:"column:description;type:text" json:"description"`
	FollowersCount   int32       `gorm:"column:followers_count;type:integer;default:0" json:"followersCount"`
	BloggerURL       string      `gorm:"column:blogger_url;type:varchar(500)" json:"bloggerUrl"`
	CaptureTimestamp int64       `gorm:"column:capture_timestamp;type:bigint;not null" json:"captureTimestamp"`
	CreatedAt        time.Time   `gorm:"column:created_at;type:timestamp with time zone;default:now();not null" json:"createdAt"`
	UpdatedAt        time.Time   `gorm:"column:updated_at;type:timestamp with time zone;default:now();not null" json:"updatedAt"`
}

// TableName 指定表名（复数 + snake_case）
func (Blogger) TableName() string {
	return "bloggers"
}

// BeforeCreate GORM hook
func (b *Blogger) BeforeCreate(tx *gorm.DB) error {
	if b.ID == "" {
		b.ID = fmt.Sprintf("blogger-%d", time.Now().UnixNano())
	}
	return nil
}
