package handler

import (
	"ai-chat/internal/common"
	"ai-chat/internal/middleware"
	"ai-chat/internal/service"
	"net/http"

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

// RegisterRequest 注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=30"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresAt    int64        `json:"expiresAt"`
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	serviceReq := &service.RegisterRequest{
		Name:     req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.authService.Register(serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "注册失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": AuthResponse{
			User: UserResponse{
				ID:        result.User.ID,
				Email:     result.User.Email,
				Username:  result.User.Name,
				CreatedAt: result.User.CreatedAt.Format(common.TimeLayout),
				UpdatedAt: result.User.UpdatedAt.Format(common.TimeLayout),
			},
			AccessToken:  result.Token.AccessToken,
			RefreshToken: result.Token.RefreshToken,
			ExpiresAt:    result.Token.ExpiresAt.Unix(),
		},
	})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	serviceReq := &service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	result, err := h.authService.Login(serviceReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "登录失败",
			"details": "邮箱或密码错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": AuthResponse{
			User: UserResponse{
				ID:        result.User.ID,
				Email:     result.User.Email,
				Username:  result.User.Name,
				CreatedAt: result.User.CreatedAt.Format(common.TimeLayout),
				UpdatedAt: result.User.UpdatedAt.Format(common.TimeLayout),
			},
			AccessToken:  result.Token.AccessToken,
			RefreshToken: result.Token.RefreshToken,
			ExpiresAt:    result.Token.ExpiresAt.Unix(),
		},
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求参数错误",
			"details": err.Error(),
		})
		return
	}

	result, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "刷新令牌失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"accessToken":  result.AccessToken,
			"refreshToken": result.RefreshToken,
			"expiresAt":    result.ExpiresAt,
		},
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Name,
			CreatedAt: user.CreatedAt.Format(common.TimeLayout),
			UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}
