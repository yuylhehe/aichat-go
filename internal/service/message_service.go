package service

import (
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
	"fmt"
	"time"
)

// MessageService 消息服务
type MessageService struct {
	repo            repository.MessageRepository
	conversationRepo repository.ConversationRepository
}

// NewMessageService 创建消息服务
func NewMessageService(repo repository.MessageRepository, conversationRepo repository.ConversationRepository) *MessageService {
	return &MessageService{
		repo:            repo,
		conversationRepo: conversationRepo,
	}
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
	count, err := s.conversationRepo.CountByIDAndUserID(req.ConversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("无权向该会话发送消息")
	}

	// 获取下一个排序值
	nextSort, err := s.repo.NextSort(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("获取消息排序失败: %w", err)
	}

	message := &model.Message{
		ConversationID:   req.ConversationID,
		Content:          req.Content,
		ReasoningContent: req.ReasoningContent,
		Sort:             nextSort,
		Type:             req.Type,
		Model:            req.Model,
		ParentID:         req.ParentID,
	}

	if err := s.repo.Create(message); err != nil {
		return nil, err
	}

	return s.toResponse(message), nil
}

// NextSort 获取下一个消息排序值
func (s *MessageService) NextSort(conversationID uint) (int, error) {
	return s.repo.NextSort(conversationID)
}

// FindByConversationID 根据会话ID查找消息
func (s *MessageService) FindByConversationID(userID, conversationID uint) ([]*MessageResponse, error) {
	// 验证会话归属权
	count, err := s.conversationRepo.CountByIDAndUserID(conversationID, userID)
	if err != nil {
		return nil, fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return nil, fmt.Errorf("无权查看该会话消息")
	}

	messages, err := s.repo.FindByConversationID(conversationID)
	if err != nil {
		return nil, err
	}

	items := make([]*MessageResponse, len(messages))
	for i, msg := range messages {
		items[i] = s.toResponse(msg)
	}

	return items, nil
}

// FindAll 获取用户的消息列表
func (s *MessageService) FindAll(userID uint) ([]*MessageResponse, error) {
	messages, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	items := make([]*MessageResponse, len(messages))
	for i, msg := range messages {
		items[i] = s.toResponse(msg)
	}

	return items, nil
}

// FindByID 根据ID查找消息
func (s *MessageService) FindByID(userID, id uint) (*MessageResponse, error) {
	message, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(message), nil
}

// Update 更新消息
func (s *MessageService) Update(userID, id uint, req *UpdateMessageRequest) (*MessageResponse, error) {
	message, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
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
		return s.toResponse(message), nil
	}

	if err := s.repo.Update(message, updates); err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	message, _ = s.repo.FindByID(id)

	return s.toResponse(message), nil
}

// Delete 删除消息
func (s *MessageService) Delete(userID, id uint) error {
	// 验证归属权
	_, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return fmt.Errorf("消息不存在或无权删除")
	}

	return s.repo.Delete(id)
}

// DeleteByConversationID 根据会话ID删除消息
func (s *MessageService) DeleteByConversationID(conversationID uint) error {
	return s.repo.DeleteByConversationID(conversationID)
}

// toResponse 转换为响应结构
func (s *MessageService) toResponse(msg *model.Message) *MessageResponse {
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
