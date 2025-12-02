package model

import (
	"time"

	"gorm.io/gorm"
)

// Message 消息模型
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

	// 关联关系
	Conversation Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationID"`
	Parent       *Message     `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies      []Message    `json:"replies,omitempty" gorm:"foreignKey:ParentID"`

	TableName string `json:"-" gorm:"tableName:message"`
}

// BeforeCreate 创建前钩子
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.Sort == 0 {
		// 获取当前会话的最大sort值并加1
		var maxSort int
		result := tx.Model(&Message{}).
			Where("conversation_id = ?", m.ConversationID).
			Pluck("COALESCE(MAX(sort), 0)", &maxSort)
		if result.Error == nil {
			m.Sort = maxSort + 1
		} else {
			m.Sort = 1
		}
	}
	return nil
}
