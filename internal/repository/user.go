package repository

import (
	"ai-chat/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Create 创建用户
	Create(user *model.User) error
	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*model.User, error)
	// FindByID 根据ID查找用户
	FindByID(id uint) (*model.User, error)
	// Update 更新用户
	Update(user *model.User) error
	// DeleteByID 删除用户
	DeleteByID(id uint) error
	// FindAll 获取所有用户
	FindAll(page, pageSize int, q string) ([]*model.User, int64, error)
	// UpdatePassword 更新用户密码
	UpdatePassword(userID uint, hashedPassword string, salt string) error
	// GetUserList 获取用户列表
	GetUserList(page, pageSize int, q string) ([]*model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) FindAll(page, pageSize int, q string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&model.User{})

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

func (r *userRepository) UpdatePassword(userID uint, hashedPassword string, salt string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password": hashedPassword,
			"salt":     salt,
		}).Error
}

// GetUserList 获取用户列表
func (r *userRepository) GetUserList(page, pageSize int, q string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&model.User{})

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
