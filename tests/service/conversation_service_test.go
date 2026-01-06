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
)

func TestConversationService_Create(t *testing.T) {
	mockConvRepo := new(mocks.MockConversationRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)

	svc := service.NewConversationService(mockConvRepo, mockMsgRepo)

	t.Run("成功创建会话", func(t *testing.T) {
		req := &dto.CreateConversationRequest{
			Name:   "测试会话",
			UserID: 1,
		}

		mockConvRepo.On("Create", mock.AnythingOfType("*model.Conversation")).
			Run(func(args mock.Arguments) {
				conv := args.Get(0).(*model.Conversation)
				conv.ID = 1
				conv.CreatedAt = time.Now()
				conv.UpdatedAt = time.Now()
			}).
			Return(nil).Once()

		result, err := svc.Create(req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试会话", result.Name)
		assert.Equal(t, uint(1), result.UserID)
		mockConvRepo.AssertExpectations(t)
	})

	t.Run("创建失败-数据库错误", func(t *testing.T) {
		req := &dto.CreateConversationRequest{
			Name:   "测试会话",
			UserID: 1,
		}

		mockConvRepo.On("Create", mock.AnythingOfType("*model.Conversation")).
			Return(errors.New("数据库错误")).Once()

		result, err := svc.Create(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockConvRepo.AssertExpectations(t)
	})
}

func TestConversationService_FindByID(t *testing.T) {
	mockConvRepo := new(mocks.MockConversationRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)

	svc := service.NewConversationService(mockConvRepo, mockMsgRepo)

	t.Run("成功查找会话", func(t *testing.T) {
		conv := &model.Conversation{
			Name:      "测试会话",
			UserID:    1,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		conv.ID = 1

		mockConvRepo.On("FindByIDAndUserID", uint(1), uint(1)).Return(conv, nil).Once()
		mockMsgRepo.On("CountByConversationID", uint(1)).Return(int64(5), nil).Once()

		result, err := svc.FindByID(1, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试会话", result.Name)
		assert.Equal(t, int64(5), result.Messages)
		mockConvRepo.AssertExpectations(t)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("会话不存在", func(t *testing.T) {
		mockConvRepo.On("FindByIDAndUserID", uint(999), uint(1)).
			Return(nil, errors.New("会话不存在")).Once()

		result, err := svc.FindByID(1, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockConvRepo.AssertExpectations(t)
	})
}

func TestConversationService_FindByUserID(t *testing.T) {
	mockConvRepo := new(mocks.MockConversationRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)

	svc := service.NewConversationService(mockConvRepo, mockMsgRepo)

	t.Run("成功查找用户的会话列表", func(t *testing.T) {
		convs := []*model.Conversation{
			{Name: "会话1", UserID: 1, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{Name: "会话2", UserID: 1, IsActive: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		convs[0].ID = 1
		convs[1].ID = 2

		mockConvRepo.On("FindByUserID", uint(1), "").Return(convs, nil).Once()
		mockMsgRepo.On("CountByConversationID", uint(1)).Return(int64(3), nil).Once()
		mockMsgRepo.On("CountByConversationID", uint(2)).Return(int64(5), nil).Once()

		result, err := svc.FindByUserID(1, "")

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "会话1", result[0].Name)
		assert.Equal(t, "会话2", result[1].Name)
		mockConvRepo.AssertExpectations(t)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("用户没有会话", func(t *testing.T) {
		mockConvRepo.On("FindByUserID", uint(2), "").
			Return([]*model.Conversation{}, nil).Once()

		result, err := svc.FindByUserID(2, "")

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		mockConvRepo.AssertExpectations(t)
	})
}

func TestConversationService_Update(t *testing.T) {
	mockConvRepo := new(mocks.MockConversationRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)

	svc := service.NewConversationService(mockConvRepo, mockMsgRepo)

	t.Run("成功更新会话名称", func(t *testing.T) {
		conv := &model.Conversation{
			Name:      "旧名称",
			UserID:    1,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		conv.ID = 1

		updatedConv := &model.Conversation{
			Name:      "新名称",
			UserID:    1,
			IsActive:  true,
			CreatedAt: conv.CreatedAt,
			UpdatedAt: time.Now(),
		}
		updatedConv.ID = 1

		newName := "新名称"
		req := &dto.UpdateConversationRequest{
			Name: &newName,
		}

		mockConvRepo.On("FindByIDAndUserID", uint(1), uint(1)).Return(conv, nil).Once()
		mockConvRepo.On("Update", conv, mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()
		mockConvRepo.On("FindByID", uint(1)).Return(updatedConv, nil).Once()
		mockMsgRepo.On("CountByConversationID", uint(1)).Return(int64(0), nil).Once()

		result, err := svc.Update(1, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockConvRepo.AssertExpectations(t)
	})
}

func TestConversationService_Delete(t *testing.T) {
	mockConvRepo := new(mocks.MockConversationRepository)
	mockMsgRepo := new(mocks.MockMessageRepository)

	svc := service.NewConversationService(mockConvRepo, mockMsgRepo)

	t.Run("成功删除会话", func(t *testing.T) {
		mockConvRepo.On("CountByIDAndUserID", uint(1), uint(1)).Return(int64(1), nil).Once()
		mockMsgRepo.On("DeleteByConversationID", uint(1)).Return(nil).Once()
		mockConvRepo.On("Delete", uint(1)).Return(nil).Once()

		err := svc.Delete(1, 1)

		assert.NoError(t, err)
		mockConvRepo.AssertExpectations(t)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("会话不存在或无权删除", func(t *testing.T) {
		mockConvRepo.On("CountByIDAndUserID", uint(999), uint(1)).Return(int64(0), nil).Once()

		err := svc.Delete(1, 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在或无权删除")
		mockConvRepo.AssertExpectations(t)
	})
}
