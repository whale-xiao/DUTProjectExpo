package model

import (
	"time"

	"gorm.io/gorm"
)

// User 管理员用户（一期仅 admin 使用，editor 预留）。
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	PasswordHash string         `gorm:"size:100;not null" json:"-"`
	Nickname     string         `gorm:"size:50" json:"nickname"`
	Role         string         `gorm:"type:enum('admin','editor');default:admin;not null" json:"role"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
