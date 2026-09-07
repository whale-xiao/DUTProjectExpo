package repository

import (
	"strings"

	"gorm.io/gorm"

	"showcase-backend/dto"
	"showcase-backend/model"
)

// AdminProjectRepo 后台项目数据访问（含事务整单保存 / 级联删除 / 统计）。
// 注意：后台"任意状态可见"，不套前台 published 过滤。
type AdminProjectRepo struct {
	db *gorm.DB
}

func NewAdminProjectRepo(db *gorm.DB) *AdminProjectRepo { return &AdminProjectRepo{db: db} }

// ListPageAdmin 后台项目列表（不含软删；可关键字/分类/状态过滤，默认按更新时间倒序）。
func (r *AdminProjectRepo) ListPageAdmin(q dto.AdminProjectQuery) ([]model.Project, int64, error) {
	db := r.db.Model(&model.Project{})
	if kw := strings.TrimSpace(q.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where(`(name LIKE ? OR summary LIKE ?
			OR EXISTS(SELECT 1 FROM project_tags pt JOIN tags t ON t.id = pt.tag_id
			         WHERE pt.project_id = projects.id AND t.name LIKE ?))`,
			like, like, like)
	}
	if q.CategoryID > 0 {
		db = db.Where("category_id = ?", q.CategoryID)
	}
	if q.Status != "" {
		db = db.Where("status = ?", q.Status)
	}

	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, size := normalizePage(q.Page, q.PageSize)
	var projects []model.Project
	err := db.Order("updated_at DESC, id DESC").
		Offset((page - 1) * size).Limit(size).
		Find(&projects).Error
	return projects, total, err
}

// FindRawByID 后台编辑回显：任意状态。
func (r *AdminProjectRepo) FindRawByID(id uint) (*model.Project, error) {
	var p model.Project
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// SaveFull 整单保存：项目基础信息 + 标签 + 成员 + 内容块，一个事务内全量替换后三者。
// id==0 表示新建，返回新 id；否则按 id 更新（保留浏览量等只读字段），返回该 id。
func (r *AdminProjectRepo) SaveFull(id uint, creatorID uint, p model.Project,
	tagNames []string, members []model.ProjectMember, blocks []model.ProjectContentBlock) (uint, error) {

	var pid uint
	err := r.db.Transaction(func(tx *gorm.DB) error {
		pid = id
		if id > 0 {
			updates := map[string]any{
				"name": p.Name, "summary": p.Summary, "description": p.Description,
				"category_id": p.CategoryID, "cover_url": p.CoverURL,
				"status": p.Status, "is_recommended": p.IsRecommended,
				"recommend_order": p.RecommendOrder,
			}
			if err := tx.Model(&model.Project{}).Where("id = ?", id).Updates(updates).Error; err != nil {
				return err
			}
		} else {
			p.CreatedBy = creatorID
			if err := tx.Create(&p).Error; err != nil {
				return err
			}
			pid = p.ID
		}

		// 标签：重建关联（按名解析，不存在则创建）
		if err := tx.Where("project_id = ?", pid).Delete(&model.ProjectTag{}).Error; err != nil {
			return err
		}
		for _, name := range tagNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			var tag model.Tag
			err := tx.Where("name = ?", name).First(&tag).Error
			if err == gorm.ErrRecordNotFound {
				tag = model.Tag{Name: name}
				if err := tx.Create(&tag).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			if err := tx.Create(&model.ProjectTag{ProjectID: pid, TagID: tag.ID}).Error; err != nil {
				return err
			}
		}

		// 成员：重建
		if err := tx.Where("project_id = ?", pid).Delete(&model.ProjectMember{}).Error; err != nil {
			return err
		}
		for i := range members {
			members[i].ProjectID = pid
			members[i].SortOrder = i + 1
			if err := tx.Create(&members[i]).Error; err != nil {
				return err
			}
		}

		// 内容块：重建并重排 sort_order
		if err := tx.Where("project_id = ?", pid).Delete(&model.ProjectContentBlock{}).Error; err != nil {
			return err
		}
		for i := range blocks {
			blocks[i].ProjectID = pid
			blocks[i].SortOrder = i + 1
			if err := tx.Create(&blocks[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return pid, nil
}

// NameExists 判断项目名是否已被占用（软删不计）。
func (r *AdminProjectRepo) NameExists(name string, exceptID uint) (bool, error) {
	var n int64
	q := r.db.Model(&model.Project{}).Where("name = ?", name)
	if exceptID > 0 {
		q = q.Where("id <> ?", exceptID)
	}
	if err := q.Count(&n).Error; err != nil {
		return false, err
	}
	return n > 0, nil
}

// DeleteCascade 删除项目：projects 软删；关联四表硬删（避免孤儿）。
func (r *AdminProjectRepo) DeleteCascade(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&model.Project{}).Error; err != nil {
			return err
		}
		for _, m := range []any{
			&model.ProjectTag{}, &model.ProjectContentBlock{},
			&model.ProjectMember{}, &model.ProjectImage{},
		} {
			if err := tx.Where("project_id = ?", id).Delete(m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// SetStatus 上/下架。
func (r *AdminProjectRepo) SetStatus(id uint, status string) error {
	return r.db.Model(&model.Project{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// SetRecommend 推荐位与排序。
func (r *AdminProjectRepo) SetRecommend(id uint, isRecommended bool, order int) error {
	return r.db.Model(&model.Project{}).
		Where("id = ?", id).
		Updates(map[string]any{"is_recommended": isRecommended, "recommend_order": order}).Error
}

// OverviewStats 仪表盘统计原始量（service 组装 DTO）。
type OverviewStats struct {
	ProjectTotal int64
	Published    int64
	Draft        int64
	Recommended  int64
	ViewTotal    uint64
	Recent       []model.Project
}

// Overview 仪表盘统计与最近编辑。
// 注意：多条统计各自独立查询，不得复用同一个 *gorm.DB 连续执行（避免状态串扰）。
func (r *AdminProjectRepo) Overview() (*OverviewStats, error) {
	out := &OverviewStats{}
	base := func() *gorm.DB { return r.db.Model(&model.Project{}) }

	if err := base().Count(&out.ProjectTotal).Error; err != nil {
		return nil, err
	}
	if err := base().Where("status = ?", model.ProjectStatusPublished).Count(&out.Published).Error; err != nil {
		return nil, err
	}
	if err := base().Where("status = ?", model.ProjectStatusDraft).Count(&out.Draft).Error; err != nil {
		return nil, err
	}
	if err := base().Where("is_recommended = ?", true).Count(&out.Recommended).Error; err != nil {
		return nil, err
	}

	var agg struct {
		Total uint64
	}
	if err := r.db.Raw("SELECT COALESCE(SUM(view_count),0) AS total FROM projects WHERE deleted_at IS NULL").
		Scan(&agg).Error; err != nil {
		return nil, err
	}
	out.ViewTotal = agg.Total

	if err := base().Order("updated_at DESC, id DESC").Limit(5).Find(&out.Recent).Error; err != nil {
		return nil, err
	}
	return out, nil
}
