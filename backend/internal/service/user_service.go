package service

import (
	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/internal/repository"
)

// UserService 用户服务
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	AuthCenterUserID string                 `json:"authCenterUserId" binding:"required"`
	Role             string                 `json:"role"`
	Profile          map[string]interface{} `json:"profile"`
}

// Create 创建用户
func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
	user := &model.User{
		AuthCenterUserID: req.AuthCenterUserID, // 直接使用字符串
		Role:             req.Role,
		Profile:          req.Profile,
	}

	if user.Role == "" {
		user.Role = "USER"
	}

	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID 根据 ID 获取用户
func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// GetByAuthCenterUserID 根据账号中心用户 ID 获取用户
func (s *UserService) GetByAuthCenterUserID(authCenterUserID string) (*model.User, error) {
	return s.userRepo.GetByAuthCenterUserID(authCenterUserID)
}

// SyncUserFromAuthCenter 从账号中心同步用户信息
func (s *UserService) SyncUserFromAuthCenter(authCenterUserID string, nickname interface{}, headimgurl interface{}) (*model.User, error) {
	// 尝试获取现有用户
	existingUser, err := s.GetByAuthCenterUserID(authCenterUserID)

	// 将 interface{} 转换为 *string
	var nicknameStr *string
	if nickname != nil {
		if str, ok := nickname.(string); ok && str != "" {
			nicknameStr = &str
		}
	}

	var avatarURLStr *string
	if headimgurl != nil {
		if str, ok := headimgurl.(string); ok && str != "" {
			avatarURLStr = &str
		}
	}

	var user *model.User
	if err == nil && existingUser != nil {
		// 用户已存在，更新独立字段
		user = existingUser
		user.Nickname = nicknameStr
		user.AvatarURL = avatarURLStr
	} else {
		// 新用户（和 PR 系统一样，直接使用字符串）
		user = &model.User{
			AuthCenterUserID: authCenterUserID, // 直接使用字符串
			Role:             "USER",
			Nickname:         nicknameStr,
			AvatarURL:        avatarURLStr,
		}
	}

	err = s.userRepo.UpsertByAuthCenterUserID(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update 更新用户
func (s *UserService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

// Delete 删除用户
func (s *UserService) Delete(id string) error {
	return s.userRepo.Delete(id)
}
