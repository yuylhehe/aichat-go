package service

import (
	"ai-chat/config"
	"ai-chat/internal/repository"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	if cfg == nil {
		panic("初始化 AuthService 失败: config 不能为 nil")
	}
	return &AuthService{
		db:  db,
		cfg: cfg,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	TokenType    string    `json:"tokenType"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	User  *repository.User `json:"user"`
	Token *TokenResponse   `json:"token"`
}

// HashPassword 密码哈希
func HashPassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成盐值
func GenerateSalt() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// VerifyPassword 验证密码
func VerifyPassword(password, salt, hashedPassword string) bool {
	return HashPassword(password, salt) == hashedPassword
}

// GenerateJWT 生成JWT令牌
func (s *AuthService) GenerateJWT(user *repository.User) (*TokenResponse, error) {
	expireTime := time.Now().Add(24 * time.Hour) // 24小时过期

	claims := jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Email,
		"exp":    expireTime.Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	// 刷新令牌（更长的过期时间）
	expireTime = time.Now().Add(7 * 24 * time.Hour) // 7天过期
	claims["exp"] = expireTime.Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		ExpiresAt:    expireTime,
		TokenType:    "Bearer",
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *RegisterRequest) (*AuthResponse, error) {
	// 检查邮箱是否已存在
	var existingUser repository.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 创建新用户
	salt := GenerateSalt()
	hashedPassword := HashPassword(req.Password, salt)

	user := &repository.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Salt:     salt,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 生成令牌
	token, err := s.GenerateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*AuthResponse, error) {
	// 查找用户
	var user repository.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}

	// 验证密码
	if !VerifyPassword(req.Password, user.Salt, user.Password) {
		return nil, errors.New("邮箱或密码错误")
	}

	// 生成令牌
	token, err := s.GenerateJWT(&user)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return &AuthResponse{
		User:  &user,
		Token: token,
	}, nil
}

// ParseJWT 解析JWT令牌
func (s *AuthService) ParseJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法，防止算法替换攻击
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	return &claims, nil
}

// GetUserFromToken 从令牌中获取用户信息
func (s *AuthService) GetUserFromToken(tokenString string) (*repository.User, error) {
	claims, err := s.ParseJWT(tokenString)
	if err != nil {
		return nil, err
	}

	userID, ok := (*claims)["userId"].(float64)
	if !ok {
		return nil, errors.New("令牌中缺少用户ID")
	}

	var user repository.User
	if err := s.db.First(&user, uint(userID)).Error; err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}

	return &user, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	// 解析刷新令牌
	claims, err := s.ParseJWT(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	userID, ok := (*claims)["userId"].(float64)
	if !ok {
		return nil, errors.New("令牌中缺少用户ID")
	}

	// 查找用户
	var user repository.User
	if err := s.db.First(&user, uint(userID)).Error; err != nil {
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}

	// 生成新的访问令牌
	return s.GenerateJWT(&user)
}

// GetUserByID 根据ID获取用户信息
func (s *AuthService) GetUserByID(userID uint) (*repository.User, error) {
	var user repository.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查找用户失败: %w", err)
	}

	return &user, nil
}
