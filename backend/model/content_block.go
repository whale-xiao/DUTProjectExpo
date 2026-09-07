package model

import "time"

// BlockType 内容块类型枚举（前台渲染由该字段分发）。
type BlockType string

const (
	BlockHeading   BlockType = "heading"   // 小节标题
	BlockParagraph BlockType = "paragraph" // 段落文本
	BlockImage     BlockType = "image"     // 图片（可带说明）
)

// ProjectContentBlock 项目详情正文的内容块（★核心表）。
// 有序排列（SortOrder），保存时整组事务内全量替换，见 10-调用链 §5.1。
type ProjectContentBlock struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ProjectID    uint      `gorm:"not null;index:idx_project_sort,priority:1" json:"projectId"`
	BlockType    BlockType `gorm:"type:enum('heading','paragraph','image');not null" json:"type"`
	TextContent  string    `gorm:"type:text" json:"text,omitempty"` // heading/paragraph 使用
	ImageURL     string    `gorm:"size:500" json:"url,omitempty"`   // image 使用
	ImageCaption string    `gorm:"size:255" json:"caption,omitempty"`
	SortOrder    int       `gorm:"not null;default:0;index:idx_project_sort,priority:2" json:"sortOrder"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (ProjectContentBlock) TableName() string { return "project_content_blocks" }
