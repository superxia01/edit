package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// JWT Claims 结构
type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

// JWT secret key
var jwtSecret = []byte(getEnv("JWT_SECRET", "change-this-secret-in-production"))

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateToken 生成 JWT token
func GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7天过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// 使用 Gin 的日志输出
		c.Set("JWT_DEBUG_HasAuthHeader", authHeader != "")
		c.Set("JWT_DEBUG_AuthHeaderLength", len(authHeader))

		// 从 Authorization header 获取 token
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未提供认证令牌",
				"debug":   "No Authorization header",
			})
			c.Abort()
			return
		}

		// Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "令牌格式错误",
				"debug":   "Invalid Bearer token format",
			})
			c.Abort()
			return
		}

		// 解析 token
		claims, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "令牌无效或已过期 - DEBUG123",
				"debug":   fmt.Sprintf("Token parse error: %v", err),
			})
			c.Abort()
			return
		}

		// 将用户 ID 存入上下文（使用 authCenterUserID 以匹配 handlers）
		c.Set("authCenterUserID", claims.UserID)
		c.Set("userId", claims.UserID) // 为了兼容性，两个都设置
		c.Next()
	}
}

// OptionalJWTAuth 可选的 JWT 认证中间件（如果有 token 则验证，没有则跳过）
func OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := ParseToken(parts[1])
		if err == nil {
			c.Set("userId", claims.UserID)
		}

		c.Next()
	}
}
