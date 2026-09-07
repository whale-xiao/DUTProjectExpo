package model

import "time"

// VisitLog 项目访问日志（一期可选：detail GET 自动 +1 view_count；
// 二期接 Redis 计数落库时启用本表做明细）。
type VisitLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index:idx_visit_project_time,priority:1" json:"projectId"`
	IP        string    `gorm:"size:45" json:"ip"` // 兼容 IPv6
	UserAgent string    `gorm:"size:255" json:"userAgent"`
	VisitedAt time.Time `gorm:"index:idx_visit_project_time,priority:2" json:"visitedAt"`
}

func (VisitLog) TableName() string { return "visit_logs" }
