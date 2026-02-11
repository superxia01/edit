package config

import (
	"fmt"
	"os"
	"strings"
)

// Config 应用配置
// 遵循 KeenChase V3.0 规范：环境变量管理
type Config struct {
	// 服务器配置
	ServerPort string
	ServerHost string

	// 数据库配置（通过 SSH 隧道）
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// 账号中心配置
	AuthCenterURL string

	// 管理后台：管理员 auth_center_user_id 列表，逗号分隔
	AdminAuthCenterUserIDs []string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),

		// KeenChase V3.0 标准：通过 SSH 隧道连接数据库
		// 隧道配置: localhost:5432 -> 47.110.82.96:5432
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "nexus_user"),
		DBPassword: getEnv("DB_PASSWORD", "hRJ9NSJApfeyFDraaDgkYowY"),
		DBName:     getEnv("DB_NAME", "edit_business_db"),

		// 账号中心 URL
		AuthCenterURL: getEnv("AUTH_CENTER_URL", "https://os.crazyaigc.com"),

		// 管理后台：EDIT_ADMIN_AUTH_CENTER_USER_IDS=id1,id2,id3
		AdminAuthCenterUserIDs: parseAdminIDs(getEnv("EDIT_ADMIN_AUTH_CENTER_USER_IDS", "")),
	}
}

// parseAdminIDs 解析逗号分隔的 admin IDs
func parseAdminIDs(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var ids []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			ids = append(ids, p)
		}
	}
	return ids
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
	)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
