package repository

import (
	"ai-chat/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// MessageRepository 消息数据访问接口
type MessageRepository interface {
	Create(message *model.Message) error
	FindByID(id uint) (*model.Message, error)
	FindByConversationID(conversationID uint) ([]*model.Message, error)
	Update(message *model.Message, updates map[string]interface{}) error
	Delete(id uint) error
	DeleteByConversationID(conversationID uint) error
	NextSort(conversationID uint) (int, error)
	FindByIDAndUserID(id, userID uint) (*model.Message, error)
	FindByUserID(userID uint) ([]*model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository 创建 MessageRepository 实例
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *model.Message) error {
	if err := r.db.Create(message).Error; err != nil {
		return fmt.Errorf("创建消息失败: %w", err)
	}
	return nil
}

func (r *messageRepository) FindByID(id uint) (*model.Message, error) {
	var message model.Message
	if err := r.db.First(&message, id).Error; err != nil {
		return nil, fmt.Errorf("查找消息失败: %w", err)
	}
	return &message, nil
}

func (r *messageRepository) FindByConversationID(conversationID uint) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Where("conversation_id = ?", conversationID).Order("sort asc").Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) Update(message *model.Message, updates map[string]interface{}) error {
	if err := r.db.Model(message).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新消息失败: %w", err)
	}
	return nil
}

func (r *messageRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Message{}, id).Error; err != nil {
		return fmt.Errorf("删除消息失败: %w", err)
	}
	return nil
}

func (r *messageRepository) DeleteByConversationID(conversationID uint) error {
	if err := r.db.Where("conversation_id = ?", conversationID).Delete(&model.Message{}).Error; err != nil {
		return fmt.Errorf("删除会话消息失败: %w", err)
	}
	return nil
}

func (r *messageRepository) NextSort(conversationID uint) (int, error) {
	var maxSort int
	err := r.db.Model(&model.Message{}).
		Where("conversation_id = ?", conversationID).
		Pluck("COALESCE(MAX(sort), 0)", &maxSort).Error
	if err != nil {
		return 0, fmt.Errorf("查询最大排序值失败: %w", err)
	}
	return maxSort + 1, nil
}

func (r *messageRepository) FindByIDAndUserID(id, userID uint) (*model.Message, error) {
	var message model.Message
	err := r.db.Model(&model.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("message.id = ? AND conversation.user_id = ?", id, userID).
		First(&message).Error
	if err != nil {
		return nil, fmt.Errorf("查找消息失败: %w", err)
	}
	return &message, nil
}

func (r *messageRepository) FindByUserID(userID uint) ([]*model.Message, error) {
	var messages []*model.Message
	err := r.db.Model(&model.Message{}).
		Joins("JOIN conversation ON conversation.id = message.conversation_id").
		Where("conversation.user_id = ?", userID).
		Order("message.sort asc").Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}
	return messages, nil
}
