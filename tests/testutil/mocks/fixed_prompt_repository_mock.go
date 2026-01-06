package mocks

import (
	"ai-chat/internal/model"

	"github.com/stretchr/testify/mock"
)

// MockFixedPromptRepository 固定提示词仓库 mock
type MockFixedPromptRepository struct {
	mock.Mock
}

func (m *MockFixedPromptRepository) Create(fp *model.FixedPrompt) error {
	args := m.Called(fp)
	return args.Error(0)
}

func (m *MockFixedPromptRepository) FindByID(id uint) (*model.FixedPrompt, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FixedPrompt), args.Error(1)
}

func (m *MockFixedPromptRepository) FindByIDAndUserID(id, userID uint) (*model.FixedPrompt, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.FixedPrompt), args.Error(1)
}

func (m *MockFixedPromptRepository) FindByUserID(userID uint, page, pageSize int, q string) ([]*model.FixedPrompt, int64, error) {
	args := m.Called(userID, page, pageSize, q)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*model.FixedPrompt), args.Get(1).(int64), args.Error(2)
}

func (m *MockFixedPromptRepository) Update(fp *model.FixedPrompt, updates map[string]interface{}) error {
	args := m.Called(fp, updates)
	return args.Error(0)
}

func (m *MockFixedPromptRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockFixedPromptRepository) DeleteByUserID(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}
