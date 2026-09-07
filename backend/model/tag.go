package model

import (
	"time"

	"gorm.io/gorm"
)

// Tag 技术栈/特性标签。
type Tag struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"uniqueIndex;size:30;not null" json:"name"`
	Color     string         `gorm:"size:20" json:"color"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProjectTag 项目-标签多对多关联（联合主键；删除项目/标签时级联清理）。
type ProjectTag struct {
	ProjectID uint `gorm:"primaryKey" json:"projectId"`
	TagID     uint `gorm:"primaryKey" json:"tagId"`
}
