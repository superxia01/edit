package service

import (
	"errors"
	"time"

	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/internal/repository"
)

// AdminService 管理后台服务
type AdminService struct {
	userRepo        *repository.UserRepository
	apiKeyRepo      *repository.APIKeyRepository
	userSettingsRepo *repository.UserSettingsRepository
	noteRepo        *repository.NoteRepository
	bloggerRepo     *repository.BloggerRepository
	statsService    *StatsService
	apiKeyService   *APIKeyService
}

// NewAdminService 创建 Admin 服务实例
func NewAdminService(
	userRepo *repository.UserRepository,
	apiKeyRepo *repository.APIKeyRepository,
	userSettingsRepo *repository.UserSettingsRepository,
	noteRepo *repository.NoteRepository,
	bloggerRepo *repository.BloggerRepository,
	statsService *StatsService,
	apiKeyService *APIKeyService,
) *AdminService {
	return &AdminService{
		userRepo:         userRepo,
		apiKeyRepo:       apiKeyRepo,
		userSettingsRepo: userSettingsRepo,
		noteRepo:         noteRepo,
		bloggerRepo:      bloggerRepo,
		statsService:     statsService,
		apiKeyService:    apiKeyService,
	}
}

// AdminUserListItem 管理后台用户列表项
type AdminUserListItem struct {
	ID               string  `json:"id"`
	AuthCenterUserID string  `json:"authCenterUserId"`
	Nickname         *string `json:"nickname,omitempty"`
	AvatarURL        *string `json:"avatarUrl,omitempty"`
	Role             string  `json:"role"`
	CreatedAt        string  `json:"createdAt"`
	TotalNotes       int64   `json:"totalNotes"`
	TotalBloggers    int64   `json:"totalBloggers"`
	HasAPIKey        bool    `json:"hasApiKey"`
}

