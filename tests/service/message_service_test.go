package service_test

import (
	"ai-chat/internal/common"
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

func TestMessageService_Create(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("成功创建消息", func(t *testing.T) {
		req := &dto.CreateMessageRequest{
			ConversationID: 1,
			Content:        "测试消息内容",
			Type:           common.RoleUser,
		}

		mockConvRepo.On("CountByIDAndUserID", uint(1), uint(1)).Return(int64(1), nil).Once()
		mockMsgRepo.On("NextSort", uint(1)).Return(1, nil).Once()
		mockMsgRepo.On("Create", mock.AnythingOfType("*model.Message")).
			Run(func(args mock.Arguments) {
				msg := args.Get(0).(*model.Message)
				msg.ID = 1
				msg.CreatedAt = time.Now()
			}).
			Return(nil).Once()

		result, err := svc.Create(1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试消息内容", result.Content)
		assert.Equal(t, common.RoleUser, result.Type)
		mockMsgRepo.AssertExpectations(t)
		mockConvRepo.AssertExpectations(t)
	})

	t.Run("无权向会话发送消息", func(t *testing.T) {
		req := &dto.CreateMessageRequest{
			ConversationID: 1,
			Content:        "测试消息内容",
			Type:           common.RoleUser,
		}

		mockConvRepo.On("CountByIDAndUserID", uint(1), uint(2)).Return(int64(0), nil).Once()

		result, err := svc.Create(2, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "无权向该会话发送消息")
		mockConvRepo.AssertExpectations(t)
	})
}

func TestMessageService_FindByConversationID(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("成功获取会话消息", func(t *testing.T) {
		messages := []*model.Message{
			{Content: "消息1", Type: common.RoleUser, Sort: 1, ConversationID: 1, CreatedAt: time.Now()},
			{Content: "消息2", Type: common.RoleAssistant, Sort: 2, ConversationID: 1, CreatedAt: time.Now()},
		}
		messages[0].ID = 1
		messages[1].ID = 2

		mockConvRepo.On("CountByIDAndUserID", uint(1), uint(1)).Return(int64(1), nil).Once()
		mockMsgRepo.On("FindByConversationID", uint(1)).Return(messages, nil).Once()

		result, err := svc.FindByConversationID(1, 1)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "消息1", result[0].Content)
		assert.Equal(t, "消息2", result[1].Content)
		mockMsgRepo.AssertExpectations(t)
		mockConvRepo.AssertExpectations(t)
	})

	t.Run("无权查看会话消息", func(t *testing.T) {
		mockConvRepo.On("CountByIDAndUserID", uint(1), uint(2)).Return(int64(0), nil).Once()

		result, err := svc.FindByConversationID(2, 1)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "无权查看该会话消息")
		mockConvRepo.AssertExpectations(t)
	})
}

func TestMessageService_FindByID(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("成功获取消息", func(t *testing.T) {
		msg := &model.Message{
			Content:        "测试消息",
			Type:           common.RoleUser,
			ConversationID: 1,
			Sort:           1,
			CreatedAt:      time.Now(),
		}
		msg.ID = 1

		mockMsgRepo.On("FindByIDAndUserID", uint(1), uint(1)).Return(msg, nil).Once()

		result, err := svc.FindByID(1, 1)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试消息", result.Content)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("消息不存在", func(t *testing.T) {
		mockMsgRepo.On("FindByIDAndUserID", uint(999), uint(1)).
			Return(nil, errors.New("消息不存在")).Once()

		result, err := svc.FindByID(1, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockMsgRepo.AssertExpectations(t)
	})
}

func TestMessageService_Update(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("成功更新消息", func(t *testing.T) {
		msg := &model.Message{
			Content:        "旧内容",
			Type:           common.RoleUser,
			ConversationID: 1,
			Sort:           1,
			CreatedAt:      time.Now(),
		}
		msg.ID = 1

		updatedMsg := &model.Message{
			Content:        "新内容",
			Type:           common.RoleUser,
			ConversationID: 1,
			Sort:           1,
			CreatedAt:      msg.CreatedAt,
		}
		updatedMsg.ID = 1

		newContent := "新内容"
		req := &dto.UpdateMessageRequest{
			Content: &newContent,
		}

		mockMsgRepo.On("FindByIDAndUserID", uint(1), uint(1)).Return(msg, nil).Once()
		mockMsgRepo.On("Update", msg, mock.AnythingOfType("map[string]interface {}")).Return(nil).Once()
		mockMsgRepo.On("FindByID", uint(1)).Return(updatedMsg, nil).Once()

		result, err := svc.Update(1, 1, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockMsgRepo.AssertExpectations(t)
	})
}

func TestMessageService_Delete(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("成功删除消息", func(t *testing.T) {
		msg := &model.Message{
			Content:        "测试消息",
			Type:           common.RoleUser,
			ConversationID: 1,
			CreatedAt:      time.Now(),
		}
		msg.ID = 1

		mockMsgRepo.On("FindByIDAndUserID", uint(1), uint(1)).Return(msg, nil).Once()
		mockMsgRepo.On("Delete", uint(1)).Return(nil).Once()

		err := svc.Delete(1, 1)

		assert.NoError(t, err)
		mockMsgRepo.AssertExpectations(t)
	})

	t.Run("消息不存在或无权删除", func(t *testing.T) {
		mockMsgRepo.On("FindByIDAndUserID", uint(999), uint(1)).
			Return(nil, errors.New("消息不存在")).Once()

		err := svc.Delete(1, 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "不存在或无权删除")
		mockMsgRepo.AssertExpectations(t)
	})
}

func TestMessageService_NextSort(t *testing.T) {
	mockMsgRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)

	svc := service.NewMessageService(mockMsgRepo, mockConvRepo)

	t.Run("获取下一个排序值", func(t *testing.T) {
		mockMsgRepo.On("NextSort", uint(1)).Return(5, nil).Once()

		result, err := svc.NextSort(1)

		assert.NoError(t, err)
		assert.Equal(t, 5, result)
		mockMsgRepo.AssertExpectations(t)
	})
}
