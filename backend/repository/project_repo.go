package repository

import (
	"strings"

	"gorm.io/gorm"

	"showcase-backend/dto"
	"showcase-backend/model"
)

// ProjectRepo 项目数据访问（前台查询一律只取 published；软删由 GORM 自动过滤）。
type ProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) *ProjectRepo { return &ProjectRepo{db: db} }

// basePublished 上架项目的基础查询。
func (r *ProjectRepo) basePublished() *gorm.DB {
	return r.db.Model(&model.Project{}).Where("status = ?", model.ProjectStatusPublished)
}

// ListPage 前台列表：COUNT + 主查询共用同一条件 builder（见 10-调用链 §3）。
func (r *ProjectRepo) ListPage(q dto.ListQuery) ([]model.Project, int64, error) {
	db := r.basePublished()
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where(`(name LIKE ? OR summary LIKE ? OR description LIKE ?
			OR EXISTS(SELECT 1 FROM project_tags pt JOIN tags t ON t.id = pt.tag_id
			         WHERE pt.project_id = projects.id AND t.name LIKE ?)
			OR EXISTS(SELECT 1 FROM project_content_blocks cb
			         WHERE cb.project_id = projects.id AND cb.text_content LIKE ?))`,
			like, like, like, like, like)
	}
	if q.CategoryID > 0 {
		db = db.Where("category_id = ?", q.CategoryID)
	}
	if q.TagID > 0 {
		db = db.Where("EXISTS(SELECT 1 FROM project_tags pt2 WHERE pt2.project_id = projects.id AND pt2.tag_id = ?)", q.TagID)
	}

	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "is_recommended DESC, recommend_order ASC, id DESC"
	switch q.Sort {
	case "newest":
		order = "created_at DESC, id DESC"
	case "views":
		order = "view_count DESC, id DESC"
	}

	page, size := normalizePage(q.Page, q.PageSize)
	var projects []model.Project
	err := db.Order(order).Offset((page - 1) * size).Limit(size).Find(&projects).Error
	return projects, total, err
}

// FindRecommended 首页推荐（按推荐位排序）。
func (r *ProjectRepo) FindRecommended(limit int) ([]model.Project, error) {
	var ps []model.Project
	err := r.basePublished().
		Where("is_recommended = ?", true).
		Order("recommend_order ASC, id DESC").
		Limit(limit).
		Find(&ps).Error
	return ps, err
}

// FindLatest 首页"最新上架"兜底。
func (r *ProjectRepo) FindLatest(limit int) ([]model.Project, error) {
	var ps []model.Project
	err := r.basePublished().
		Order("created_at DESC, id DESC").
		Limit(limit).
		Find(&ps).Error
	return ps, err
}

// SearchSuggest 联想：名称命中优先（SQL 侧排序），上限由调用方给。
func (r *ProjectRepo) SearchSuggest(kw string, limit int) ([]model.Project, error) {
	like := "%" + kw + "%"
	db := r.basePublished().Where(`(name LIKE ? OR summary LIKE ? OR description LIKE ?
		OR EXISTS(SELECT 1 FROM project_tags pt JOIN tags t ON t.id = pt.tag_id
		         WHERE pt.project_id = projects.id AND t.name LIKE ?))`,
		like, like, like, like)
	db = db.Order(gorm.Expr("(name LIKE ?) DESC", like)).Order("id ASC").Limit(limit)
	var ps []model.Project
	err := db.Find(&ps).Error
	return ps, err
}

// FindPublishedDetail 前台详情（仅上架项目）。
func (r *ProjectRepo) FindPublishedDetail(id uint) (*model.Project, error) {
	var p model.Project
	err := r.db.Model(&model.Project{}).
		Where("id = ? AND status = ?", id, model.ProjectStatusPublished).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// FindBlocks 详情正文块（按 sort_order 有序）。
func (r *ProjectRepo) FindBlocks(projectID uint) ([]model.ProjectContentBlock, error) {
	var blocks []model.ProjectContentBlock
	err := r.db.Where("project_id = ?", projectID).
		Order("sort_order ASC, id ASC").
		Find(&blocks).Error
	return blocks, err
}

// FindMembers 详情成员（按 sort_order 有序）。
func (r *ProjectRepo) FindMembers(projectID uint) ([]model.ProjectMember, error) {
	var members []model.ProjectMember
	err := r.db.Where("project_id = ?", projectID).
		Order("sort_order ASC, id ASC").
		Find(&members).Error
	return members, err
}

// IncrementView 浏览量 +1（单条原子 UPDATE）。
func (r *ProjectRepo) IncrementView(projectID uint) error {
	return r.db.Model(&model.Project{}).
		Where("id = ?", projectID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// TagsByProjectIDs 批量取标签（一次 IN 查询，避免 N+1）。
func (r *ProjectRepo) TagsByProjectIDs(ids []uint) (map[uint][]string, error) {
	m := map[uint][]string{}
	if len(ids) == 0 {
		return m, nil
	}
	var rows []struct {
		ProjectID uint
		Name      string
	}
	err := r.db.Table("project_tags pt").
		Select("pt.project_id, t.name").
		Joins("JOIN tags t ON t.id = pt.tag_id").
		Where("pt.project_id IN ?", ids).
		Order("t.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		m[row.ProjectID] = append(m[row.ProjectID], row.Name)
	}
	return m, nil
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
