package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AuthCenterService 账号中心认证服务
type AuthCenterService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAuthCenterService 创建账号中心服务
func NewAuthCenterService() *AuthCenterService {
	return &AuthCenterService{
		BaseURL: "https://os.crazyaigc.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VerifyTokenRequest 验证 Token 请求
type VerifyTokenRequest struct {
	Token string `json:"token"`
}

// VerifyTokenResponse 验证 Token 响应
type VerifyTokenResponse struct {
	Success bool `json:"success"`
	Data    struct {
		UserID string `json:"userId"`
		UnionID string `json:"unionId"`
	} `json:"data"`
	Error interface{} `json:"error,omitempty"`
}

// VerifyToken 验证 Token
func (s *AuthCenterService) VerifyToken(token string) (string, error) {
	reqBody := VerifyTokenRequest{Token: token}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("序列化请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", s.BaseURL+"/api/auth/verify-token", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var result VerifyTokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("token 无效")
	}

	return result.Data.UserID, nil
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	Success bool `json:"success"`
	Data    struct {
		UserID      string `json:"userId"`
		UnionID     string `json:"unionId"`
		PhoneNumber string `json:"phoneNumber"`
		Email       string `json:"email"`
		CreatedAt   string `json:"createdAt"`
		LastLoginAt string `json:"lastLoginAt"`
		Profile     struct {
			Nickname  string `json:"nickname"`
			AvatarURL string `json:"avatarUrl"`
		} `json:"profile"`
		Accounts []struct {
			Provider   string `json:"provider"`
			Type       string `json:"type"`
			Nickname   string `json:"nickname"`
			AvatarURL  string `json:"avatarUrl"`
			CreatedAt  string `json:"createdAt"`
		} `json:"accounts"`
	} `json:"data"`
}

// GetUserInfoFromToken 用 token 获取账号中心的用户信息
func (s *AuthCenterService) GetUserInfoFromToken(token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", s.BaseURL+"/api/auth/user-info", nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result UserInfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 转换为 map 返回
	userInfo := map[string]interface{}{
		"userId":      result.Data.UserID,
		"unionId":     result.Data.UnionID,
		"phoneNumber": result.Data.PhoneNumber,
		"email":       result.Data.Email,
		"nickname":    result.Data.Profile.Nickname,
		"avatarUrl":   result.Data.Profile.AvatarURL,
	}

	return userInfo, nil
}
