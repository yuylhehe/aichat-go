package dto

// UserResponse 用户响应
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Username string `json:"username" example:"john_doe"`
	Email    string `json:"email" binding:"email" example:"john@example.com"`
	Nickname string `json:"nickname" example:"约翰"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" example:"old_password123"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"new_password123"`
}

// GetUsersRequest 获取用户列表请求
type GetUsersRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1" example:"1"`
	PageSize int    `form:"page_size,default=10" binding:"min=1,max=100" example:"10"`
	Keyword  string `form:"keyword" example:"john"`
}

// GetUsersResponse 获取用户列表响应
type GetUsersResponse struct {
	Data  []*UserResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Size  int             `json:"size"`
}