package service

import (
	"errors"
	"time"

	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/internal/repository"
)

var (
	ErrMaxAPIKeysReached = errors.New("each user can only have one API key")
	ErrInvalidAPIKey     = errors.New("invalid API key")
	ErrAPIKeyNotFound    = errors.New("API key not found, please contact admin")
)

// GetOrCreateAPIKeyByUser gets or creates API key using the user object directly (avoids lookup race for new users)
func (s *APIKeyService) GetOrCreateAPIKeyByUser(user *model.User) (*APIKeyResponse, error) {
	if user == nil || user.ID == "" {
		return nil, errors.New("invalid user")
	}
	return s.getOrCreateByUserID(user.ID)
}

// GetOrCreateAPIKey gets existing API key or creates a new one for the user
func (s *APIKeyService) GetOrCreateAPIKey(authCenterUserID string) (*APIKeyResponse, error) {
	// Get user by auth center user ID first
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}
	return s.getOrCreateByUserID(user.ID)
}

func (s *APIKeyService) getOrCreateByUserID(userID string) (*APIKeyResponse, error) {
	// 仅获取，不再自动创建。无 API Key 时返回 ErrAPIKeyNotFound
	apiKeys, err := s.apiKeyRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	for _, apiKey := range apiKeys {
		if apiKey.IsActive {
			return &APIKeyResponse{
				ID:        apiKey.ID,
				Name:      apiKey.Name,
				Key:       apiKey.Key,
				IsActive:  apiKey.IsActive,
				LastUsed:  apiKey.LastUsed,
				ExpiresAt: apiKey.ExpiresAt,
				CreatedAt: apiKey.CreatedAt,
			}, nil
		}
	}
	return nil, ErrAPIKeyNotFound
}

// APIKeyService handles API key business logic
type APIKeyService struct {
	apiKeyRepo *repository.APIKeyRepository
	userRepo   *repository.UserRepository
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(apiKeyRepo *repository.APIKeyRepository, userRepo *repository.UserRepository) *APIKeyService {
	return &APIKeyService{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

// CreateAPIKeyRequest represents the request to create an API key
type CreateAPIKeyRequest struct {
	Name      string `json:"name" binding:"required,min=1,max=255"`
	ExpiresIn *int   `json:"expiresIn"` // Optional expiration in days
}

// APIKeyResponse represents the API key response
type APIKeyResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Key       string      `json:"key"` // Only shown on creation
	IsActive  bool        `json:"isActive"`
	LastUsed  *time.Time  `json:"lastUsed"`
	ExpiresAt *time.Time  `json:"expiresAt"`
	CreatedAt time.Time   `json:"createdAt"`
}

// CreateForUserID 管理员为指定用户创建 API Key（支持设置有效期，expiresIn 为天数，nil 表示永不过期）
func (s *APIKeyService) CreateForUserID(userID string, expiresIn *int) (*APIKeyResponse, error) {
	// 先停用已过期的 Key，便于在旧 Key 过期后管理员可重新生成
	_ = s.apiKeyRepo.DeactivateExpiredForUser(userID)

	count, err := s.apiKeyRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, ErrMaxAPIKeysReached
	}
	keyString := s.apiKeyRepo.GenerateKey()
	apiKey := &model.APIKey{
		UserID:   userID,
		Name:     "默认API Key",
		Key:      keyString,
		IsActive: true,
	}
	if expiresIn != nil && *expiresIn > 0 {
		expires := time.Now().AddDate(0, 0, *expiresIn)
		apiKey.ExpiresAt = &expires
	}
	if err := s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, err
	}
	return &APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       keyString,
		IsActive:  apiKey.IsActive,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

// Create creates a new API key for a user
func (s *APIKeyService) Create(authCenterUserID string, req CreateAPIKeyRequest) (*APIKeyResponse, error) {
	// Get user by auth center user ID
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}

	// Check if user already has an API key (max 1 per user)
	count, err := s.apiKeyRepo.CountByUserID(user.ID)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, ErrMaxAPIKeysReached
	}

	// Generate API key
	keyString := s.apiKeyRepo.GenerateKey()

	// Calculate expiration if provided
	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expires := time.Now().AddDate(0, 0, *req.ExpiresIn)
		expiresAt = &expires
	}

	// Create API key record
	apiKey := &model.APIKey{
		UserID:    user.ID,
		Name:      req.Name,
		Key:       keyString,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := s.apiKeyRepo.Create(apiKey); err != nil {
		return nil, err
	}

	return &APIKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		Key:       keyString, // Only return the key on creation
		IsActive:  apiKey.IsActive,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

// List returns all API keys for a user
func (s *APIKeyService) List(authCenterUserID string) ([]APIKeyResponse, error) {
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}

	apiKeys, err := s.apiKeyRepo.GetByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	responses := make([]APIKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		responses[i] = APIKeyResponse{
			ID:        apiKey.ID,
			Name:      apiKey.Name,
			Key:       maskAPIKey(apiKey.Key),
			IsActive:  apiKey.IsActive,
			LastUsed:  apiKey.LastUsed,
			ExpiresAt: apiKey.ExpiresAt,
			CreatedAt: apiKey.CreatedAt,
		}
	}

	return responses, nil
}

