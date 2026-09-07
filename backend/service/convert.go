package service

import (
	"showcase-backend/dto"
	"showcase-backend/model"
)

// cardFrom 把 Project + 已批量取的 tags 组装成卡片 DTO（列表/首页复用）。
func cardFrom(p model.Project, tags []string) dto.CardDTO {
	if tags == nil {
		tags = []string{}
	}
	return dto.CardDTO{
		ID:            p.ID,
		Name:          p.Name,
		Summary:       p.Summary,
		CoverURL:      p.CoverURL,
		Tags:          tags,
		IsRecommended: p.IsRecommended,
		CategoryID:    p.CategoryID,
		CreatedAt:     p.CreatedAt,
	}
}

// collectIDs 批量取 ids 供一次 IN 查标签。
func collectIDs(projects []model.Project) []uint {
	ids := make([]uint, 0, len(projects))
	for i := range projects {
		ids = append(ids, projects[i].ID)
	}
	return ids
}
