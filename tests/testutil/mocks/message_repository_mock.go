package mocks

import (
	"ai-chat/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockMessageRepository 消息仓库 mock
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(message *model.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockMessageRepository) FindByID(id uint) (*model.Message, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageRepository) FindByConversationID(conversationID uint) ([]*model.Message, error) {
	args := m.Called(conversationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func (m *MockMessageRepository) Update(message *model.Message, updates map[string]interface{}) error {
	args := m.Called(message, updates)
	return args.Error(0)
}

func (m *MockMessageRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMessageRepository) DeleteByConversationID(conversationID uint) error {
	args := m.Called(conversationID)
	return args.Error(0)
}

func (m *MockMessageRepository) NextSort(conversationID uint) (int, error) {
	args := m.Called(conversationID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockMessageRepository) FindByIDAndUserID(id, userID uint) (*model.Message, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Message), args.Error(1)
}

func (m *MockMessageRepository) FindByUserID(userID uint) ([]*model.Message, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Message), args.Error(1)
}

func (m *MockMessageRepository) CountByConversationID(conversationID uint) (int64, error) {
	args := m.Called(conversationID)
	return args.Get(0).(int64), args.Error(1)
}
