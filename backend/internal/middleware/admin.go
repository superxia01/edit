package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 校验当前登录用户是否为管理员
// 必须在 AuthCenterMiddleware 之后使用，依赖 authCenterUserID
func AdminMiddleware(adminIDs []string) gin.HandlerFunc {
	adminSet := make(map[string]bool)
	for _, id := range adminIDs {
		adminSet[id] = true
	}

	return func(c *gin.Context) {
		authCenterUserID, exists := c.Get("authCenterUserID")
		if !exists || authCenterUserID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "未登录",
			})
			c.Abort()
			return
		}

		id, ok := authCenterUserID.(string)
		if !ok || id == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "无效用户",
			})
			c.Abort()
			return
		}

		if !adminSet[id] {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "无管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
