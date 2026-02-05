package database

import (
	"fmt"
	"log"

	"github.com/keenchase/edit-business/internal/config"
	"github.com/keenchase/edit-business/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
// 遵循 KeenChase V3.0 规范：通过 SSH 隧道连接 PostgreSQL
func InitDatabase(cfg *config.Config) error {
	dsn := cfg.GetDSN()

	// 配置 GORM
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		// NamingStrategy: 默认使用 NamingStrategy，会在表名和字段名之间自动转换
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("Database connection established successfully")
	return nil
}

// AutoMigrate 自动迁移数据库表结构
// 开发环境使用，生产环境应使用迁移脚本
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}

	// 迁移所有模型
	err := DB.AutoMigrate(
		&model.User{},
		&model.Note{},
		&model.Blogger{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase 关闭数据库连接
func CloseDatabase() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
