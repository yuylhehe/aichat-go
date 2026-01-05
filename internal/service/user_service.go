package service

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/model"
	"ai-chat/internal/pkg/crypto"
	"ai-chat/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create 创建用户
func (s *UserService) Create(user *model.User) error {
	return s.repo.Create(user)
}

// FindByEmail 根据邮箱查找用户
func (s *UserService) FindByEmail(email string) (*model.User, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindByID 根据ID查找用户
func (s *UserService) FindByID(id uint) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update 更新用户
func (s *UserService) Update(user *model.User) error {
	return s.repo.Update(user)
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	return s.repo.DeleteByID(id)
}

// FindAll 获取所有用户
func (s *UserService) FindAll(page, pageSize int, q string) ([]*model.User, int64, error) {
	users, total, err := s.repo.FindAll(page, pageSize, q)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetProfile 获取用户信息
func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: user.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
	}, nil
}

// UpdateProfile 更新用户信息
func (s *UserService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	// 先查询用户数据
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	// 更新用户信息
	if req.Username != "" {
		user.Name = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	// 保存更改
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: user.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
	}, nil
}

// UpdatePassword 更新用户密码
func (s *UserService) UpdatePassword(userID uint, req *dto.UpdatePasswordRequest) error {
	// 1. 查找用户
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 2. 验证旧密码
	if !crypto.VerifyPassword(req.OldPassword, user.Salt, user.Password) {
		return errors.New("旧密码错误")
	}

	// 3. 生成新盐值和哈希密码
	newSalt := crypto.GenerateSalt()
	hashedPassword := crypto.HashPassword(req.NewPassword, newSalt)

	// 4. 更新用户密码
	user.Password = hashedPassword
	user.Salt = newSalt

	if err := s.repo.UpdatePassword(userID, hashedPassword, newSalt); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeleteAccount 删除用户账户
func (s *UserService) DeleteAccount(userID uint) error {
	err := s.repo.DeleteByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}
	return err
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Name,
		CreatedAt: user.CreatedAt.Format(common.TimeLayout),
		UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
	}, nil
}
