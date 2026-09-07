package model

import (
	"time"

	"gorm.io/gorm"
)

// 项目状态枚举（前台只消费 Published；后台管理三者）。
const (
	ProjectStatusDraft     = "draft"
	ProjectStatusPublished = "published"
	ProjectStatusArchived  = "archived"
)

// Project 项目主表（详情正文由 ProjectContentBlock 承载）。
type Project struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Summary        string         `gorm:"size:255;not null" json:"summary"` // 一句话介绍（卡片/搜索）
	Description    string         `gorm:"type:text" json:"description"`     // 简介补充（兼容过渡）
	CategoryID     uint           `gorm:"not null;index" json:"categoryId"` // 所属分类
	CoverURL       string         `gorm:"size:500" json:"coverUrl"`
	Status         string         `gorm:"type:enum('draft','published','archived');default:draft;not null;index" json:"status"`
	IsRecommended  bool           `gorm:"not null;default:false;index" json:"isRecommended"`
	RecommendOrder int            `gorm:"not null;default:0" json:"recommendOrder"`
	ViewCount      uint           `gorm:"not null;default:0" json:"viewCount"`
	CreatedBy      uint           `gorm:"not null" json:"createdBy"` // 创建人 user.id
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
