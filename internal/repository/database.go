package repository

import (
	"ai-chat/internal/model"
	"fmt"
	"log"
	"time"

	"ai-chat/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB 创建数据库连接
func NewDB(cfg *config.Config) (*gorm.DB, error) {
	// 使用配置中的数据库连接字符串
	dsn := cfg.DatabaseURL()

	log.Printf("Attempting to connect to database with DSN: %s", dsn)

	// 配置连接参数
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	var db *gorm.DB
	var err error

	// 重试连接3次
	for i := 1; i <= 3; i++ {
		log.Printf("Database connection attempt %d/3", i)

		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		log.Printf("Connection attempt %d failed: %v", i, err)
		if i < 3 {
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after 3 attempts: %w", err)
	}

	// 获取通用数据库对象设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	// 测试数据库连接
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return db, nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&model.User{},
		&model.Conversation{},
		&model.Message{},
		&model.FixedPrompt{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
