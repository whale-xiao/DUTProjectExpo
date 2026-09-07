package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"showcase-backend/dto"
	"showcase-backend/model"
	"showcase-backend/pkg/errcode"
	"showcase-backend/repository"
)

// ProjectService 前台项目业务。
type ProjectService struct {
	repo    *repository.ProjectRepo
	catRepo *repository.CategoryRepo
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		repo:    repository.NewProjectRepo(db),
		catRepo: repository.NewCategoryRepo(db),
	}
}

// List 项目列表（全部 / 搜索 / 分类与标签筛选 / 排序 / 分页）。
func (s *ProjectService) List(q dto.ListQuery) (*dto.PageData[dto.CardDTO], error) {
	projects, total, err := s.repo.ListPage(q)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	tagMap, err := s.repo.TagsByProjectIDs(collectIDs(projects))
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}

	cards := make([]dto.CardDTO, 0, len(projects))
	for _, p := range projects {
		cards = append(cards, cardFrom(p, tagMap[p.ID]))
	}

	page, size := normalizePage(q.Page, q.PageSize)
	return &dto.PageData[dto.CardDTO]{
		List:     cards,
		Total:    total,
		Page:     page,
		PageSize: size,
		HasMore:  int64(page)*int64(size) < total,
	}, nil
}

// Suggest 搜索联想（最多 6 条，名称优先）。
func (s *ProjectService) Suggest(kw string) ([]dto.SuggestItemDTO, error) {
	kw = strings.TrimSpace(kw)
	if kw == "" {
		return []dto.SuggestItemDTO{}, nil
	}
	projects, err := s.repo.SearchSuggest(kw, 6)
	if err != nil {
		return nil, NewError(errcode.InternalError, "搜索失败")
	}
	tagMap, err := s.repo.TagsByProjectIDs(collectIDs(projects))
	if err != nil {
		return nil, NewError(errcode.InternalError, "搜索失败")
	}
	items := make([]dto.SuggestItemDTO, 0, len(projects))
	for _, p := range projects {
		tags := tagMap[p.ID]
		if tags == nil {
			tags = []string{}
		}
		items = append(items, dto.SuggestItemDTO{
			ID:       p.ID,
			Name:     p.Name,
			Summary:  p.Summary,
			CoverURL: p.CoverURL,
			Tags:     tags,
		})
	}
	return items, nil
}

// Detail 前台详情：组装 category/tags/members/blocks 并执行浏览量 +1（一期策略）。
func (s *ProjectService) Detail(id uint) (*dto.DetailDTO, error) {
	p, err := s.repo.FindPublishedDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewError(errcode.NotFound, "项目不存在或已下架")
		}
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}

	blocks, err := s.repo.FindBlocks(id)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	members, err := s.repo.FindMembers(id)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	tagMap, err := s.repo.TagsByProjectIDs([]uint{id})
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}

	detail := &dto.DetailDTO{
		ID:          p.ID,
		Name:        p.Name,
		Summary:     p.Summary,
		CoverURL:    p.CoverURL,
		Tags:        tagMap[id],
		Members:     make([]dto.MemberBriefDTO, 0, len(members)),
		Blocks:      make([]dto.BlockDTO, 0, len(blocks)),
		PublishedAt: p.CreatedAt.Format("2006-01-02"),
	}
	if tags := detail.Tags; tags == nil {
		detail.Tags = []string{}
	}

	if c, cerr := s.catRepo.FindByID(p.CategoryID); cerr == nil {
		detail.Category = &dto.CategoryBriefDTO{ID: c.ID, Name: c.Name}
	}

	for _, m := range members {
		detail.Members = append(detail.Members, dto.MemberBriefDTO{Name: m.Name, Role: m.Role})
	}
	for _, b := range blocks {
		bd := dto.BlockDTO{Type: string(b.BlockType)}
		switch b.BlockType {
		case model.BlockHeading, model.BlockParagraph:
			bd.Text = b.TextContent
		case model.BlockImage:
			bd.URL = b.ImageURL
			bd.Caption = b.ImageCaption
		}
		detail.Blocks = append(detail.Blocks, bd)
	}

	// 浏览量 +1 并返回现值（失败不阻断详情；二期改 POST /views + Redis）。
	_ = s.repo.IncrementView(id)
	detail.ViewCount = p.ViewCount + 1
	return detail, nil
}

func normalizePage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}
