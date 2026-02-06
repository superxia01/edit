package middleware

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// generateInvitationCode 生成8位随机邀请码
func generateInvitationCode() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 8)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}

// AuthCenterMiddlewareInterface 账号中心服务接口（避免循环依赖）
type AuthCenterService interface {
	VerifyToken(token string) (string, error)
	GetUserInfoFromToken(token string) (map[string]interface{}, error)
}

// UserRepositoryInterface 用户仓库接口
type UserRepository interface {
	GetByAuthCenterUserIDInterface(authCenterUserID string) (interface{}, error)
	CreateWithAuthCenter(authCenterUserID string, unionID string, nickname string, avatarURL string) (interface{}, error)
}

// 为了避免循环依赖，这里使用 interface{} 作为返回类型
// 实际使用时需要类型断言：user.(*model.User)

// AuthCenterMiddleware 账号中心认证中间件
func AuthCenterMiddleware(authService AuthCenterService, userRepo UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未登录",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token 格式错误",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// 1. 调用账号中心验证 token
		authCenterUserID, err := authService.VerifyToken(token)
		if err != nil {
			fmt.Printf("[AuthCenterMiddleware] VerifyToken failed: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Token 无效或已过期",
			})
			c.Abort()
			return
		}

		fmt.Printf("[AuthCenterMiddleware] AuthCenterUserID: %s\n", authCenterUserID)

		// 2. 获取本地用户
		localUser, err := userRepo.GetByAuthCenterUserIDInterface(authCenterUserID)
		if err != nil {
			// 本地用户不存在，从 auth-center 获取并创建
			fmt.Printf("[AuthCenterMiddleware] Local user not found, creating from auth-center: %v\n", err)

			// 调用 auth-center 获取用户信息
			userInfo, err := authService.GetUserInfoFromToken(token)
			if err != nil {
				fmt.Printf("[AuthCenterMiddleware] GetUserInfoFromToken failed: %v\n", err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "获取账号中心用户信息失败",
				})
				c.Abort()
				return
			}

			// 提取用户信息
			unionID := getStringValue(userInfo, "unionId")
			nickname := getStringValue(userInfo, "nickname")
			avatarURL := getStringValue(userInfo, "avatarUrl")

			// 创建本地用户
			localUser, err = userRepo.CreateWithAuthCenter(authCenterUserID, unionID, nickname, avatarURL)
			if err != nil {
				fmt.Printf("[AuthCenterMiddleware] Create local user failed: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "创建本地用户失败",
				})
				c.Abort()
				return
			}

			fmt.Printf("[AuthCenterMiddleware] Created new local user\n")
		}

		// 3. 将用户信息存入上下文
		c.Set("user", localUser)
		c.Set("authCenterUserID", authCenterUserID)
		c.Set("userId", authCenterUserID)

		c.Next()
	}
}

// getStringValue 从 map 中安全获取字符串值
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
