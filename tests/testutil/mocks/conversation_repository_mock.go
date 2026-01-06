package mocks

import (
	"ai-chat/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockConversationRepository 会话仓库 mock
type MockConversationRepository struct {
	mock.Mock
}

func (m *MockConversationRepository) Create(conversation *model.Conversation) error {
	args := m.Called(conversation)
	return args.Error(0)
}

func (m *MockConversationRepository) FindByID(id uint) (*model.Conversation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) FindByIDAndUserID(id, userID uint) (*model.Conversation, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) FindByUserID(userID uint, query string) ([]*model.Conversation, error) {
	args := m.Called(userID, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) FindAll(query string) ([]*model.Conversation, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Conversation), args.Error(1)
}

func (m *MockConversationRepository) Update(conversation *model.Conversation, updates map[string]interface{}) error {
	args := m.Called(conversation, updates)
	return args.Error(0)
}

func (m *MockConversationRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockConversationRepository) CountByIDAndUserID(id, userID uint) (int64, error) {
	args := m.Called(id, userID)
	return args.Get(0).(int64), args.Error(1)
}
