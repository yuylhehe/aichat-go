package handler

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	serviceReq := &service.RegisterRequest{
		Name:     req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.authService.Register(serviceReq)
	if err != nil {
		SendInternalError(c, "注册失败: "+err.Error())
		return
	}

	SendCreated(c, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			Username:  result.User.Name,
			CreatedAt: result.User.CreatedAt.Format(common.TimeLayout),
			UpdatedAt: result.User.UpdatedAt.Format(common.TimeLayout),
		},
		AccessToken:  result.Token.AccessToken,
		RefreshToken: result.Token.RefreshToken,
		ExpiresAt:    result.Token.ExpiresAt.Unix(),
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	serviceReq := &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.authService.Login(serviceReq)
	if err != nil {
		SendUnauthorized(c, "邮箱或密码错误")
		return
	}

	SendSuccess(c, dto.AuthResponse{
		User: dto.UserResponse{
			ID:        result.User.ID,
			Email:     result.User.Email,
			Username:  result.User.Name,
			CreatedAt: result.User.CreatedAt.Format(common.TimeLayout),
			UpdatedAt: result.User.UpdatedAt.Format(common.TimeLayout),
		},
		AccessToken:  result.Token.AccessToken,
		RefreshToken: result.Token.RefreshToken,
		ExpiresAt:    result.Token.ExpiresAt.Unix(),
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		SendBadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		SendUnauthorized(c, "刷新令牌失败: "+err.Error())
		return
	}

	SendSuccess(c, gin.H{
		"accessToken":  result.AccessToken,
		"refreshToken": result.RefreshToken,
		"expiresAt":    result.ExpiresAt.Unix(),
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		SendNotFound(c, "用户不存在")
		return
	}

	SendSuccess(c, dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: user.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	SendSuccess(c, "登出成功")
}
