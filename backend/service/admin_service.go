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

// AdminService 后台项目业务（管理端任意状态可见；整单保存走事务）。
type AdminService struct {
	db       *gorm.DB
	repo     *repository.AdminProjectRepo
	projRepo *repository.ProjectRepo
	catRepo  *repository.CategoryRepo
}

func NewAdminService(db *gorm.DB) *AdminService {
	return &AdminService{
		db:       db,
		repo:     repository.NewAdminProjectRepo(db),
		projRepo: repository.NewProjectRepo(db),
		catRepo:  repository.NewCategoryRepo(db),
	}
}

// catNames 分类 id -> 名称 映射。
func (s *AdminService) catNames() (map[uint]string, error) {
	cats, err := s.catRepo.ListAll()
	if err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(cats))
	for _, c := range cats {
		m[c.ID] = c.Name
	}
	return m, nil
}

func (s *AdminService) toItems(projects []model.Project) ([]dto.AdminProjectItem, error) {
	ids := collectIDs(projects)
	tagMap, err := s.projRepo.TagsByProjectIDs(ids)
	if err != nil {
		return nil, err
	}
	catMap, err := s.catNames()
	if err != nil {
		return nil, err
	}
	items := make([]dto.AdminProjectItem, 0, len(projects))
	for _, p := range projects {
		tags := tagMap[p.ID]
		if tags == nil {
			tags = []string{}
		}
		items = append(items, dto.AdminProjectItem{
			ID:             p.ID,
			Name:           p.Name,
			Summary:        p.Summary,
			CoverURL:       p.CoverURL,
			CategoryID:     p.CategoryID,
			CategoryName:   catMap[p.CategoryID],
			Status:         p.Status,
			IsRecommended:  p.IsRecommended,
			RecommendOrder: p.RecommendOrder,
			ViewCount:      p.ViewCount,
			Tags:           tags,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		})
	}
	return items, nil
}

// ListProjects 后台项目列表。
func (s *AdminService) ListProjects(q dto.AdminProjectQuery) (*dto.PageData[dto.AdminProjectItem], error) {
	projects, total, err := s.repo.ListPageAdmin(q)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	items, err := s.toItems(projects)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	page, size := normalizePage(q.Page, q.PageSize)
	return &dto.PageData[dto.AdminProjectItem]{
		List:     items,
		Total:    total,
		Page:     page,
		PageSize: size,
		HasMore:  int64(page)*int64(size) < total,
	}, nil
}

// Detail 后台编辑回显（任意状态）。
func (s *AdminService) Detail(id uint) (*dto.AdminDetailDTO, error) {
	p, err := s.repo.FindRawByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewError(errcode.NotFound, "项目不存在")
		}
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	blocks, err := s.projRepo.FindBlocks(id)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	members, err := s.projRepo.FindMembers(id)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}
	tagMap, err := s.projRepo.TagsByProjectIDs([]uint{id})
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询项目失败")
	}

	d := &dto.AdminDetailDTO{
		ID: p.ID, Name: p.Name, Summary: p.Summary, CoverURL: p.CoverURL,
		Status: p.Status, IsRecommended: p.IsRecommended, RecommendOrder: p.RecommendOrder,
		ViewCount: p.ViewCount, Tags: tagMap[id],
		Members:   make([]dto.MemberBriefDTO, 0, len(members)),
		Blocks:    make([]dto.BlockDTO, 0, len(blocks)),
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
	if d.Tags == nil {
		d.Tags = []string{}
	}
	if c, cerr := s.catRepo.FindByID(p.CategoryID); cerr == nil {
		d.Category = &dto.CategoryBriefDTO{ID: c.ID, Name: c.Name}
	}
	for _, m := range members {
		d.Members = append(d.Members, dto.MemberBriefDTO{Name: m.Name, Role: m.Role})
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
		d.Blocks = append(d.Blocks, bd)
	}
	return d, nil
}

