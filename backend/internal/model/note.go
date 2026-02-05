package model

import (
    "time"

    "github.com/google/uuid"
    "github.com/lib/pq"
    "gorm.io/gorm"
)

// Note 笔记数据模型
// 遵循 KeenChase V3.0 规范：
// - 结构体 PascalCase
// - JSON camelCase
// - GORM column snake_case
// - 主键 UUID
type Note struct {
    ID              UUID              `gorm:"primaryKey;column:id;type:uuid;default:gen_random_uuid()" json:"id"`
    UserID          UUID              `gorm:"column:user_id;type:uuid;not null;index:idx_notes_user_id" json:"userId"`
    URL             string            `gorm:"column:url;type:varchar(500);not null;index:idx_notes_url" json:"url"`
    Title           string            `gorm:"column:title;type:varchar(500)" json:"title"`
    Author          string            `gorm:"column:author;type:varchar(100)" json:"author"`
    Content         string            `gorm:"column:content;type:text" json:"content"`
    Tags            pq.StringArray    `gorm:"column:tags;type:text[]" json:"tags"`
    ImageURLs       pq.StringArray    `gorm:"column:image_urls;type:text[]" json:"imageUrls"`
    VideoURL        *string           `gorm:"column:video_url;type:varchar(500)" json:"videoUrl,omitempty"`
    NoteType        string            `gorm:"column:note_type;type:varchar(20)" json:"noteType"`
    CoverImageURL   string            `gorm:"column:cover_image_url;type:varchar(500)" json:"coverImageUrl"`
    Likes           int32             `gorm:"column:likes;type:integer;default:0" json:"likes"`
    Collects        int32             `gorm:"column:collects;type:integer;default:0" json:"collects"`
    Comments        int32             `gorm:"column:comments;type:integer;default:0" json:"comments"`
    PublishDate     int64             `gorm:"column:publish_date;type:bigint" json:"publishDate"`
    Source          string            `gorm:"column:source;type:varchar(20);default:'single'" json:"source"` // single: 单篇, batch: 批量
    CaptureTimestamp int64            `gorm:"column:capture_timestamp;type:bigint;not null" json:"captureTimestamp"`
    CreatedAt       time.Time         `gorm:"column:created_at;type:timestamp with time zone;default:now();not null" json:"createdAt"`
    UpdatedAt       time.Time         `gorm:"column:updated_at;type:timestamp with time zone;default:now();not null" json:"updatedAt"`
}

// TableName 指定表名（复数 + snake_case）
func (Note) TableName() string {
    return "notes"
}

// BeforeCreate GORM hook
func (n *Note) BeforeCreate(tx *gorm.DB) error {
    if n.ID == uuid.Nil {
        n.ID = uuid.New()
    }
    return nil
}
