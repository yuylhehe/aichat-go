package repository

import (
	"time"

	"gorm.io/gorm"
)

// Message 消息数据库模型
type Message struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	ConversationID   uint           `json:"conversationId" gorm:"not null;index"`
	Content          string         `json:"content" gorm:"type:text;not null"`
	ReasoningContent string         `json:"reasoningContent" gorm:"type:text"` // 新增思考内容字段
	Sort             int            `json:"sort" gorm:"not null"`
	Type             string         `json:"type" gorm:"size:20;not null;index"` // system user assistant
	Tokens           *int           `json:"tokens" gorm:"index"`
	Model            *string        `json:"model" gorm:"size:100"`
	ParentID         *uint          `json:"parentId" gorm:"index"`
	Metadata         *string        `json:"metadata" gorm:"type:jsonb"`
	CreatedAt        time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	TableName string `json:"-" gorm:"tableName:message"`
}
