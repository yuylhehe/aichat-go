package model

import (
	"time"

	"gorm.io/gorm"
)

// FixedPrompt 固定提示词模型
type FixedPrompt struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:255;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"userId" gorm:"not null;index"` // Added UserID
	IsActive  bool           `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	TableName string `json:"-" gorm:"tableName:fixed_prompt"`
}

// BeforeCreate 创建前钩子
func (fp *FixedPrompt) BeforeCreate(tx *gorm.DB) error {
	fp.IsActive = true
	return nil
}
