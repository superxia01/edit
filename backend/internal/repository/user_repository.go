package repository

import (
	"github.com/keenchase/edit-business/internal/model"

	"gorm.io/gorm"
)

// UserRepository 用户仓库
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID 根据 ID 获取用户
func (r *UserRepository) GetByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAuthCenterUserID 根据账号中心用户 ID 获取用户（返回具体类型，用于内部调用）
func (r *UserRepository) GetByAuthCenterUserID(authCenterUserID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("auth_center_user_id = ?", authCenterUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByAuthCenterUserIDInterface 根据账号中心用户 ID 获取用户（返回interface，用于AuthCenterMiddleware）
func (r *UserRepository) GetByAuthCenterUserIDInterface(authCenterUserID string) (interface{}, error) {
	return r.GetByAuthCenterUserID(authCenterUserID)
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

// UpsertByAuthCenterUserID 根据账号中心用户 ID 插入或更新用户
func (r *UserRepository) UpsertByAuthCenterUserID(user *model.User) error {
	var existing model.User
	err := r.db.Where("auth_center_user_id = ?", user.AuthCenterUserID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		return r.db.Create(user).Error
	} else if err != nil {
		return err
	}

	// 存在则更新
	user.ID = existing.ID
	return r.db.Save(user).Error
}

// CreateWithAuthCenter 根据账号中心用户信息创建本地用户（用于 AuthCenterMiddleware）
func (r *UserRepository) CreateWithAuthCenter(authCenterUserID string, unionID string, nickname string, avatarURL string) (interface{}, error) {
	// 先检查是否已存在
	existing, err := r.GetByAuthCenterUserID(authCenterUserID)
	if err == nil && existing != nil {
		return existing, nil // 已存在则直接返回
	}

	// 不存在则创建新用户（和 PR 系统一样，直接使用字符串）
	var nicknamePtr *string
	if nickname != "" {
		nicknamePtr = &nickname
	}

	var avatarURLPtr *string
	if avatarURL != "" {
		avatarURLPtr = &avatarURL
	}

	var unionIDPtr *string
	if unionID != "" {
		unionIDPtr = &unionID
	}

	user := &model.User{
		AuthCenterUserID: authCenterUserID, // 直接使用字符串，不需要 uuid.Parse
		Role:             "USER",
		UnionID:          unionIDPtr,
		Nickname:         nicknamePtr,
		AvatarURL:        avatarURLPtr,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