// GetStats returns API key statistics for a user
func (s *APIKeyService) GetStats(authCenterUserID string) (*repository.APIKeyStats, error) {
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}

	return s.apiKeyRepo.GetStats(user.ID)
}

// Delete deletes an API key
func (s *APIKeyService) Delete(authCenterUserID string, apiKeyID string) error {
	// Get user
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return err
	}

	// Get API key
	apiKey, err := s.apiKeyRepo.GetByID(apiKeyID)
	if err != nil {
		return err
	}

	// Verify ownership
	if apiKey.UserID != user.ID {
		return ErrInvalidAPIKey
	}

	return s.apiKeyRepo.Delete(apiKeyID)
}

// Deactivate deactivates an API key
func (s *APIKeyService) Deactivate(authCenterUserID string, apiKeyID string) error {
	// Get user
	user, err := s.userRepo.GetByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return err
	}

	// Get API key
	apiKey, err := s.apiKeyRepo.GetByID(apiKeyID)
	if err != nil {
		return err
	}

	// Verify ownership
	if apiKey.UserID != user.ID {
		return ErrInvalidAPIKey
	}

	return s.apiKeyRepo.Deactivate(apiKeyID)
}

// ValidateAPIKey validates an API key and returns the user ID
func (s *APIKeyService) ValidateAPIKey(key string) (string, error) {
	apiKey, err := s.apiKeyRepo.GetByKey(key)
	if err != nil {
		return "", ErrInvalidAPIKey
	}

	// Check if expired
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return "", ErrInvalidAPIKey
	}

	// Update last used timestamp
	go s.apiKeyRepo.UpdateLastUsed(apiKey.ID)

	return apiKey.UserID, nil
}

// ValidateAPIKeyWithAuthCenterID validates an API key and returns both user ID and auth center user ID
func (s *APIKeyService) ValidateAPIKeyWithAuthCenterID(key string) (string, string, error) {
	apiKey, err := s.apiKeyRepo.GetByKey(key)
	if err != nil {
		return "", "", ErrInvalidAPIKey
	}

	// Check if expired
	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return "", "", ErrInvalidAPIKey
	}

	// Get user to retrieve authCenterUserID
	user, err := s.userRepo.GetByID(apiKey.UserID)
	if err != nil {
		return "", "", ErrInvalidAPIKey
	}

	// Update last used timestamp
	go s.apiKeyRepo.UpdateLastUsed(apiKey.ID)

	return apiKey.UserID, user.AuthCenterUserID, nil
}

// maskAPIKey masks the API key for display (only show first 8 and last 4 characters)
func maskAPIKey(key string) string {
	if len(key) < 12 {
		return "****"
	}
	return key[:8] + "..." + key[len(key)-4:]
}
