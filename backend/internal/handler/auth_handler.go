package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// WechatLoginProxy 发起微信登录（重定向到 auth-center）
// @Summary 发起微信登录
// @Description 重定向到账号中心进行微信登录
// @Tags auth
// @Accept json
// @Produce json
// @Param callbackUrl query string false "回调地址"
// @Success 200 {object} Response
// @Router /api/v1/auth/wechat/login [get]
func (h *AuthHandler) WechatLoginProxy(c *gin.Context) {
	// 构建回调 URL
	callbackURL := "https://edit.crazyaigc.com/api/v1/auth/wechat/callback"
	authCenterURL := fmt.Sprintf(
		"https://os.crazyaigc.com/api/auth/wechat/login?callbackUrl=%s",
		callbackURL,
	)

	// 重定向到账号中心
	c.Redirect(http.StatusFound, authCenterURL)
}

// WechatLoginRequest 微信登录请求
type WechatLoginRequest struct {
	AuthCode string `json:"authCode" binding:"required"`
	Type     string `json:"type"` // "open" (开放平台) 或 "mp" (公众号)
}

// WechatLogin PC扫码微信登录
// @Summary PC扫码微信登录
// @Description 使用 authCode 换取用户信息并登录
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

	// 调用账号中心的微信登录API，用 code 换取用户信息
	// POST https://os.crazyaigc.com/api/auth/wechat/login
	// Body: {"code": "xxx", "type": "open"}
	loginReqBody, _ := json.Marshal(map[string]string{
		"code": req.AuthCode,
		"type": req.Type,
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
			UserID      string                 `json:"userId"`
			Token       string                 `json:"token"`
			UnionID     string                 `json:"unionId"`
			PhoneNumber string                 `json:"phoneNumber"`
			Profile     map[string]interface{} `json:"profile"`
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

	// 从 data.profile 直接获取用户信息（PC扫码登录不需要再调用 user-info）
	var nicknameStr, avatarUrlStr string
	if loginResult.Data.Profile != nil {
		if val, ok := loginResult.Data.Profile["nickname"]; ok && val != nil {
			switch v := val.(type) {
			case string:
				nicknameStr = v
			case []byte:
				nicknameStr = string(v)
			}
		}
		if val, ok := loginResult.Data.Profile["avatarUrl"]; ok && val != nil {
			switch v := val.(type) {
			case string:
				avatarUrlStr = v
			case []byte:
				avatarUrlStr = string(v)
			}
		}
	}

	// 同步或创建用户
	user, err := h.userService.SyncUserFromAuthCenter(loginResult.Data.UserID, nicknameStr, avatarUrlStr)
	if err != nil {
		InternalError(c, "创建用户失败")
		return
	}

	// 生成 JWT token（本地系统的 token）
	jwtToken, err := middleware.GenerateToken(user.ID)
	if err != nil {
		InternalError(c, "生成令牌失败")
		return
	}

	SuccessResponse(c, gin.H{
		"token": jwtToken,
		"user":  user,
	})
}

// WechatCallback 微信登录回调（V3.1 统一 Token 模式）
// @Summary 微信登录回调
// @Description 接收账号中心回调的token，重定向到前端
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Auth-Center Token"
// @Success 302 {string} string "重定向到前端"
// @Router /api/v1/auth/wechat/callback [get]
func (h *AuthHandler) WechatCallback(c *gin.Context) {
	token := c.Query("token")

	// V3.1: 验证 token 参数
	if token == "" {
		c.Redirect(http.StatusFound, "/login?error=missing_token")
		return
	}

	// ✅ 直接重定向到前端页面，带上 token 参数
	// 前端会调用 /api/v1/user/me 来验证 token 并获取用户信息
	frontendURL := fmt.Sprintf("https://edit.crazyaigc.com/auth/callback?token=%s", token)
	c.Redirect(http.StatusFound, frontendURL)
}

// Me 获取当前用户信息（使用 AuthCenterMiddleware）
// @Summary 获取当前用户
// @Description 获取当前登录用户的信息，AuthCenterMiddleware 已经处理了认证
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /api/v1/user/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	// AuthCenterMiddleware 已经处理了：
	//   1. 验证 auth-center token
	//   2. 获取/创建本地用户
	//   3. 存入上下文 ("user" 和 "userId")

	// 直接从上下文获取用户信息
	user, exists := c.Get("user")
	if !exists {
		Unauthorized(c, "用户未登录")
		return
	}

	SuccessResponse(c, gin.H{
		"success": true,
		"data":    user,
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
