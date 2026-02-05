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

// GetByAuthCenterUserID 根据账号中心用户 ID 获取用户
func (r *UserRepository) GetByAuthCenterUserID(authCenterUserID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("auth_center_user_id = ?", authCenterUserID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
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
