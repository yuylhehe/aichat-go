package service

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/repository"
	"fmt"

	"gorm.io/gorm"
)

// FixedPromptService 固定提示词服务
type FixedPromptService struct {
	db *gorm.DB
}

// NewFixedPromptService 创建固定提示词服务
func NewFixedPromptService(db *gorm.DB) *FixedPromptService {
	return &FixedPromptService{db: db}
}

// Create 创建固定提示词
func (s *FixedPromptService) Create(userID uint, req *dto.CreateFixedPromptRequest) (*dto.FixedPromptResponse, error) {
	fixedPrompt := &repository.FixedPrompt{
		Name:    req.Name,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.db.Create(fixedPrompt).Error; err != nil {
		return nil, fmt.Errorf("创建固定提示词失败: %w", err)
	}

	return s.toResponse(fixedPrompt), nil
}

// FindAll 获取用户的固定提示词
func (s *FixedPromptService) FindAll(userID uint, page, pageSize int, q string) (*dto.PaginatedFixedPrompts, error) {
	var fixedPrompts []*repository.FixedPrompt
	var total int64

	query := s.db.Model(&repository.FixedPrompt{}).Where("user_id = ?", userID)

	if q != "" {
		query = query.Where("name ILIKE ? OR content ILIKE ?", "%"+q+"%", "%"+q+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询固定提示词总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&fixedPrompts).Error; err != nil {
		return nil, fmt.Errorf("查询固定提示词列表失败: %w", err)
	}

	items := make([]*dto.FixedPromptResponse, len(fixedPrompts))
	for i, fp := range fixedPrompts {
		items[i] = s.toResponse(fp)
	}

	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	return &dto.PaginatedFixedPrompts{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// FindByID 根据ID查找固定提示词
func (s *FixedPromptService) FindByID(userID, id uint) (*dto.FixedPromptResponse, error) {
	var fixedPrompt repository.FixedPrompt
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&fixedPrompt).Error; err != nil {
		return nil, fmt.Errorf("查找固定提示词失败或无权访问: %w", err)
	}

	return s.toResponse(&fixedPrompt), nil
}

// Update 更新固定提示词
func (s *FixedPromptService) Update(userID, id uint, req *dto.UpdateFixedPromptRequest) (*dto.FixedPromptResponse, error) {
	var fixedPrompt repository.FixedPrompt
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&fixedPrompt).Error; err != nil {
		return nil, fmt.Errorf("查找固定提示词失败或无权访问: %w", err)
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return s.toResponse(&fixedPrompt), nil
	}

	if err := s.db.Model(&fixedPrompt).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新固定提示词失败: %w", err)
	}

	// 重新获取更新后的数据
	s.db.First(&fixedPrompt, id)

	return s.toResponse(&fixedPrompt), nil
}

// Delete 删除固定提示词
func (s *FixedPromptService) Delete(userID, id uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&repository.FixedPrompt{}).Error; err != nil {
		return fmt.Errorf("删除固定提示词失败或无权删除: %w", err)
	}

	return nil
}

// toResponse 转换为响应结构
func (s *FixedPromptService) toResponse(fp *repository.FixedPrompt) *dto.FixedPromptResponse {
	return &dto.FixedPromptResponse{
		ID:        fp.ID,
		Name:      fp.Name,
		Content:   fp.Content,
		IsActive:  fp.IsActive,
		CreatedAt: fp.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: fp.UpdatedAt.Format(common.TimeLayout),
	}
}
