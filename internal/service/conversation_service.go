package service

import (
	"ai-chat/internal/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// ConversationService 会话服务
type ConversationService struct {
	db *gorm.DB
}

// NewConversationService 创建会话服务
func NewConversationService(db *gorm.DB) *ConversationService {
	return &ConversationService{db: db}
}

// CreateConversationRequest 创建会话请求
type CreateConversationRequest struct {
	Name         string   `json:"name" binding:"required,min=1,max=255"`
	UserID       uint     `json:"userId" binding:"required"`
	SystemPrompt *string  `json:"systemPrompt,omitempty"`
	Model        *string  `json:"model,omitempty"`
	Temperature  *float64 `json:"temperature,omitempty"`
}

// UpdateConversationRequest 更新会话请求
type UpdateConversationRequest struct {
	Name         *string  `json:"name,omitempty"`
	SystemPrompt *string  `json:"systemPrompt,omitempty"`
	Model        *string  `json:"model,omitempty"`
	Temperature  *float64 `json:"temperature,omitempty"`
	IsActive     *bool    `json:"isActive,omitempty"`
}

// ConversationResponse 会话响应
type ConversationResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	UserID       uint      `json:"userId"`
	IsActive     bool      `json:"isActive"`
	SystemPrompt *string   `json:"systemPrompt"`
	Model        *string   `json:"model"`
	Temperature  *float64  `json:"temperature"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	Messages     int64     `json:"messageCount"`
}

// Create 创建会话
func (s *ConversationService) Create(req *CreateConversationRequest) (*ConversationResponse, error) {
	conversation := &repository.Conversation{
		Name:         req.Name,
		UserID:       req.UserID,
		SystemPrompt: req.SystemPrompt,
		Model:        req.Model,
		Temperature:  req.Temperature,
	}

	if err := s.db.Create(conversation).Error; err != nil {
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}

	return s.toResponse(conversation, 0), nil
}

// FindByID 根据ID查找会话
func (s *ConversationService) FindByID(userID, id uint) (*ConversationResponse, error) {
	conversation, err := s.findByID(userID, id)
	if err != nil {
		return nil, err
	}

	// 获取消息数量
	var messageCount int64
	s.db.Model(&repository.Message{}).Where("conversation_id = ?", id).Count(&messageCount)

	return s.toResponse(conversation, messageCount), nil
}

// findByID 内部方法：根据ID查找会话
func (s *ConversationService) findByID(userID, id uint) (*repository.Conversation, error) {
	var conversation repository.Conversation
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("会话不存在")
		}
		return nil, fmt.Errorf("查找会话失败: %w", err)
	}
	return &conversation, nil
}

// FindByUserID 根据用户ID查找会话
func (s *ConversationService) FindByUserID(userID uint, q string) ([]*ConversationResponse, error) {
	var conversations []*repository.Conversation

	query := s.db.Model(&repository.Conversation{}).Where("user_id = ?", userID)

	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	// 查询列表
	if err := query.Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}

	// 获取每个会话的消息数量
	items := make([]*ConversationResponse, len(conversations))
	for i, conv := range conversations {
		var messageCount int64
		s.db.Model(&repository.Message{}).Where("conversation_id = ?", conv.ID).Count(&messageCount)
		items[i] = s.toResponse(conv, messageCount)
	}

	return items, nil
}

// FindAll 查找所有会话
func (s *ConversationService) FindAll(q string) ([]*ConversationResponse, error) {
	var conversations []*repository.Conversation

	query := s.db.Model(&repository.Conversation{})

	if q != "" {
		query = query.Where("name ILIKE ?", "%"+q+"%")
	}

	// 查询列表
	if err := query.Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}

	// 获取每个会话的消息数量
	items := make([]*ConversationResponse, len(conversations))
	for i, conv := range conversations {
		var messageCount int64
		s.db.Model(&repository.Message{}).Where("conversation_id = ?", conv.ID).Count(&messageCount)
		items[i] = s.toResponse(conv, messageCount)
	}

	return items, nil
}

// Update 更新会话
func (s *ConversationService) Update(userID, id uint, req *UpdateConversationRequest) (*ConversationResponse, error) {
	conversation, err := s.findByID(userID, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.SystemPrompt != nil {
		updates["system_prompt"] = req.SystemPrompt
	}
	if req.Model != nil {
		updates["model"] = req.Model
	}
	if req.Temperature != nil {
		updates["temperature"] = req.Temperature
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return s.toResponse(conversation, 0), nil
	}

	if err := s.db.Model(conversation).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新会话失败: %w", err)
	}

	// 重新获取更新后的数据
	s.db.First(conversation, id)

	var messageCount int64
	s.db.Model(&repository.Message{}).Where("conversation_id = ?", id).Count(&messageCount)

	return s.toResponse(conversation, messageCount), nil
}

// Delete 删除会话
func (s *ConversationService) Delete(userID, id uint) error {
	// 验证归属权
	var count int64
	if err := s.db.Model(&repository.Conversation{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
		return fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("会话不存在或无权删除")
	}

	// 先删除相关的消息
	if err := s.db.Where("conversation_id = ?", id).Delete(&repository.Message{}).Error; err != nil {
		return fmt.Errorf("删除会话消息失败: %w", err)
	}

	// 删除会话
	if err := s.db.Delete(&repository.Conversation{}, id).Error; err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}

	return nil
}

// NextSort 获取下一个消息排序值
func (s *ConversationService) NextSort(conversationID uint) (int, error) {
	var maxSort int
	err := s.db.Model(&repository.Message{}).
		Where("conversation_id = ?", conversationID).
		Pluck("COALESCE(MAX(sort), 0)", &maxSort).Error
	if err != nil {
		return 0, fmt.Errorf("查询最大排序值失败: %w", err)
	}
	return maxSort + 1, nil
}

// toResponse 转换为响应结构
func (s *ConversationService) toResponse(conv *repository.Conversation, messageCount int64) *ConversationResponse {
	return &ConversationResponse{
		ID:           conv.ID,
		Name:         conv.Name,
		UserID:       conv.UserID,
		IsActive:     conv.IsActive,
		SystemPrompt: conv.SystemPrompt,
		Model:        conv.Model,
		Temperature:  conv.Temperature,
		CreatedAt:    conv.CreatedAt,
		UpdatedAt:    conv.UpdatedAt,
		Messages:     messageCount,
	}
}
