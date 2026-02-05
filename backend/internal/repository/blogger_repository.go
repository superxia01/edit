package repository

import (
	"github.com/keenchase/edit-business/internal/model"

	"gorm.io/gorm"
)

// BloggerRepository 博主信息仓库
type BloggerRepository struct {
	db *gorm.DB
}

// NewBloggerRepository 创建博主仓库实例
func NewBloggerRepository(db *gorm.DB) *BloggerRepository {
	return &BloggerRepository{db: db}
}

// Create 创建博主信息
func (r *BloggerRepository) Create(blogger *model.Blogger) error {
	return r.db.Create(blogger).Error
}

// GetByID 根据 ID 获取博主信息
func (r *BloggerRepository) GetByID(id string) (*model.Blogger, error) {
	var blogger model.Blogger
	err := r.db.Where("id = ?", id).First(&blogger).Error
	if err != nil {
		return nil, err
	}
	return &blogger, nil
}

// GetByXhsID 根据小红书 ID 获取博主信息
func (r *BloggerRepository) GetByXhsID(xhsID string) (*model.Blogger, error) {
	var blogger model.Blogger
	err := r.db.Where("xhs_id = ?", xhsID).First(&blogger).Error
	if err != nil {
		return nil, err
	}
	return &blogger, nil
}

// List 获取博主列表
func (r *BloggerRepository) List(offset, limit int) ([]*model.Blogger, int64, error) {
	var bloggers []*model.Blogger
	var total int64

	// 计算总数
	if err := r.db.Model(&model.Blogger{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按粉丝数排序
	err := r.db.Order("followers_count DESC").
		Offset(offset).
		Limit(limit).
		Find(&bloggers).Error

	return bloggers, total, err
}

// Update 更新博主信息
func (r *BloggerRepository) Update(blogger *model.Blogger) error {
	return r.db.Save(blogger).Error
}

// Delete 删除博主信息
func (r *BloggerRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Blogger{}).Error
}

// BatchCreate 批量创建博主信息
func (r *BloggerRepository) BatchCreate(bloggers []*model.Blogger) error {
	if len(bloggers) == 0 {
		return nil
	}
	return r.db.Create(&bloggers).Error
}

// UpsertByXhsID 根据 xhs_id 插入或更新博主信息
func (r *BloggerRepository) UpsertByXhsID(blogger *model.Blogger) error {
	var existing model.Blogger
	err := r.db.Where("xhs_id = ?", blogger.XhsID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		return r.db.Create(blogger).Error
	} else if err != nil {
		return err
	}

	// 存在则更新
	blogger.ID = existing.ID
	return r.db.Save(blogger).Error
}

// Count 获取博主总数
func (r *BloggerRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Blogger{}).Count(&count).Error
	return count, err
}
