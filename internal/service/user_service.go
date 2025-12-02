package service

import (
	"ai-chat/internal/common"
	"ai-chat/internal/dto"
	"ai-chat/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// userService 用户服务结构体（私有）
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB) *userService {
	return &userService{db: db}
}

// UserService 用户服务接口（在service包中定义）
type UserService interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error)
	UpdatePassword(userID uint, req *dto.UpdatePasswordRequest) error
	DeleteAccount(userID uint) error
	GetUserList(req *dto.GetUsersRequest) ([]*dto.UserResponse, int64, error)
	GetUserByID(id uint) (*dto.UserResponse, error)
}

// UserRepository 用户仓库接口
type UserRepository interface {
	Create(user *repository.User) error
	FindByEmail(email string) (*repository.User, error)
	FindByID(id uint) (*repository.User, error)
	Update(user *repository.User) error
	Delete(id uint) error
	FindAll(page, pageSize int, q string) ([]*repository.User, int64, error)
}

// userRepository 用户仓库实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *repository.User) error {
	return r.db.Create(user).Error
}

// FindByEmail 根据邮箱查找用户
func (r *userRepository) FindByEmail(email string) (*repository.User, error) {
	var user repository.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *userRepository) FindByID(id uint) (*repository.User, error) {
	var user repository.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *repository.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&repository.User{}, id).Error
}

// FindAll 获取所有用户
func (r *userRepository) FindAll(page, pageSize int, q string) ([]*repository.User, int64, error) {
	var users []*repository.User
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&repository.User{})

	if q != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+q+"%", "%"+q+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// 实现UserService接口方法

// GetProfile 获取用户信息
func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
	var user repository.User
	err := s.db.First(&user, userID).Error
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
func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	// 先查询用户数据
	var user repository.User
	err := s.db.First(&user, userID).Error
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
	if err := s.db.Save(&user).Error; err != nil {
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
func (s *userService) UpdatePassword(userID uint, req *dto.UpdatePasswordRequest) error {
	// 1. 查找用户
	var user repository.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 2. 验证旧密码
	if !VerifyPassword(req.OldPassword, user.Salt, user.Password) {
		return errors.New("旧密码错误")
	}

	// 3. 生成新盐值和哈希密码
	newSalt := GenerateSalt()
	hashedPassword := HashPassword(req.NewPassword, newSalt)

	// 4. 更新用户密码
	user.Password = hashedPassword
	user.Salt = newSalt

	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeleteAccount 删除用户账户
func (s *userService) DeleteAccount(userID uint) error {
	var user repository.User
	err := s.db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	return s.db.Delete(&user).Error
}

// GetUserList 获取用户列表
func (s *userService) GetUserList(req *dto.GetUsersRequest) ([]*dto.UserResponse, int64, error) {
	var users []*repository.User
	var total int64

	offset := (req.Page - 1) * req.PageSize

	query := s.db.Model(&repository.User{})

	if req.Keyword != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset(offset).Limit(req.PageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, &dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Name,
			CreatedAt: user.CreatedAt.Format(common.TimeLayout),
			UpdatedAt: user.UpdatedAt.Format(common.TimeLayout),
		})
	}

	return responses, total, nil
}

// GetUserByID 根据ID获取用户信息
func (s *userService) GetUserByID(id uint) (*dto.UserResponse, error) {
	var user repository.User
	err := s.db.First(&user, id).Error
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
