package service

import (
	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/internal/repository"
)

// BloggerService 博主服务
type BloggerService struct {
	bloggerRepo      *repository.BloggerRepository
	settingsService  *UserSettingsService
}

// NewBloggerService 创建博主服务实例
func NewBloggerService(bloggerRepo *repository.BloggerRepository, settingsService *UserSettingsService) *BloggerService {
	return &BloggerService{
		bloggerRepo:     bloggerRepo,
		settingsService: settingsService,
	}
}

// CreateBloggerRequest 创建博主请求
type CreateBloggerRequest struct {
	XhsID            string `json:"xhsId" binding:"required"`
	BloggerName      string `json:"bloggerName"`
	AvatarURL        string `json:"avatarUrl"`
	Description      string `json:"description"`
	FollowersCount   int32  `json:"followersCount"`
	BloggerURL       string `json:"bloggerUrl"`
	CaptureTimestamp int64  `json:"captureTimestamp" binding:"required"`
}

// ListBloggersRequest 列表查询请求
type ListBloggersRequest struct {
	Page int `form:"page" binding:"min=1"`
	Size int `form:"size" binding:"min=1,max=100"`
}

// ListBloggersResponse 列表查询响应
type ListBloggersResponse struct {
	Bloggers   []*model.Blogger `json:"bloggers"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
	TotalPages int              `json:"totalPages"`
}

// Create 创建博主信息
func (s *BloggerService) Create(authCenterUserID string, req *CreateBloggerRequest) (*model.Blogger, error) {
	// Check if collection is enabled
	enabled, err := s.settingsService.IsCollectionEnabled(authCenterUserID)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, ErrCollectionDisabled
	}

	// Get user
	user, err := s.settingsService.GetUserByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}

	blogger := &model.Blogger{
		UserID:           user.ID,
		XhsID:            req.XhsID,
		BloggerName:      req.BloggerName,
		AvatarURL:        req.AvatarURL,
		Description:      req.Description,
		FollowersCount:   req.FollowersCount,
		BloggerURL:       req.BloggerURL,
		CaptureTimestamp: req.CaptureTimestamp,
	}

	err = s.bloggerRepo.Create(blogger)
	if err != nil {
		return nil, err
	}

	return blogger, nil
}

// GetByID 根据 ID 获取博主信息
func (s *BloggerService) GetByID(id string) (*model.Blogger, error) {
	return s.bloggerRepo.GetByID(id)
}

// GetByXhsID 根据小红书 ID 获取博主信息
func (s *BloggerService) GetByXhsID(xhsID string) (*model.Blogger, error) {
	return s.bloggerRepo.GetByXhsID(xhsID)
}

// List 获取博主列表
func (s *BloggerService) List(req *ListBloggersRequest) (*ListBloggersResponse, error) {
	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Size < 1 {
		req.Size = 20
	}

	offset := (req.Page - 1) * req.Size

	bloggers, total, err := s.bloggerRepo.List(offset, req.Size)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / req.Size
	if int(total)%req.Size > 0 {
		totalPages++
	}

	return &ListBloggersResponse{
		Bloggers:   bloggers,
		Total:      total,
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}, nil
}

// UpsertByXhsID 根据 xhs_id 插入或更新博主信息
func (s *BloggerService) UpsertByXhsID(authCenterUserID string, req *CreateBloggerRequest) (*model.Blogger, error) {
	// Check if collection is enabled
	enabled, err := s.settingsService.IsCollectionEnabled(authCenterUserID)
	if err != nil {
		return nil, err
	}
	if !enabled {
		return nil, ErrCollectionDisabled
	}

	// Get user
	user, err := s.settingsService.GetUserByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return nil, err
	}

	blogger := &model.Blogger{
		UserID:           user.ID,
		XhsID:            req.XhsID,
		BloggerName:      req.BloggerName,
		AvatarURL:        req.AvatarURL,
		Description:      req.Description,
		FollowersCount:   req.FollowersCount,
		BloggerURL:       req.BloggerURL,
		CaptureTimestamp: req.CaptureTimestamp,
	}

	err = s.bloggerRepo.UpsertByXhsID(blogger)
	if err != nil {
		return nil, err
	}

	return blogger, nil
}

// BatchCreate 批量创建博主信息（用于 Chrome 插件同步）
func (s *BloggerService) BatchCreate(authCenterUserID string, reqs []*CreateBloggerRequest) error {
	if len(reqs) == 0 {
		return nil
	}

	// Check if collection is enabled
	enabled, err := s.settingsService.IsCollectionEnabled(authCenterUserID)
	if err != nil {
		return err
	}
	if !enabled {
		return ErrCollectionDisabled
	}

	// Get user
	user, err := s.settingsService.GetUserByAuthCenterUserID(authCenterUserID)
	if err != nil {
		return err
	}

	bloggers := make([]*model.Blogger, len(reqs))
	for i, req := range reqs {
		bloggers[i] = &model.Blogger{
			UserID:           user.ID,
			XhsID:            req.XhsID,
			BloggerName:      req.BloggerName,
			AvatarURL:        req.AvatarURL,
			Description:      req.Description,
			FollowersCount:   req.FollowersCount,
			BloggerURL:       req.BloggerURL,
			CaptureTimestamp: req.CaptureTimestamp,
		}
	}

	return s.bloggerRepo.BatchCreate(bloggers)
}

// Update 更新博主信息
func (s *BloggerService) Update(blogger *model.Blogger) error {
	return s.bloggerRepo.Update(blogger)
}

// Delete 删除博主信息
func (s *BloggerService) Delete(id string) error {
	return s.bloggerRepo.Delete(id)
}
