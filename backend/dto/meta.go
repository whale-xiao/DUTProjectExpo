package dto

// CategoryItemDTO 分类列表项（前端 chips / 下拉）。
type CategoryItemDTO struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	SortOrder int    `json:"sortOrder"`
}

// CategoryBriefDTO 详情内嵌的分类简述。
type CategoryBriefDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// TagItemDTO 标签列表项。
type TagItemDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// HomeDTO 首页聚合（09 §4.3 GET /api/home）。
type HomeDTO struct {
	Recommended []CardDTO         `json:"recommended"`
	Categories  []CategoryItemDTO `json:"categories"`
	Latest      []CardDTO         `json:"latest"`
}