// ListUsersResponse 分页用户列表响应
type ListUsersResponse struct {
	Items      []AdminUserListItem `json:"items"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	Size       int                 `json:"size"`
	TotalPages int                 `json:"totalPages"`
}

// ListUsers 获取所有用户（含统计，分页）
func (s *AdminService) ListUsers(page, size int) (*ListUsersResponse, error) {
	users, total, err := s.userRepo.ListAll(page, size)
	if err != nil {
		return nil, err
	}
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}

	items := make([]AdminUserListItem, 0, len(users))
	for _, u := range users {
		stats, _ := s.statsService.GetStats(u.ID)
		totalNotes := int64(0)
		totalBloggers := int64(0)
		if stats != nil {
			totalNotes = stats.TotalNotes
			totalBloggers = stats.TotalBloggers
		}

		apiKeys, _ := s.apiKeyRepo.GetByUserID(u.ID)
		hasAPIKey := len(apiKeys) > 0 && apiKeys[0].IsActive

		items = append(items, AdminUserListItem{
			ID:               u.ID,
			AuthCenterUserID: u.AuthCenterUserID,
			Nickname:         u.Nickname,
			AvatarURL:        u.AvatarURL,
			Role:             u.Role,
			CreatedAt:        u.CreatedAt.Format("2006-01-02 15:04:05"),
			TotalNotes:       totalNotes,
			TotalBloggers:    totalBloggers,
			HasAPIKey:        hasAPIKey,
		})
	}
	return &ListUsersResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		Size:       size,
		TotalPages: totalPages,
	}, nil
}

// AdminUserDetail 管理后台用户详情
type AdminUserDetail struct {
	User     *model.User        `json:"user"`
	Stats    *StatsResponse     `json:"stats"`
	Settings *model.UserSettings `json:"settings"`
	APIKeys  []APIKeyMasked     `json:"apiKeys"`
}

// APIKeyMasked 脱敏后的 API Key
type APIKeyMasked struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	IsActive  bool    `json:"isActive"`
	LastUsed  *string `json:"lastUsed,omitempty"`
	ExpiresAt *string `json:"expiresAt,omitempty"` // 到期日，nil 表示永不过期
	CreatedAt string  `json:"createdAt"`
}

// GetUserDetail 获取用户详情
func (s *AdminService) GetUserDetail(userID string) (*AdminUserDetail, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	stats, _ := s.statsService.GetStats(user.ID)
	settings, _ := s.userSettingsRepo.GetByUserID(user.ID)
	if settings == nil {
		settings = &model.UserSettings{
			UserID:               user.ID,
			CollectionEnabled:    false,
			CollectionDailyLimit: 500,
			CollectionBatchLimit: 50,
		}
	}

	apiKeys, _ := s.apiKeyRepo.GetByUserID(user.ID)
	keys := make([]APIKeyMasked, 0, len(apiKeys))
	for _, k := range apiKeys {
		lastUsed := ""
		if k.LastUsed != nil {
			lastUsed = k.LastUsed.Format("2006-01-02 15:04:05")
		}
		var expiresAt *string
		if k.ExpiresAt != nil {
			s := k.ExpiresAt.Format("2006-01-02 15:04:05")
			expiresAt = &s
		}
		keys = append(keys, APIKeyMasked{
			ID:        k.ID,
			Name:      k.Name,
			IsActive:  k.IsActive,
			LastUsed:  ptrString(lastUsed),
			ExpiresAt: expiresAt,
			CreatedAt: k.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &AdminUserDetail{
		User:     user,
		Stats:    stats,
		Settings: settings,
		APIKeys:  keys,
	}, nil
}

func ptrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// CreateAPIKeyForUser 管理员为用户创建 API Key（expiresIn 为天数，nil 表示永不过期）
func (s *AdminService) CreateAPIKeyForUser(userID string, expiresIn *int) (*APIKeyResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return s.apiKeyService.CreateForUserID(user.ID, expiresIn)
}

// UpdateAPIKeyExpiry 管理员修改 API Key 有效期（expiresIn 为天数，nil 表示永不过期）
func (s *AdminService) UpdateAPIKeyExpiry(userID string, apiKeyID string, expiresIn *int) error {
	key, err := s.apiKeyRepo.GetByID(apiKeyID)
	if err != nil {
		return err
	}
	if key.UserID != userID {
		return errors.New("API Key 不属于该用户")
	}
	var expiresAt *time.Time
	if expiresIn != nil && *expiresIn > 0 {
		t := time.Now().AddDate(0, 0, *expiresIn)
		expiresAt = &t
	}
	return s.apiKeyRepo.UpdateExpiresAt(apiKeyID, expiresAt)
}

// UpdateUserSettings 更新用户采集设置（dailyLimit、batchLimit、collectionEnabled）
func (s *AdminService) UpdateUserSettings(userID string, dailyLimit, batchLimit *int, collectionEnabled *bool) error {
	settings, err := s.userSettingsRepo.GetOrCreate(userID)
	if err != nil {
		return err
	}
	if dailyLimit != nil && *dailyLimit >= 0 {
		settings.CollectionDailyLimit = *dailyLimit
	}
	if batchLimit != nil && *batchLimit >= 0 {
		settings.CollectionBatchLimit = *batchLimit
	}
	if collectionEnabled != nil {
		settings.CollectionEnabled = *collectionEnabled
	}
	return s.userSettingsRepo.Update(settings)
}

// StatsOverview 全局统计（总用户数、总采集量等）
type StatsOverview struct {
	TotalUsers    int64 `json:"totalUsers"`
	TotalNotes    int64 `json:"totalNotes"`
	TotalBloggers int64 `json:"totalBloggers"`
}

// GetStatsOverview 获取全局统计
func (s *AdminService) GetStatsOverview() (*StatsOverview, error) {
	totalUsers, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}
	totalNotes, err := s.noteRepo.TotalCount()
	if err != nil {
		return nil, err
	}
	totalBloggers, err := s.bloggerRepo.TotalCount()
	if err != nil {
		return nil, err
	}
	return &StatsOverview{
		TotalUsers:    totalUsers,
		TotalNotes:    totalNotes,
		TotalBloggers: totalBloggers,
	}, nil
}
