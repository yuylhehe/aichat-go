package repository

import (
	"time"

	"gorm.io/gorm"
)

// User 用户数据库模型
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Email     string         `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	Salt      string         `json:"-" gorm:"size:255;not null"`
	Avatar    *string        `json:"avatar" gorm:"size:255"`
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	TableName string `json:"-" gorm:"tableName:user"`
}
