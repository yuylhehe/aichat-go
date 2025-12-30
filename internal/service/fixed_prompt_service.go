package service

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/model"
	"ai-chat/internal/repository"
)

// FixedPromptService 固定提示词服务
type FixedPromptService struct {
	repo repository.FixedPromptRepository
}

// NewFixedPromptService 创建固定提示词服务
func NewFixedPromptService(repo repository.FixedPromptRepository) *FixedPromptService {
	return &FixedPromptService{repo: repo}
}

// Create 创建固定提示词
func (s *FixedPromptService) Create(userID uint, req *dto.CreateFixedPromptRequest) (*dto.FixedPromptResponse, error) {
	fixedPrompt := &model.FixedPrompt{
		Name:    req.Name,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(fixedPrompt); err != nil {
		return nil, err
	}

	return s.toResponse(fixedPrompt), nil
}

// FindAll 获取用户的固定提示词
func (s *FixedPromptService) FindAll(userID uint, page, pageSize int, q string) (*dto.PaginatedFixedPrompts, error) {
	fixedPrompts, total, err := s.repo.FindByUserID(userID, page, pageSize, q)
	if err != nil {
		return nil, err
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
	fixedPrompt, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
	}

	return s.toResponse(fixedPrompt), nil
}

// Update 更新固定提示词
func (s *FixedPromptService) Update(userID, id uint, req *dto.UpdateFixedPromptRequest) (*dto.FixedPromptResponse, error) {
	fixedPrompt, err := s.repo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, err
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
		return s.toResponse(fixedPrompt), nil
	}

	if err := s.repo.Update(fixedPrompt, updates); err != nil {
		return nil, err
	}

	// 重新获取更新后的数据
	fixedPrompt, _ = s.repo.FindByID(id)

	return s.toResponse(fixedPrompt), nil
}

// Delete 删除固定提示词
func (s *FixedPromptService) Delete(userID, id uint) error {
	return s.repo.DeleteByUserID(id, userID)
}

// toResponse 转换为响应结构
func (s *FixedPromptService) toResponse(fp *model.FixedPrompt) *dto.FixedPromptResponse {
	return &dto.FixedPromptResponse{
		ID:        fp.ID,
		Name:      fp.Name,
		Content:   fp.Content,
		IsActive:  fp.IsActive,
		CreatedAt: fp.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: fp.UpdatedAt.Format(common.TimeLayout),
	}
}
