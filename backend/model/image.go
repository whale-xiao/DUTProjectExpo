package model

import "time"

// ImageType 图片用途枚举。
type ImageType string

const (
	ImageCover      ImageType = "cover"      // 封面
	ImageScreenshot ImageType = "screenshot" // 项目截图
	ImageOther      ImageType = "other"
)

// ProjectImage 项目图片（URL 存库；V1 指向本地 /uploads，二期迁 OSS）。
// 封面同时冗余在 projects.cover_url，列表页无需 join 本表。
type ProjectImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `gorm:"not null;index" json:"projectId"`
	ImageType ImageType `gorm:"type:enum('cover','screenshot','other');not null;default:screenshot" json:"imageType"`
	URL       string    `gorm:"size:500;not null" json:"url"`
	Caption   string    `gorm:"size:255" json:"caption"`
	SortOrder int       `gorm:"not null;default:0" json:"sortOrder"`
	CreatedAt time.Time `json:"createdAt"`
}

func (ProjectImage) TableName() string { return "project_images" }
