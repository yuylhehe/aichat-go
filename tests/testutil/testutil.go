package testutil

import (
	"ai-chat/internal/model"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB 创建测试用的内存数据库
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// 自动迁移
	err = db.AutoMigrate(
		&model.User{},
		&model.Conversation{},
		&model.Message{},
		&model.FixedPrompt{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	return db
}

// CreateTestUser 创建测试用户
func CreateTestUser(db *gorm.DB, name, email, password, salt string) *model.User {
	user := &model.User{
		Name:     name,
		Email:    email,
		Password: password,
		Salt:     salt,
		IsActive: true,
	}
	db.Create(user)
	return user
}

// CreateTestConversation 创建测试会话
func CreateTestConversation(db *gorm.DB, name string, userID uint) *model.Conversation {
	conv := &model.Conversation{
		Name:     name,
		UserID:   userID,
		IsActive: true,
	}
	db.Create(conv)
	return conv
}

// CreateTestMessage 创建测试消息
func CreateTestMessage(db *gorm.DB, conversationID uint, content, msgType string, sort int) *model.Message {
	msg := &model.Message{
		ConversationID: conversationID,
		Content:        content,
		Type:           msgType,
		Sort:           sort,
	}
	db.Create(msg)
	return msg
}

// CreateTestFixedPrompt 创建测试固定提示词
func CreateTestFixedPrompt(db *gorm.DB, userID uint, name, content string) *model.FixedPrompt {
	prompt := &model.FixedPrompt{
		UserID:   userID,
		Name:     name,
		Content:  content,
		IsActive: true,
	}
	db.Create(prompt)
	return prompt
}

// TimePtr 返回时间指针
func TimePtr(t time.Time) *time.Time {
	return &t
}

// StringPtr 返回字符串指针
func StringPtr(s string) *string {
	return &s
}

// UintPtr 返回 uint 指针
func UintPtr(n uint) *uint {
	return &n
}

// IntPtr 返回 int 指针
func IntPtr(n int) *int {
	return &n
}

// Float64Ptr 返回 float64 指针
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr 返回 bool 指针
func BoolPtr(b bool) *bool {
	return &b
}
