package dto

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	Name string `json:"name" binding:"required,min=1,max=255"`
}

// UpdateConversationRequest 更新对话请求
type UpdateConversationRequest struct {
	Name *string `json:"name,omitempty"`
}

// ConversationResponse 对话响应
type ConversationResponse struct {
	ID           uint     `json:"id"`
	Name         string   `json:"name"`
	UserID       uint     `json:"userId"`
	IsActive     bool     `json:"isActive"`
	SystemPrompt *string  `json:"systemPrompt,omitempty"`
	Model        *string  `json:"model,omitempty"`
	Temperature  *float64 `json:"temperature,omitempty"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
	Messages     int64    `json:"messageCount"`
}

