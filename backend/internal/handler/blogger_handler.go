package handler

import (
	"strconv"

	"github.com/keenchase/edit-business/internal/service"
	"github.com/gin-gonic/gin"
)

// BloggerHandler 博主处理器
type BloggerHandler struct {
	bloggerService *service.BloggerService
}

// NewBloggerHandler 创建博主处理器实例
func NewBloggerHandler(bloggerService *service.BloggerService) *BloggerHandler {
	return &BloggerHandler{bloggerService: bloggerService}
}

// Create 创建博主信息
// @Summary 创建博主信息
// @Description 创建新的博主记录
// @Tags bloggers
// @Accept json
// @Produce json
// @Param request body service.CreateBloggerRequest true "创建博主请求"
// @Success 200 {object} Response
// @Router /api/v1/bloggers [post]
func (h *BloggerHandler) Create(c *gin.Context) {
	var req service.CreateBloggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Get authCenterUserID from context (set by auth middleware)
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(401, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	blogger, err := h.bloggerService.Create(authCenterUserID.(string), &req)
	if err != nil {
		if err == service.ErrCollectionDisabled {
			c.JSON(403, Response{
				Code:    403,
				Message: "采集功能已关闭，请在网站设置中开启",
			})
			return
		}
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, blogger)
}

// GetByID 根据 ID 获取博主信息
// @Summary 获取博主详情
// @Description 根据 ID 获取博主详情
// @Tags bloggers
// @Accept json
// @Produce json
// @Param id path string true "博主 ID"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/{id} [get]
func (h *BloggerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "id is required")
		return
	}

	blogger, err := h.bloggerService.GetByID(id)
	if err != nil {
		NotFound(c, "blogger not found")
		return
	}

	SuccessResponse(c, blogger)
}

// GetByXhsID 根据小红书 ID 获取博主信息
// @Summary 根据 xhs_id 获取博主详情
// @Description 根据小红书 ID 获取博主详情
// @Tags bloggers
// @Accept json
// @Produce json
// @Param xhsId path string true "小红书 ID"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/xhs/{xhsId} [get]
func (h *BloggerHandler) GetByXhsID(c *gin.Context) {
	xhsID := c.Param("xhsId")
	if xhsID == "" {
		BadRequest(c, "xhsId is required")
		return
	}

	blogger, err := h.bloggerService.GetByXhsID(xhsID)
	if err != nil {
		NotFound(c, "blogger not found")
		return
	}

	SuccessResponse(c, blogger)
}

// List 获取博主列表
// @Summary 获取博主列表
// @Description 分页获取博主列表，按粉丝数排序
// @Tags bloggers
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(20)
// @Success 200 {object} Response
// @Router /api/v1/bloggers [get]
func (h *BloggerHandler) List(c *gin.Context) {
	var req service.ListBloggersRequest

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil {
			req.Page = page
		}
	}

	if sizeStr := c.Query("size"); sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err == nil {
			req.Size = size
		}
	}

	// 调用服务层
	result, err := h.bloggerService.List(&req)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, result)
}

// BatchCreate 批量创建博主信息
// @Summary 批量创建博主信息
// @Description 批量创建博主记录（用于 Chrome 插件同步）
// @Tags bloggers
// @Accept json
// @Produce json
// @Param request body []service.CreateBloggerRequest true "批量创建博主请求"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/batch [post]
func (h *BloggerHandler) BatchCreate(c *gin.Context) {
	var reqs []*service.CreateBloggerRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Get authCenterUserID from context (set by auth middleware)
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(401, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	err := h.bloggerService.BatchCreate(authCenterUserID.(string), reqs)
	if err != nil {
		if err == service.ErrCollectionDisabled {
			c.JSON(403, Response{
				Code:    403,
				Message: "采集功能已关闭，请在网站设置中开启",
			})
			return
		}
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"count":  len(reqs),
		"status": "success",
	})
}

// UpsertByXhsID 根据 xhs_id 插入或更新博主信息
// @Summary 插入或更新博主信息
// @Description 根据小红书 ID 插入或更新博主信息
// @Tags bloggers
// @Accept json
// @Produce json
// @Param request body service.CreateBloggerRequest true "博主信息"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/upsert [post]
func (h *BloggerHandler) UpsertByXhsID(c *gin.Context) {
	var req service.CreateBloggerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// Get authCenterUserID from context (set by auth middleware)
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(401, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	blogger, err := h.bloggerService.UpsertByXhsID(authCenterUserID.(string), &req)
	if err != nil {
		if err == service.ErrCollectionDisabled {
			c.JSON(403, Response{
				Code:    403,
				Message: "采集功能已关闭，请在网站设置中开启",
			})
			return
		}
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, blogger)
}

// Update 更新博主信息
// @Summary 更新博主信息
// @Description 更新博主信息
// @Tags bloggers
// @Accept json
// @Produce json
// @Param id path string true "博主 ID"
// @Param request body model.Blogger true "博主信息"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/{id} [put]
func (h *BloggerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "id is required")
		return
	}

	// 获取现有博主
	blogger, err := h.bloggerService.GetByID(id)
	if err != nil {
		NotFound(c, "blogger not found")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// 更新字段（这里简化处理）
	_ = req

	err = h.bloggerService.Update(blogger)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, blogger)
}

// Delete 删除博主信息
// @Summary 删除博主信息
// @Description 根据 ID 删除博主
// @Tags bloggers
// @Accept json
// @Produce json
// @Param id path string true "博主 ID"
// @Success 200 {object} Response
// @Router /api/v1/bloggers/{id} [delete]
func (h *BloggerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		BadRequest(c, "id is required")
		return
	}

	err := h.bloggerService.Delete(id)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	SuccessResponse(c, gin.H{
		"id":     id,
		"status": "deleted",
	})
}
