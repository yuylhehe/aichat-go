package handler

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// 错误响应结构
type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

// 成功响应结构
type SuccessResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Code: 500, Error: "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Code: 0, Data: user})
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: 400, Error: "请求参数错误"})
		return
	}

	response, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Code: 500, Error: "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Code: 0, Data: response})
}

// UpdatePassword 更新用户密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	err := h.userService.UpdatePassword(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "更新密码失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码更新成功",
	})
}
