package repository

import (
	"ai-chat/internal/model"
	"fmt"

	"gorm.io/gorm"
)

// FixedPromptRepository 固定提示词数据访问接口
type FixedPromptRepository interface {
	Create(fp *model.FixedPrompt) error
	FindByID(id uint) (*model.FixedPrompt, error)
	FindByIDAndUserID(id, userID uint) (*model.FixedPrompt, error)
	FindByUserID(userID uint, page, pageSize int, q string) ([]*model.FixedPrompt, int64, error)
	Update(fp *model.FixedPrompt, updates map[string]interface{}) error
	Delete(id uint) error
	DeleteByUserID(id, userID uint) error
}

type fixedPromptRepository struct {
	db *gorm.DB
}

// NewFixedPromptRepository 创建 FixedPromptRepository 实例
func NewFixedPromptRepository(db *gorm.DB) FixedPromptRepository {
	return &fixedPromptRepository{db: db}
}

func (r *fixedPromptRepository) Create(fp *model.FixedPrompt) error {
	if err := r.db.Create(fp).Error; err != nil {
		return fmt.Errorf("创建固定提示词失败: %w", err)
	}
	return nil
}

func (r *fixedPromptRepository) FindByID(id uint) (*model.FixedPrompt, error) {
	var fp model.FixedPrompt
	if err := r.db.First(&fp, id).Error; err != nil {
		return nil, fmt.Errorf("查找固定提示词失败: %w", err)
	}
	return &fp, nil
}

func (r *fixedPromptRepository) FindByIDAndUserID(id, userID uint) (*model.FixedPrompt, error) {
	var fp model.FixedPrompt
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&fp).Error; err != nil {
		return nil, fmt.Errorf("查找固定提示词失败: %w", err)
	}
	return &fp, nil
}

func (r *fixedPromptRepository) FindByUserID(userID uint, page, pageSize int, q string) ([]*model.FixedPrompt, int64, error) {
	var fps []*model.FixedPrompt
	var total int64

	query := r.db.Model(&model.FixedPrompt{}).Where("user_id = ?", userID)

	if q != "" {
		query = query.Where("name ILIKE ? OR content ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询固定提示词总数失败: %w", err)
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&fps).Error; err != nil {
		return nil, 0, fmt.Errorf("查询固定提示词列表失败: %w", err)
	}

	return fps, total, nil
}

func (r *fixedPromptRepository) Update(fp *model.FixedPrompt, updates map[string]interface{}) error {
	if err := r.db.Model(fp).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新固定提示词失败: %w", err)
	}
	return nil
}

func (r *fixedPromptRepository) Delete(id uint) error {
	if err := r.db.Delete(&model.FixedPrompt{}, id).Error; err != nil {
		return fmt.Errorf("删除固定提示词失败: %w", err)
	}
	return nil
}

func (r *fixedPromptRepository) DeleteByUserID(id, userID uint) error {
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.FixedPrompt{}).Error; err != nil {
		return fmt.Errorf("删除固定提示词失败: %w", err)
	}
	return nil
}
