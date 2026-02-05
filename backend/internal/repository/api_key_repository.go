package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/keenchase/edit-business/internal/model"
	"gorm.io/gorm"
)

var (
	ErrAPIKeyNotFound = errors.New("API key not found")
)

// APIKeyRepository handles API key data operations
type APIKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new API key repository
func NewAPIKeyRepository(db *gorm.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// Create creates a new API key
func (r *APIKeyRepository) Create(apiKey *model.APIKey) error {
	return r.db.Create(apiKey).Error
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(id string) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := r.db.Where("id = ?", id).First(&apiKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAPIKeyNotFound
		}
		return nil, err
	}
	return &apiKey, nil
}

// GetByKey retrieves an API key by the key value (hashed)
func (r *APIKeyRepository) GetByKey(key string) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := r.db.Where("key = ? AND is_active = ?", key, true).First(&apiKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAPIKeyNotFound
		}
		return nil, err
	}
	return &apiKey, nil
}

// GetByUserID retrieves all API keys for a user
func (r *APIKeyRepository) GetByUserID(userID string) ([]model.APIKey, error) {
	var apiKeys []model.APIKey
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&apiKeys).Error
	return apiKeys, err
}

// UpdateLastUsed updates the last used timestamp
func (r *APIKeyRepository) UpdateLastUsed(id string) error {
	now := time.Now()
	return r.db.Model(&model.APIKey{}).Where("id = ?", id).Update("last_used", now).Error
}

// Deactivate deactivates an API key
func (r *APIKeyRepository) Deactivate(id string) error {
	return r.db.Model(&model.APIKey{}).Where("id = ?", id).Update("is_active", false).Error
}

// Delete permanently deletes an API key
func (r *APIKeyRepository) Delete(id string) error {
	return r.db.Delete(&model.APIKey{}, "id = ?", id).Error
}

// GenerateKey generates a random API key string
func (r *APIKeyRepository) GenerateKey() string {
	b := make([]byte, 32)
	rand.Read(b)
	return "eb_" + hex.EncodeToString(b)
}

// CountByUserID returns the count of active API keys for a user
func (r *APIKeyRepository) CountByUserID(userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.APIKey{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&count).Error
	return count, err
}

// TotalCount returns total count of API keys for analytics
func (r *APIKeyRepository) TotalCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.APIKey{}).Count(&count).Error
	return count, err
}

// GetStats returns API key statistics for a user
type APIKeyStats struct {
	TotalCount    int64     `json:"totalCount"`
	ActiveCount   int64     `json:"activeCount"`
	TotalUsage    int64     `json:"totalUsage"`
	LastUsed      *time.Time `json:"lastUsed"`
}

func (r *APIKeyRepository) GetStats(userID string) (*APIKeyStats, error) {
	var stats APIKeyStats

	// Count total keys
	r.db.Model(&model.APIKey{}).Where("user_id = ?", userID).Count(&stats.TotalCount)

	// Count active keys
	r.db.Model(&model.APIKey{}).Where("user_id = ? AND is_active = ?", userID, true).Count(&stats.ActiveCount)

	// Get last used time
	var apiKey model.APIKey
	err := r.db.Where("user_id = ? AND last_used IS NOT NULL", userID).
		Order("last_used DESC").
		First(&apiKey).Error
	if err == nil {
		stats.LastUsed = apiKey.LastUsed
	}

	return &stats, nil
}
