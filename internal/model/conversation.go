package model

import (
	"time"

	"gorm.io/gorm"
)

// Conversation 会话模型
type Conversation struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name" gorm:"size:255;not null"`
	UserID       uint           `json:"userId" gorm:"not null;index"`
	IsActive     bool           `json:"isActive" gorm:"default:true"`
	SystemPrompt *string        `json:"systemPrompt" gorm:"type:text"`
	Model        *string        `json:"model" gorm:"size:100"`
	Temperature  *float64       `json:"temperature" gorm:"type:decimal(3,2)"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Messages []Message `json:"messages,omitempty" gorm:"foreignKey:ConversationID"`

	TableName string `json:"-" gorm:"tableName:conversation"`
}

// BeforeCreate 创建前钩子
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	c.IsActive = true
	return nil
}
