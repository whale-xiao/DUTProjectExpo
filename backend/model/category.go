package model

import (
	"time"

	"gorm.io/gorm"
)

// Category 项目分类（新版模型，04 已确认采用）。
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	SortOrder   int            `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt   time.Time      `json:"createdAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
