package service

import (
	"ai-chat/internal/repository"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// MessageService 消息服务
type MessageService struct {
	db *gorm.DB
}

// NewMessageService 创建消息服务
func NewMessageService(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

// CreateMessageRequest 创建消息请求
type CreateMessageRequest struct {
	ConversationID   uint    `json:"conversationId" binding:"required"`
	Content          string  `json:"content" binding:"required"`
	ReasoningContent string  `json:"reasoningContent,omitempty"`
	Type             string  `json:"type" binding:"required,oneof=system user assistant"`
	Model            *string `json:"model,omitempty"`
	ParentID         *uint   `json:"parentId,omitempty"`
}

// UpdateMessageRequest 更新消息请求
type UpdateMessageRequest struct {
	Content *string `json:"content,omitempty"`
	Type    *string `json:"type,omitempty"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID               uint      `json:"id"`
	ConversationID   uint      `json:"conversationId"`
	Content          string    `json:"content"`
	ReasoningContent string    `json:"reasoningContent,omitempty"`
	Sort             int       `json:"sort"`
	Type             string    `json:"type"`
	Tokens           *int      `json:"tokens"`
	Model            *string   `json:"model"`
	ParentID         *uint     `json:"parentId"`
	Metadata         *string   `json:"metadata"`
	CreatedAt        time.Time `json:"createdAt"`
}

// Create 创建消息
func (s *MessageService) Create(userID uint, req *CreateMessageRequest) (*MessageResponse, error) {
	// 验证会话归属权
	var count int64
	if err := s.db.Model(&repository.Conversation{}).
		Where(&repository.Conversation{ID: req.ConversationID, UserID: userID}).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("无权向该会话发送消息")
	}

	// 获取下一个排序值
	nextSort, err := s.NextSort(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("获取消息排序失败: %w", err)
	}

	message := &repository.Message{
		ConversationID:   req.ConversationID,
		Content:          req.Content,
		ReasoningContent: req.ReasoningContent,
		Sort:             nextSort,
		Type:             req.Type,
		Model:            req.Model,
		ParentID:         req.ParentID,
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("创建消息失败: %w", err)
	}

	return s.toResponse(message), nil
}

// NextSort 获取下一个消息排序值
func (s *MessageService) NextSort(conversationID uint) (int, error) {
	var maxSort int
	err := s.db.Model(&repository.Message{}).
		Where("conversation_id = ?", conversationID).
		Pluck("COALESCE(MAX(sort), 0)", &maxSort).Error
	if err != nil {
		return 0, fmt.Errorf("查询最大排序值失败: %w", err)
	}
	return maxSort + 1, nil
}

// FindByConversationID 根据会话ID查找消息
func (s *MessageService) FindByConversationID(userID, conversationID uint) ([]*MessageResponse, error) {
	// 验证会话归属权
	var count int64
	if err := s.db.Model(&repository.Conversation{}).Where(&repository.Conversation{ID: conversationID, UserID: userID}).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("无权查看该会话消息")
	}

	var messages []*repository.Message

	query := s.db.Model(&repository.Message{}).Where("conversation_id = ?", conversationID)

	// 查询列表
	if err := query.Order("sort asc").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}

	items := make([]*MessageResponse, len(messages))
	for i, msg := range messages {
		items[i] = s.toResponse(msg)
	}

	return items, nil
}

// FindAll 获取用户的消息列表
func (s *MessageService) FindAll(userID uint) ([]*MessageResponse, error) {
	var messages []*repository.Message

	// 使用 Join 查询属于该用户的消息
	query := s.db.Model(&repository.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("conversation.user_id = ?", userID)

	// 查询列表
	if err := query.Order("message.sort asc").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}

	items := make([]*MessageResponse, len(messages))
	for i, msg := range messages {
		items[i] = s.toResponse(msg)
	}

	return items, nil
}

// FindByID 根据ID查找消息
func (s *MessageService) FindByID(userID, id uint) (*MessageResponse, error) {
	var message repository.Message

	// 使用 Join 验证归属权
	err := s.db.Model(&repository.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("message.id = ? AND conversation.user_id = ?", id, userID).
		First(&message).Error

	if err != nil {
		return nil, fmt.Errorf("查找消息失败: %w", err)
	}

	return s.toResponse(&message), nil
}

// Update 更新消息
func (s *MessageService) Update(userID, id uint, req *UpdateMessageRequest) (*MessageResponse, error) {
	var message repository.Message

	// 验证归属权并获取消息
	err := s.db.Model(&repository.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("message.id = ? AND conversation.user_id = ?", id, userID).
		First(&message).Error

	if err != nil {
		return nil, fmt.Errorf("查找消息失败: %w", err)
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if len(updates) == 0 {
		return s.toResponse(&message), nil
	}

	if err := s.db.Model(&message).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新消息失败: %w", err)
	}

	// 重新获取更新后的数据
	s.db.First(&message, id)

	return s.toResponse(&message), nil
}

// Delete 删除消息
func (s *MessageService) Delete(userID, id uint) error {
	// 验证归属权
	var count int64
	err := s.db.Model(&repository.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("message.id = ? AND conversation.user_id = ?", id, userID).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("验证消息归属权失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("消息不存在或无权删除")
	}

	if err := s.db.Delete(&repository.Message{}, id).Error; err != nil {
		return fmt.Errorf("删除消息失败: %w", err)
	}

	return nil
}

// DeleteByConversationID 根据会话ID删除消息
func (s *MessageService) DeleteByConversationID(conversationID uint) error {
	if err := s.db.Where("conversation_id = ?", conversationID).Delete(&repository.Message{}).Error; err != nil {
		return fmt.Errorf("删除会话消息失败: %w", err)
	}

	return nil
}

// toResponse 转换为响应结构
func (s *MessageService) toResponse(msg *repository.Message) *MessageResponse {
	return &MessageResponse{
		ID:               msg.ID,
		ConversationID:   msg.ConversationID,
		Content:          msg.Content,
		ReasoningContent: msg.ReasoningContent,
		Sort:             msg.Sort,
		Type:             msg.Type,
		Tokens:           msg.Tokens,
		Model:            msg.Model,
		ParentID:         msg.ParentID,
		Metadata:         msg.Metadata,
		CreatedAt:        msg.CreatedAt,
	}
}