// Save 后台整单保存（新建 id=0 / 编辑 id>0），返回项目 id。
func (s *AdminService) Save(id uint, req dto.SaveProjectReq, creatorID uint) (uint, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Summary = strings.TrimSpace(req.Summary)
	if req.Name == "" || req.Summary == "" {
		return 0, NewError(errcode.InvalidParam, "项目名称与一句话介绍不能为空")
	}
	exists, err := s.repo.NameExists(req.Name, id)
	if err != nil {
		return 0, NewError(errcode.InternalError, "保存失败")
	}
	if exists {
		return 0, NewError(errcode.Conflict, "项目名称已存在，请换一个")
	}
	if _, err := s.catRepo.FindByID(req.CategoryID); err != nil {
		return 0, NewError(errcode.InvalidParam, "所选分类不存在")
	}
	if id > 0 {
		if _, err := s.repo.FindRawByID(id); err != nil {
			return 0, NewError(errcode.NotFound, "项目不存在")
		}
	}

	status := req.Status
	if status == "" {
		status = model.ProjectStatusDraft
	}
	p := model.Project{
		Name: req.Name, Summary: req.Summary, CategoryID: req.CategoryID,
		CoverURL: req.CoverURL, Status: status,
		IsRecommended: req.IsRecommended, RecommendOrder: req.RecommendOrder,
	}

	tagNames := uniqueNonEmpty(req.Tags)
	members := make([]model.ProjectMember, 0, len(req.Members))
	for _, m := range req.Members {
		if strings.TrimSpace(m.Name) == "" {
			continue
		}
		members = append(members, model.ProjectMember{Name: strings.TrimSpace(m.Name), Role: strings.TrimSpace(m.Role)})
	}
	blocks := make([]model.ProjectContentBlock, 0, len(req.Blocks))
	for _, b := range req.Blocks {
		mb := model.ProjectContentBlock{BlockType: model.BlockType(b.Type)}
		switch mb.BlockType {
		case model.BlockHeading, model.BlockParagraph:
			mb.TextContent = b.Text
		case model.BlockImage:
			mb.ImageURL = b.URL
			mb.ImageCaption = b.Caption
		}
		blocks = append(blocks, mb)
	}

	newID, err := s.repo.SaveFull(id, creatorID, p, tagNames, members, blocks)
	if err != nil {
		return 0, NewError(errcode.InternalError, "保存失败")
	}
	return newID, nil
}

// Delete 删除项目（级联清理）。
func (s *AdminService) Delete(id uint) error {
	if _, err := s.repo.FindRawByID(id); err != nil {
		return NewError(errcode.NotFound, "项目不存在")
	}
	if err := s.repo.DeleteCascade(id); err != nil {
		return NewError(errcode.InternalError, "删除失败")
	}
	return nil
}

// SetStatus 上下架。
func (s *AdminService) SetStatus(id uint, status string) error {
	if _, err := s.repo.FindRawByID(id); err != nil {
		return NewError(errcode.NotFound, "项目不存在")
	}
	if err := s.repo.SetStatus(id, status); err != nil {
		return NewError(errcode.InternalError, "操作失败")
	}
	return nil
}

// SetRecommend 推荐位与排序。
func (s *AdminService) SetRecommend(id uint, req dto.RecommendReq) error {
	if _, err := s.repo.FindRawByID(id); err != nil {
		return NewError(errcode.NotFound, "项目不存在")
	}
	if err := s.repo.SetRecommend(id, req.IsRecommended, req.RecommendOrder); err != nil {
		return NewError(errcode.InternalError, "操作失败")
	}
	return nil
}

// Overview 仪表盘统计。
func (s *AdminService) Overview() (*dto.OverviewDTO, error) {
	stats, err := s.repo.Overview()
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询统计失败")
	}
	items, err := s.toItems(stats.Recent)
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询统计失败")
	}
	return &dto.OverviewDTO{
		ProjectTotal: stats.ProjectTotal,
		Published:    stats.Published,
		Draft:        stats.Draft,
		Recommended:  stats.Recommended,
		ViewTotal:    stats.ViewTotal,
		Recent:       items,
	}, nil
}

func uniqueNonEmpty(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	return out
}
