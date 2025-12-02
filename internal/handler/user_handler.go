package handler

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/repository"
	"ai-chat/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// userToUserResponse 将repository.User转换为dto.UserResponse
func userToUserResponse(user *repository.User) *dto.UserResponse {
	if user == nil {
		return nil
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: user.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
	}
}

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// UserService 用户服务接口（在service包中定义）

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
func NewUserHandler(userService service.UserService) *UserHandler {
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

// DeleteAccount 删除用户账户
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	err := h.userService.DeleteAccount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"error": "删除账户失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "账户删除成功",
	})
}

// GetUserList 获取用户列表（管理员功能）
func (h *UserHandler) GetUserList(c *gin.Context) {
	var req dto.GetUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Code: 400, Error: "请求参数错误"})
		return
	}

	data, total, err := h.userService.GetUserList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Code: 500, Error: "获取用户列表失败"})
		return
	}

	response := dto.GetUsersResponse{
		Data:  data,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}

	c.JSON(http.StatusOK, SuccessResponse{Code: 0, Data: response})

	// 添加分页信息到响应头
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Page", strconv.Itoa(req.Page))
	c.Header("X-Page-Size", strconv.Itoa(req.PageSize))
}

// GetUserByID 根据ID获取用户信息（管理员功能）
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"error": "无效的ID格式",
		})
		return
	}

	// TODO: 检查管理员权限

	response, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  404,
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}
