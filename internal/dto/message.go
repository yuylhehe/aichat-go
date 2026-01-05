package dto

// MessageResponse 消息响应
type MessageResponse struct {
	ID               uint        `json:"id"`
	ConversationID   uint        `json:"conversationId"`
	Content          string      `json:"content"`
	ReasoningContent string      `json:"reasoningContent"`
	Sort             int         `json:"sort"`
	Type             string      `json:"type"`
	Tokens           int         `json:"tokens"`
	Model            string      `json:"model"`
	ParentID         *uint       `json:"parentId"`
	Metadata         interface{} `json:"metadata"`
	CreatedAt        string      `json:"createdAt"`
}

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	ConversationID   uint    `json:"conversationId" binding:"required"`
	Content          string  `json:"content" binding:"required"`
	ReasoningContent string  `json:"reasoningContent,omitempty"`
	Type             string  `json:"type" binding:"required,oneof=system user assistant"`
	Model            *string `json:"model,omitempty"`
	ParentID         *uint   `json:"parentId,omitempty"`
}

// UpdateMessageRequest 更新消息请求
type UpdateMessageRequest struct {
	Content          *string `json:"content,omitempty"`
	ReasoningContent *string `json:"reasoningContent,omitempty"`
	Sort             *int    `json:"sort,omitempty"`
	Type             *string `json:"type,omitempty"`
	Tokens           *int    `json:"tokens,omitempty"`
	Model            *string `json:"model,omitempty"`
	ParentID         *uint   `json:"parentId,omitempty"`
	Metadata         *string `json:"metadata,omitempty"`
	CreatedAt        string  `json:"createdAt"`
}
