package service

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
	"fmt"
)

// ConversationService 会话服务
type ConversationService struct {
	repo        repository.ConversationRepository
	messageRepo repository.MessageRepository
}

// NewConversationService 创建会话服务
func NewConversationService(repo repository.ConversationRepository, messageRepo repository.MessageRepository) *ConversationService {
	return &ConversationService{
		repo:        repo,
		messageRepo: messageRepo,
	}
}

// Create 创建会话
func (s *ConversationService) Create(userID uint, req *dto.CreateConversationRequest) (*dto.ConversationResponse, error) {
	conversation := &model.Conversation{
		Name:         req.Name,
		UserID:       userID,
		SystemPrompt: req.SystemPrompt,
		Model:        req.Model,
		Temperature:  req.Temperature,
	}

	if err := s.repo.Create(conversation); err != nil {
		return nil, err
	}

	return s.toResponse(conversation, 0), nil
}

// FindByID 根据ID查找会话
func (s *ConversationService) FindByID(userID, id uint) (*dto.ConversationResponse, error) {
	conversation, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	// 获取消息数量
	messageCount, _ := s.messageRepo.CountByConversationID(id)

	return s.toResponse(conversation, messageCount), nil
}

// FindByUserID 根据用户ID查找会话
func (s *ConversationService) FindByUserID(userID uint, q string) ([]*dto.ConversationResponse, error) {
	conversations, err := s.repo.FindByUserID(userID, q)
	if err != nil {
		return nil, err
	}

	// 获取每个会话的消息数量
	items := make([]*dto.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		messageCount, _ := s.messageRepo.CountByConversationID(conv.ID)
		items[i] = s.toResponse(conv, messageCount)
	}

	return items, nil
}

// FindAll 查找所有会话
func (s *ConversationService) FindAll(q string) ([]*dto.ConversationResponse, error) {
	conversations, err := s.repo.FindAll(q)
	if err != nil {
		return nil, err
	}

	// 获取每个会话的消息数量
	items := make([]*dto.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		messageCount, _ := s.messageRepo.CountByConversationID(conv.ID)
		items[i] = s.toResponse(conv, messageCount)
	}

	return items, nil
}

// Update 更新会话
func (s *ConversationService) Update(userID, id uint, req *dto.UpdateConversationRequest) (*dto.ConversationResponse, error) {
	conversation, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	// 构建更新字段
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

	if err := s.repo.Update(conversation, updates); err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	conversation, _ = s.repo.FindByID(id)
	messageCount, _ := s.messageRepo.CountByConversationID(id)

	return s.toResponse(conversation, messageCount), nil
}

// Delete 删除会话
func (s *ConversationService) Delete(userID, id uint) error {
	// 验证归属权
	count, err := s.repo.CountByIDAndUserID(id, userID)
	if err != nil {
		return fmt.Errorf("验证会话归属权失败: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("会话不存在或无权删除")
	}

	// 先删除相关的消息
	if err := s.messageRepo.DeleteByConversationID(id); err != nil {
		return fmt.Errorf("删除会话消息失败: %w", err)
	}

	// 删除会话
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

// NextSort 获取下一个消息排序值
func (s *ConversationService) NextSort(conversationID uint) (int, error) {
	return s.messageRepo.NextSort(conversationID)
}

// toResponse 转换为响应结构
func (s *ConversationService) toResponse(conv *model.Conversation, messageCount int64) *dto.ConversationResponse {
	return &dto.ConversationResponse{
		ID:           conv.ID,
		Name:         conv.Name,
		UserID:       conv.UserID,
		IsActive:     conv.IsActive,
		SystemPrompt: conv.SystemPrompt,
		Model:        conv.Model,
		Temperature:  conv.Temperature,
		CreatedAt:    conv.CreatedAt.Format(common.TimeLayout),
		UpdatedAt:    conv.UpdatedAt.Format(common.TimeLayout),
		Messages:     messageCount,
	}
}
