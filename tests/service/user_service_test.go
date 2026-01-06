package service_test

import (
	"ai-chat/internal/dto"
	"ai-chat/internal/model"
	"ai-chat/internal/service"
	"ai-chat/tests/testutil/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserService_Create(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功创建用户", func(t *testing.T) {
		user := &model.User{
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}

		mockRepo.On("Create", user).Return(nil).Once()

		err := svc.Create(user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("创建失败-数据库错误", func(t *testing.T) {
		user := &model.User{
			Name:     "测试用户",
			Email:    "test@example.com",
			Password: "hashedPassword",
			Salt:     "salt",
		}

		mockRepo.On("Create", user).Return(errors.New("数据库错误")).Once()

		err := svc.Create(user)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_FindByEmail(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功查找用户", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		mockRepo.On("FindByEmail", "test@example.com").Return(user, nil).Once()

		result, err := svc.FindByEmail("test@example.com")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test@example.com", result.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo.On("FindByEmail", "notfound@example.com").
			Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := svc.FindByEmail("notfound@example.com")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_FindByID(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功查找用户", func(t *testing.T) {
		user := &model.User{
			Name:  "测试用户",
			Email: "test@example.com",
		}
		user.ID = 1

		mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()

		result, err := svc.FindByID(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo.On("FindByID", uint(999)).
			Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := svc.FindByID(999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Update(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功更新用户", func(t *testing.T) {
		user := &model.User{
			Name:  "更新后的用户",
			Email: "test@example.com",
		}
		user.ID = 1

		mockRepo.On("Update", user).Return(nil).Once()

		err := svc.Update(user)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Delete(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功删除用户", func(t *testing.T) {
		mockRepo.On("DeleteByID", uint(1)).Return(nil).Once()

		err := svc.Delete(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetProfile(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功获取用户信息", func(t *testing.T) {
		user := &model.User{
			Name:      "测试用户",
			Email:     "test@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user.ID = 1

		mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()

		result, err := svc.GetProfile(1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "test@example.com", result.Email)
		assert.Equal(t, "测试用户", result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("用户不存在", func(t *testing.T) {
		mockRepo.On("FindByID", uint(999)).
			Return(nil, gorm.ErrRecordNotFound).Once()

		result, err := svc.GetProfile(999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "用户不存在")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateProfile(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功更新用户信息", func(t *testing.T) {
		user := &model.User{
			Name:      "旧名称",
			Email:     "old@example.com",
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user.ID = 1

		req := &dto.UpdateProfileRequest{
			Username: "新名称",
			Email:    "new@example.com",
		}

		mockRepo.On("FindByID", uint(1)).Return(user, nil).Once()
		mockRepo.On("Update", mock.AnythingOfType("*model.User")).Return(nil).Once()

		result, err := svc.UpdateProfile(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_FindAll(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	svc := service.NewUserService(mockRepo)

	t.Run("成功获取用户列表", func(t *testing.T) {
		users := []*model.User{
			{Name: "用户1", Email: "user1@example.com"},
			{Name: "用户2", Email: "user2@example.com"},
		}
		users[0].ID = 1
		users[1].ID = 2

		mockRepo.On("FindAll", 1, 10, "").Return(users, int64(2), nil).Once()

		result, total, err := svc.FindAll(1, 10, "")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("带搜索条件的用户列表", func(t *testing.T) {
		users := []*model.User{
			{Name: "测试用户", Email: "test@example.com"},
		}
		users[0].ID = 1

		mockRepo.On("FindAll", 1, 10, "测试").Return(users, int64(1), nil).Once()

		result, total, err := svc.FindAll(1, 10, "测试")

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), total)
		mockRepo.AssertExpectations(t)
	})
}
