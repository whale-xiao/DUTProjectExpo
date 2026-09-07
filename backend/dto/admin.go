package dto

import "time"

// AdminProjectQuery 后台项目列表查询参数（不过滤状态，含草稿/下架）。
type AdminProjectQuery struct {
	Keyword    string `form:"keyword"`
	CategoryID uint   `form:"category_id"`
	Status     string `form:"status"` // draft | published | archived
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
}

// AdminProjectItem 后台项目列表行。
type AdminProjectItem struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Summary        string    `json:"summary"`
	CoverURL       string    `json:"coverUrl"`
	CategoryID     uint      `json:"categoryId"`
	CategoryName   string    `json:"categoryName,omitempty"`
	Status         string    `json:"status"`
	IsRecommended  bool      `json:"isRecommended"`
	RecommendOrder int       `json:"recommendOrder"`
	ViewCount      uint      `json:"viewCount"`
	Tags           []string  `json:"tags"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
}

// MemberReq 成员编辑项。
type MemberReq struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// BlockReq 内容块编辑项（type ∈ heading|paragraph|image）。
type BlockReq struct {
	Type    string `json:"type" binding:"required,oneof=heading paragraph image"`
	Text    string `json:"text"`
	URL     string `json:"url"`
	Caption string `json:"caption"`
}

// SaveProjectReq 后台整单保存（新建/编辑共用）：基础信息 + 标签 + 成员 + 内容块。
// 后端在一个事务内全量替换 tags/members/blocks（对齐 10-调用链 §5.1）。
type SaveProjectReq struct {
	Name           string      `json:"name" binding:"required,min=1,max=100"`
	Summary        string      `json:"summary" binding:"required,max=255"`
	CategoryID     uint        `json:"categoryId" binding:"required,gt=0"`
	CoverURL       string      `json:"coverUrl"`
	Status         string      `json:"status" binding:"omitempty,oneof=draft published archived"`
	IsRecommended  bool        `json:"isRecommended"`
	RecommendOrder int         `json:"recommendOrder"`
	Tags           []string    `json:"tags"`
	Members        []MemberReq `json:"members"`
	Blocks         []BlockReq  `json:"blocks"`
}

// SetStatusReq 上下架请求。
type SetStatusReq struct {
	Status string `json:"status" binding:"required,oneof=draft published archived"`
}

// RecommendReq 推荐位请求。
type RecommendReq struct {
	IsRecommended  bool `json:"isRecommended"`
	RecommendOrder int  `json:"recommendOrder"`
}

// CategorySaveReq 分类新建/编辑。
type CategorySaveReq struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}

// TagSaveReq 标签新建/编辑。
type TagSaveReq struct {
	Name  string `json:"name" binding:"required,max=30"`
	Color string `json:"color"`
}

// AdminDetailDTO 后台编辑回显详情（任意状态可见）。
type AdminDetailDTO struct {
	ID             uint              `json:"id"`
	Name           string            `json:"name"`
	Summary        string            `json:"summary"`
	CoverURL       string            `json:"coverUrl"`
	Category       *CategoryBriefDTO `json:"category,omitempty"`
	Status         string            `json:"status"`
	IsRecommended  bool              `json:"isRecommended"`
	RecommendOrder int               `json:"recommendOrder"`
	ViewCount      uint              `json:"viewCount"`
	Tags           []string          `json:"tags"`
	Members        []MemberBriefDTO  `json:"members"`
	Blocks         []BlockDTO        `json:"blocks"`
	CreatedAt      time.Time         `json:"createdAt"`
	UpdatedAt      time.Time         `json:"updatedAt"`
}

// OverviewDTO 仪表盘统计。
type OverviewDTO struct {
	ProjectTotal int64              `json:"projectTotal"`
	Published    int64              `json:"published"`
	Draft        int64              `json:"draft"`
	Recommended  int64              `json:"recommended"`
	ViewTotal    uint64             `json:"viewTotal"`
	Recent       []AdminProjectItem `json:"recent"`
}
