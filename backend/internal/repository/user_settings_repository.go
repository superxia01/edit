package repository

import (
	"time"

	"github.com/keenchase/edit-business/internal/model"
	"gorm.io/gorm"
)

// UserSettingsRepository handles user settings data operations
type UserSettingsRepository struct {
	db *gorm.DB
}

// NewUserSettingsRepository creates a new user settings repository
func NewUserSettingsRepository(db *gorm.DB) *UserSettingsRepository {
	return &UserSettingsRepository{db: db}
}

// GetByUserID retrieves user settings by user ID
func (r *UserSettingsRepository) GetByUserID(userID string) (*model.UserSettings, error) {
	var settings model.UserSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default settings if not found
			return &model.UserSettings{
				UserID:               userID, // 直接使用字符串
				CollectionEnabled:    false,
				CollectionDailyLimit: 500,
				CollectionBatchLimit: 50,
			}, nil
		}
		return nil, err
	}
	return &settings, nil
}

// GetOrCreate 获取或创建用户设置（用于管理后台更新限额时确保记录存在）
func (r *UserSettingsRepository) GetOrCreate(userID string) (*model.UserSettings, error) {
	settings, err := r.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	// 若为内存默认值（未持久化），先创建
	var count int64
	if r.db.Model(&model.UserSettings{}).Where("user_id = ?", userID).Count(&count).Error == nil && count == 0 {
		if err := r.Create(settings); err != nil {
			return nil, err
		}
	}
	return settings, nil
}

// Create creates user settings
func (r *UserSettingsRepository) Create(settings *model.UserSettings) error {
	return r.db.Create(settings).Error
}

// Update updates user settings
func (r *UserSettingsRepository) Update(settings *model.UserSettings) error {
	return r.db.Save(settings).Error
}

// UpdateCollectionEnabled toggles collection enabled status
func (r *UserSettingsRepository) UpdateCollectionEnabled(userID string, enabled bool) error {
	return r.db.Model(&model.UserSettings{}).
		Where("user_id = ?", userID).
		Update("collection_enabled", enabled).Error
}

// IncrementDailyCount increments the daily collection count (stored separately)
// This is called when a note is collected
func (r *UserSettingsRepository) GetDailyCount(userID string) (int64, error) {
	// Check if there's a daily count record for today
	var count int64
	today := time.Now().Format("2006-01-02")

	// For simplicity, we'll use a simple approach:
	// Check notes table for today's count
	err := r.db.Table("notes").
		Where("user_id = ? AND DATE(created_at) = ?", userID, today).
		Count(&count).Error

	return count, err
}
