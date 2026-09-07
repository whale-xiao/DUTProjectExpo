package model

import "time"

// ProjectMember 项目成员。
type ProjectMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index" json:"projectId"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Role      string    `gorm:"size:100" json:"role"` // 角色/职责，如"后端开发"
	SortOrder int       `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ProjectMember) TableName() string { return "project_members" }
