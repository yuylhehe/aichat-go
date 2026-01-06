package service_test

import (
	"ai-chat/config"
	"ai-chat/internal/model"
	"ai-chat/internal/pkg/crypto"
	"ai-chat/internal/service"
	"ai-chat/tests/testutil/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func newTestConfig() *config.Config {
	return &config.Config{
		JWTSecret: "test-secret-key-for-testing-purposes",
	}
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功注册用户", func(t *testing.T) {
		req := &service.RegisterRequest{
			Name:     "测试用户",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", "newuser@example.com").
			Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("Create", mock.AnythingOfType("*model.User")).
			Run(func(args mock.Arguments) {
				user := args.Get(0).(*model.User)
				user.ID = 1
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
			}).
			Return(nil).Once()

		result, err := svc.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.NotNil(t, result.Token)
		assert.Equal(t, "测试用户", result.User.Name)
		assert.NotEmpty(t, result.Token.AccessToken)
		assert.NotEmpty(t, result.Token.RefreshToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("邮箱已被注册", func(t *testing.T) {
		existingUser := &model.User{
			Name:  "已存在用户",
			Email: "existing@example.com",
		}
		existingUser.ID = 1

		req := &service.RegisterRequest{
			Name:     "新用户",
			Email:    "existing@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", "existing@example.com").
			Return(existingUser, nil).Once()

		result, err := svc.Register(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "邮箱已被注册")
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功登录", func(t *testing.T) {
		salt := crypto.GenerateSalt()
		hashedPassword := crypto.HashPassword("password123", salt)

		user := &model.User{
			Name:      "测试用户",
			Email:     "test@example.com",
			Password:  hashedPassword,
			Salt:      salt,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user.ID = 1

		req := &service.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil).Once()

		result, err := svc.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.User)
		assert.NotNil(t, result.Token)
		assert.Equal(t, "test@example.com", result.User.Email)
		assert.NotEmpty(t, result.Token.AccessToken)
		mockRepo.AssertExpectations(t)
	})

	t.Run("邮箱不存在", func(t *testing.T) {
		req := &service.LoginRequest{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", "notfound@example.com").
			Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := svc.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "邮箱或密码错误")
		mockRepo.AssertExpectations(t)
	})

	t.Run("密码错误", func(t *testing.T) {
		salt := crypto.GenerateSalt()
		hashedPassword := crypto.HashPassword("correctpassword", salt)

		user := &model.User{
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: hashedPassword,
			Salt:     salt,
		}
		user.ID = 1

		req := &service.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil).Once()

		result, err := svc.Login(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "邮箱或密码错误")
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthService_GenerateJWT(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功生成JWT", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		result, err := svc.GenerateJWT(user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.Equal(t, "Bearer", result.TokenType)
		assert.True(t, result.ExpiresAt.After(time.Now()))
	})
}

func TestAuthService_ParseJWT(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功解析有效的JWT", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		tokenResp, _ := svc.GenerateJWT(user)

		claims, err := svc.ParseJWT(tokenResp.AccessToken)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, float64(1), (*claims)["userId"])
		assert.Equal(t, "test@example.com", (*claims)["email"])
	})

	t.Run("解析无效的JWT", func(t *testing.T) {
		claims, err := svc.ParseJWT("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功刷新令牌", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		tokenResp, _ := svc.GenerateJWT(user)

		mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()

		newToken, err := svc.RefreshToken(tokenResp.RefreshToken)

		assert.NoError(t, err)
		assert.NotNil(t, newToken)
		assert.NotEmpty(t, newToken.AccessToken)
		assert.NotEmpty(t, newToken.RefreshToken)
		assert.Equal(t, "Bearer", newToken.TokenType)
		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的刷新令牌", func(t *testing.T) {
		newToken, err := svc.RefreshToken("invalid-refresh-token")

		assert.Error(t, err)
		assert.Nil(t, newToken)
		assert.Contains(t, err.Error(), "无效的刷新令牌")
	})
}

func TestAuthService_GetUserFromToken(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	cfg := newTestConfig()
	svc := service.NewAuthService(mockRepo, cfg)

	t.Run("成功从令牌获取用户", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		tokenResp, _ := svc.GenerateJWT(user)

		mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()

		result, err := svc.GetUserFromToken(tokenResp.AccessToken)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "test@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("无效的令牌", func(t *testing.T) {
		result, err := svc.GetUserFromToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
