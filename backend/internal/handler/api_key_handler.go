package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/keenchase/edit-business/internal/model"
	"github.com/keenchase/edit-business/internal/service"
)

// APIKeyHandler handles API key HTTP requests
type APIKeyHandler struct {
	apiKeyService *service.APIKeyService
}

// NewAPIKeyHandler creates a new API key handler
func NewAPIKeyHandler(apiKeyService *service.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// Create handles API key creation
// @Summary Create API key
// @Description Create a new API key for plugin authentication
// @Tags api-keys
// @Accept json
// @Produce json
// @Param request body service.CreateAPIKeyRequest true "Create API key request"
// @Success 200 {object} handler.Response{data=service.APIKeyResponse}
// @Failure 400 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys [post]
func (h *APIKeyHandler) Create(c *gin.Context) {
	var req service.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "Invalid request",
		})
		return
	}

	// Get auth center user ID from context (set by auth middleware)
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	apiKey, err := h.apiKeyService.Create(authCenterUserID.(string), req)
	if err != nil {
		if err == service.ErrMaxAPIKeysReached {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "Maximum API keys limit reached (max 10)",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to create API key",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "API key created successfully",
		Data:    apiKey,
	})
}

// List handles listing all API keys for the current user
// @Summary List API keys
// @Description Get all API keys for the current user
// @Tags api-keys
// @Produce json
// @Success 200 {object} handler.Response{data=[]service.APIKeyResponse}
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys [get]
func (h *APIKeyHandler) List(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	apiKeys, err := h.apiKeyService.List(authCenterUserID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to list API keys",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    apiKeys,
	})
}

// GetOrCreate handles getting or creating API key for the user
// @Summary Get or create API key
// @Description Get existing API key or create a new one automatically
// @Tags api-keys
// @Produce json
// @Success 200 {object} handler.Response{data=service.APIKeyResponse}
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys/get-or-create [get]
func (h *APIKeyHandler) GetOrCreate(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	// 优先使用 context 中的 user（AuthCenterMiddleware 已创建/获取），避免新用户竞态
	var apiKey *service.APIKeyResponse
	var err error
	if userVal, ok := c.Get("user"); ok && userVal != nil {
		if user, ok := userVal.(*model.User); ok && user != nil {
			apiKey, err = h.apiKeyService.GetOrCreateAPIKeyByUser(user)
		}
	}
	if apiKey == nil {
		apiKey, err = h.apiKeyService.GetOrCreateAPIKey(authCenterUserID.(string))
	}
	if err != nil {
		if err == service.ErrAPIKeyNotFound {
			c.JSON(http.StatusNotFound, Response{
				Code:    404,
				Message: "尚未分配 API Key，请联系管理员",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to get API key",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    apiKey,
	})
}

// GetStats handles getting API key statistics
// @Summary Get API key statistics
// @Description Get statistics about API keys usage
// @Tags api-keys
// @Produce json
// @Success 200 {object} handler.Response{data=repository.APIKeyStats}
// @Failure 401 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys/stats [get]
func (h *APIKeyHandler) GetStats(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	stats, err := h.apiKeyService.GetStats(authCenterUserID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to get API key statistics",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "Success",
		Data:    stats,
	})
}

// Delete handles deleting an API key
// @Summary Delete API key
// @Description Delete an API key by ID
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys/{id} [delete]
func (h *APIKeyHandler) Delete(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "API Key ID is required",
		})
		return
	}

	err := h.apiKeyService.Delete(authCenterUserID.(string), apiKeyID)
	if err != nil {
		if err == service.ErrInvalidAPIKey {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "Invalid API Key ID or permission denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to delete API key",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "API key deleted successfully",
	})
}

// Deactivate handles deactivating an API key
// @Summary Deactivate API key
// @Description Deactivate an API key by ID
// @Tags api-keys
// @Produce json
// @Param id path string true "API Key ID"
// @Success 200 {object} handler.Response
// @Failure 400 {object} handler.Response
// @Failure 401 {object} handler.Response
// @Failure 404 {object} handler.Response
// @Failure 500 {object} handler.Response
// @Router /api/v1/api-keys/{id}/deactivate [patch]
func (h *APIKeyHandler) Deactivate(c *gin.Context) {
	authCenterUserID, exists := c.Get("authCenterUserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	apiKeyID := c.Param("id")
	if apiKeyID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "API Key ID is required",
		})
		return
	}

	err := h.apiKeyService.Deactivate(authCenterUserID.(string), apiKeyID)
	if err != nil {
		if err == service.ErrInvalidAPIKey {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "Invalid API Key ID or permission denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "Failed to deactivate API key",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "API key deactivated successfully",
	})
}

// Validate handles validating an API key
// @Summary Validate API key
// @Description Validate if an API key is valid
// @Tags api-keys
// @Produce json
// @Success 200 {object} handler.Response{data=map[string]string}
// @Failure 401 {object} handler.Response
// @Router /api/v1/api-keys/validate [get]
func (h *APIKeyHandler) Validate(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "Unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Code:    0,
		Message: "API key is valid",
		Data: map[string]string{
			"userId": userID.(string),
		},
	})
}

// ValidateAPIKeyMiddleware validates API key from X-API-Key header
// This is used for plugin authentication
func (h *APIKeyHandler) ValidateAPIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header first (for normal JWT auth)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			// Let the JWT middleware handle this
			c.Next()
			return
		}

		// Check X-API-Key header (for plugin authentication)
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "Unauthorized: Missing API key or token",
			})
			c.Abort()
			return
		}

		// Validate API key and get both userId and authCenterUserID
		userID, authCenterUserID, err := h.apiKeyService.ValidateAPIKeyWithAuthCenterID(apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, Response{
				Code:    401,
				Message: "Unauthorized: Invalid API key",
			})
			c.Abort()
			return
		}

		// Set user IDs in context for downstream handlers
		c.Set("userId", userID)
		c.Set("authCenterUserID", authCenterUserID)
		c.Set("authType", "api_key")
		c.Next()
	}
}
