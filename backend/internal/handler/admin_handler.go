package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/keenchase/edit-business/internal/service"
)

// AdminHandler 管理后台 HTTP 处理器
type AdminHandler struct {
	adminService         *service.AdminService
	adminAuthCenterIDs   []string
}

// NewAdminHandler 创建 Admin Handler
func NewAdminHandler(adminService *service.AdminService, adminAuthCenterIDs []string) *AdminHandler {
	return &AdminHandler{
		adminService:       adminService,
		adminAuthCenterIDs: adminAuthCenterIDs,
	}
}

// CheckAdmin 检查当前用户是否为管理员（供前端判断是否显示管理后台入口）
func (h *AdminHandler) CheckAdmin(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists || authCenterUserID == nil {
		c.JSON(200, gin.H{"isAdmin": false})
		return
	}
	id, _ := authCenterUserID.(string)
	for _, adminID := range h.adminAuthCenterIDs {
		if adminID == id {
			c.JSON(200, gin.H{"isAdmin": true})
			return
		}
	}
	c.JSON(200, gin.H{"isAdmin": false})
}

// ListUsers 分页获取所有用户（含采集统计）
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	res, err := h.adminService.ListUsers(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取用户列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    res,
	})
}

// GetUserDetail 获取用户详情
func (h *AdminHandler) GetUserDetail(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "用户ID不能为空",
		})
		return
	}
	detail, err := h.adminService.GetUserDetail(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    detail,
	})
}

// CreateAPIKeyForUserRequest 管理员创建 API Key 请求
type CreateAPIKeyForUserRequest struct {
	ExpiresIn *int `json:"expiresIn"` // 有效期天数，nil 或 0 表示永不过期
}

// CreateAPIKeyForUser 管理员为用户创建 API Key
func (h *AdminHandler) CreateAPIKeyForUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "用户ID不能为空",
		})
		return
	}
	var req CreateAPIKeyForUserRequest
	_ = c.ShouldBindJSON(&req) // 可选，无 body 时 expiresIn 为 nil
	apiKey, err := h.adminService.CreateAPIKeyForUser(userID, req.ExpiresIn)
	if err != nil {
		if err == service.ErrMaxAPIKeysReached {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "该用户已有 API Key",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建 API Key 失败",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "API Key 创建成功",
		Data:    apiKey,
	})
}

// UpdateAPIKeyExpiryRequest 修改 API Key 有效期请求
type UpdateAPIKeyExpiryRequest struct {
	ExpiresIn *int `json:"expiresIn"` // 有效期天数，nil 或 0 表示永不过期
}

// UpdateAPIKeyExpiry 管理员修改 API Key 有效期
func (h *AdminHandler) UpdateAPIKeyExpiry(c *gin.Context) {
	userID := c.Param("id")
	apiKeyID := c.Param("keyId")
	if userID == "" || apiKeyID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "用户ID和API Key ID不能为空",
		})
		return
	}
	var req UpdateAPIKeyExpiryRequest
	_ = c.ShouldBindJSON(&req)
	if err := h.adminService.UpdateAPIKeyExpiry(userID, apiKeyID, req.ExpiresIn); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "已更新有效期",
	})
}

// UpdateUserSettingsRequest 更新用户设置请求（仅 dailyLimit、batchLimit，允许数据收藏由用户自主控制）
type UpdateUserSettingsRequest struct {
	CollectionDailyLimit *int `json:"collectionDailyLimit"`
	CollectionBatchLimit *int `json:"collectionBatchLimit"`
}

// UpdateUserSettings 管理员修改用户采集限额（每日限额、单次限额）
func (h *AdminHandler) UpdateUserSettings(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "用户ID不能为空",
		})
		return
	}
	var req UpdateUserSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}
	if req.CollectionDailyLimit == nil && req.CollectionBatchLimit == nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请至少提供 collectionDailyLimit 或 collectionBatchLimit",
		})
		return
	}
	if err := h.adminService.UpdateUserSettings(userID, req.CollectionDailyLimit, req.CollectionBatchLimit, nil); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新失败",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "更新成功",
	})
}

// GetStatsOverview 全局统计（总用户数、总采集量等）
func (h *AdminHandler) GetStatsOverview(c *gin.Context) {
	overview, err := h.adminService.GetStatsOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取统计失败",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    overview,
	})
}
