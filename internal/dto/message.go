package dto

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	ConversationID uint   `json:"conversationId" binding:"required"`
	Content        string `json:"content" binding:"required,min=1"`
	Type           string `json:"type" binding:"required,oneof=system user assistant"`
}

// UpdateMessageRequest 更新消息请求
type UpdateMessageRequest struct {
	Content *string `json:"content,omitempty"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID               uint    `json:"id"`
	ConversationID   uint    `json:"conversationId"`
	Content          string  `json:"content"`
	ReasoningContent string  `json:"reasoningContent,omitempty"`
	Sort             int     `json:"sort"`
	Type             string  `json:"type"`
	Tokens           *int    `json:"tokens,omitempty"`
	Model            *string `json:"model,omitempty"`
	ParentID         *uint   `json:"parentId,omitempty"`
	Metadata         *string `json:"metadata,omitempty"`
	CreatedAt        string  `json:"createdAt"`
}

