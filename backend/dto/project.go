package dto

import "time"

// ListQuery 前台项目列表查询参数（对齐 09-前台设计方案 §4.2）。
type ListQuery struct {
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	Sort       string `form:"sort"` // recommended | newest | views
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

// PageData 分页容器：total/has_more 用于"加载更多"。
type PageData[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	HasMore  bool  `json:"has_more"`
}

// CardDTO 列表/首页卡片（一眼懂项目，字段克制）。
type CardDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Summary       string    `json:"summary"`
	CoverURL      string    `json:"coverUrl"`
	Tags          []string  `json:"tags"`
	IsRecommended bool      `json:"isRecommended"`
	CategoryID    uint      `json:"categoryId"`
	CreatedAt     time.Time `json:"createdAt"`
}

// SuggestItemDTO 联想结果（轻量，不含 blocks/members）。
type SuggestItemDTO struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Summary  string   `json:"summary"`
	CoverURL string   `json:"coverUrl"`
	Tags     []string `json:"tags"`
}

// BlockDTO 详情正文内容块（type ∈ heading|paragraph|image）。
type BlockDTO struct {
	Type    string `json:"type"`
	Text    string `json:"text,omitempty"`
	URL     string `json:"url,omitempty"`
	Caption string `json:"caption,omitempty"`
}

// MemberBriefDTO 项目成员（详情展示）。
type MemberBriefDTO struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// DetailDTO 项目详情（作品文章，09 §4.3）。
type DetailDTO struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Summary     string            `json:"summary"`
	CoverURL    string            `json:"coverUrl"`
	Category    *CategoryBriefDTO `json:"category,omitempty"`
	Tags        []string          `json:"tags"`
	Members     []MemberBriefDTO  `json:"members"`
	ViewCount   uint              `json:"viewCount"`
	PublishedAt string            `json:"publishedAt"`
	Blocks      []BlockDTO        `json:"blocks"`
}
