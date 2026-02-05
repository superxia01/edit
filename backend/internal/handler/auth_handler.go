package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/keenchase/edit-business/internal/middleware"
	"github.com/keenchase/edit-business/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	Token string `json:"token" binding:"required"`
}

// WechatLogin 微信登录（账号中心集成）
// @Summary 微信登录
// @Description 通过账号中心进行微信登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body WechatLoginRequest true "登录请求"
// @Success 200 {object} Response
// @Router /api/v1/auth/wechat [post]
func (h *AuthHandler) WechatLogin(c *gin.Context) {
	var req WechatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	// TODO: 调用账号中心验证 token
	// 暂时简化处理：直接使用 token 作为 userId
	// 实际应该调用账号中心的 API 验证 token 并获取 userId

	authCenterUserID := req.Token

	// 同步或创建用户
	user, err := h.userService.SyncUserFromAuthCenter(authCenterUserID)
	if err != nil {
		InternalError(c, err.Error())
		return
	}

	// 生成 JWT token
	jwtToken, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		InternalError(c, "生成令牌失败")
		return
	}

	SuccessResponse(c, gin.H{
		"token": jwtToken,
		"user":  user,
	})
}

// WechatCallback 微信登录回调（V3.0 标准）
// @Summary 微信登录回调
// @Description 处理从账号中心回调的微信登录请求
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "微信授权码"
// @Param type query string true "登录类型 (open/mp)"
// @Success 200 {object} Response
// @Router /api/v1/auth/wechat/callback [get]
func (h *AuthHandler) WechatCallback(c *gin.Context) {
	code := c.Query("code")
	loginType := c.Query("type") // "open" (开放平台) 或 "mp" (公众号)

	if code == "" {
		BadRequest(c, "缺少code参数")
		return
	}

	// 调用账号中心的微信登录API，用code换取token
	// POST https://os.crazyaigc.com/api/auth/wechat/login
	// Body: {"code": "xxx", "type": "open"}
	loginReqBody, _ := json.Marshal(map[string]string{
		"code": code,
		"type": loginType,
	})

	loginResp, err := http.Post(
		"https://os.crazyaigc.com/api/auth/wechat/login",
		"application/json",
		bytes.NewBuffer(loginReqBody),
	)
	if err != nil {
		InternalError(c, "调用账号中心失败")
		return
	}
	defer loginResp.Body.Close()

	var loginResult struct {
		Success bool `json:"success"`
		Data    struct {
			UserID      string `json:"userId"`
			Token       string `json:"token"`
			UnionID     string `json:"unionId"`
			PhoneNumber string `json:"phoneNumber"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginResult); err != nil {
		InternalError(c, "解析响应失败")
		return
	}

	if !loginResult.Success {
		BadRequest(c, loginResult.Message)
		return
	}

	// 同步或创建用户
	user, err := h.userService.SyncUserFromAuthCenter(loginResult.Data.UserID)
	if err != nil {
		InternalError(c, "创建用户失败")
		return
	}

	// 生成 JWT token（本地系统的token）
	jwtToken, err := middleware.GenerateToken(user.ID.String())
	if err != nil {
		InternalError(c, "生成令牌失败")
		return
	}

	// 返回 JSON 给前端（前端会保存 token 并跳转）
	SuccessResponse(c, gin.H{
		"token":               jwtToken,
		"auth_center_user_id": loginResult.Data.UserID,
		"user":                user,
	})
}

// GetCurrentUser 获取当前登录用户信息
// @Summary 获取当前用户
// @Description 获取当前登录用户的信息
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		NotFound(c, "用户未登录")
		return
	}

	user, err := h.userService.GetByID(userID.(string))
	if err != nil {
		NotFound(c, "用户不存在")
		return
	}

	SuccessResponse(c, user)
}
