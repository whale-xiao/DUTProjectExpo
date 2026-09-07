package service

import (
	"gorm.io/gorm"

	"showcase-backend/dto"
	"showcase-backend/model"
	"showcase-backend/pkg/errcode"
	"showcase-backend/repository"
)

// HomeService 首页聚合业务。
type HomeService struct {
	repo    *repository.ProjectRepo
	catRepo *repository.CategoryRepo
}

func NewHomeService(db *gorm.DB) *HomeService {
	return &HomeService{
		repo:    repository.NewProjectRepo(db),
		catRepo: repository.NewCategoryRepo(db),
	}
}

// Home 首页聚合：推荐（≤12）+ 分类 + 最新（3 兜底）。
func (s *HomeService) Home() (*dto.HomeDTO, error) {
	recs, err := s.repo.FindRecommended(12)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询推荐项目失败")
	}
	latest, err := s.repo.FindLatest(3)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询最新项目失败")
	}
	cats, err := s.catRepo.ListAll()
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询分类失败")
	}

	// 推荐 + 最新合并一次批量取标签，避免 N+1。
	all := make([]model.Project, 0, len(recs)+len(latest))
	all = append(all, recs...)
	all = append(all, latest...)
	tagMap, err := s.repo.TagsByProjectIDs(collectIDs(all))
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}

	recommended := make([]dto.CardDTO, 0, len(recs))
	for _, p := range recs {
		recommended = append(recommended, cardFrom(p, tagMap[p.ID]))
	}
	latestCards := make([]dto.CardDTO, 0, len(latest))
	for _, p := range latest {
		latestCards = append(latestCards, cardFrom(p, tagMap[p.ID]))
	}
	categoryItems := make([]dto.CategoryItemDTO, 0, len(cats))
	for _, c := range cats {
		categoryItems = append(categoryItems, dto.CategoryItemDTO{ID: c.ID, Name: c.Name, SortOrder: c.SortOrder})
	}

	return &dto.HomeDTO{
		Recommended: recommended,
		Categories:  categoryItems,
		Latest:      latestCards,
	}, nil
}
