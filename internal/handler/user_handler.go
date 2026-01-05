package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile 获取用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		SendInternalError(c, "获取用户信息失败")
		return
	}

	SendSuccess(c, user)
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	response, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		SendInternalError(c, "更新用户信息失败")
		return
	}

	SendSuccess(c, response)
}

// UpdatePassword 更新用户密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		SendBadRequest(c, "更新密码失败: "+err.Error())
		return
	}

	SendSuccess(c, "密码更新成功")
}
