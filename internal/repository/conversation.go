package repository

import (
	"ai-chat/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// ConversationRepository 会话数据访问接口
// 定义所有会话相关的数据库操作方法
type ConversationRepository interface {
	// Create 创建会话
	Create(conversation *model.Conversation) error
	// FindByID 根据ID查找会话
	FindByID(id uint) (*model.Conversation, error)
	// FindByIDAndUserID 根据ID和用户ID查找会话（验证归属权）
	FindByIDAndUserID(id, userID uint) (*model.Conversation, error)
	// FindByUserID 根据用户ID查找会话列表
	FindByUserID(userID uint, query string) ([]*model.Conversation, error)
	// FindAll 查找所有会话
	FindAll(query string) ([]*model.Conversation, error)
	// Update 更新会话
	Update(conversation *model.Conversation, updates map[string]interface{}) error
	// Delete 删除会话
	Delete(id uint) error
	// CountByIDAndUserID 统计指定ID和用户ID的会话数量
	CountByIDAndUserID(id, userID uint) (int64, error)
}

// conversationRepository ConversationRepository 接口的具体实现
type conversationRepository struct {
	db *gorm.DB
}

// NewConversationRepository 创建 ConversationRepository 实例
// 这是一个工厂函数，用于依赖注入
func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

// Create 创建会话
func (r *conversationRepository) Create(conversation *model.Conversation) error {
	if err := r.db.Create(conversation).Error; err != nil {
		return fmt.Errorf("创建会话失败: %w", err)
	}
	return nil
}

// FindByID 根据ID查找会话
func (r *conversationRepository) FindByID(id uint) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := r.db.First(&conversation, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("会话不存在")
		}
		return nil, fmt.Errorf("查找会话失败: %w", err)
	}
	return &conversation, nil
}

// FindByIDAndUserID 根据ID和用户ID查找会话（验证归属权）
func (r *conversationRepository) FindByIDAndUserID(id, userID uint) (*model.Conversation, error) {
	var conversation model.Conversation
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&conversation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("会话不存在")
		}
		return nil, fmt.Errorf("查找会话失败: %w", err)
	}
	return &conversation, nil
}

// FindByUserID 根据用户ID查找会话列表
func (r *conversationRepository) FindByUserID(userID uint, query string) ([]*model.Conversation, error) {
	var conversations []*model.Conversation

	q := r.db.Model(&model.Conversation{}).Where("user_id = ?", userID)

	if query != "" {
		q = q.Where("name ILIKE ?", "%"+query+"%")
	}

	if err := q.Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}

	return conversations, nil
}

// FindAll 查找所有会话
func (r *conversationRepository) FindAll(query string) ([]*model.Conversation, error) {
	var conversations []*model.Conversation

	q := r.db.Model(&model.Conversation{})

	if query != "" {
		q = q.Where("name ILIKE ?", "%"+query+"%")
	}

	if err := q.Order("created_at DESC").Find(&conversations).Error; err != nil {
		return nil, fmt.Errorf("查询会话列表失败: %w", err)
	}

	return conversations, nil
}

// Update 更新会话
func (r *conversationRepository) Update(conversation *model.Conversation, updates map[string]interface{}) error {
	if err := r.db.Model(conversation).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新会话失败: %w", err)
	}
	return nil
}

// Delete 删除会话
func (r *conversationRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.Conversation{}, id).Error; err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}
	return nil
}

// CountByIDAndUserID 统计指定ID和用户ID的会话数量
func (r *conversationRepository) CountByIDAndUserID(id, userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Conversation{}).Where("id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("统计会话数量失败: %w", err)
	}
	return count, nil
}
