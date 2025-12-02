package dto

// CreateFixedPromptRequest 创建固定提示词请求
type CreateFixedPromptRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=255"`
	Content string `json:"content" binding:"required"`
}

// UpdateFixedPromptRequest 更新固定提示词请求
type UpdateFixedPromptRequest struct {
	Name     *string `json:"name,omitempty"`
	Content  *string `json:"content,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}

// FixedPromptResponse 固定提示词响应
type FixedPromptResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// GetFixedPromptsRequest 获取固定提示词列表请求
type GetFixedPromptsRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100" example:"10"`
	Search   string `form:"search"`
}

// PaginatedFixedPrompts 分页固定提示词列表
type PaginatedFixedPrompts struct {
	Items      []*FixedPromptResponse `json:"items"`
	Total      int64                  `json:"total"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"pageSize"`
	TotalPages int64                  `json:"totalPages"`
}
