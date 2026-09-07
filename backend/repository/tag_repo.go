package repository

import (
	"gorm.io/gorm"

	"showcase-backend/model"
)

// TagRepo 标签数据访问。
type TagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo { return &TagRepo{db: db} }

// ListAll 全部标签。
func (r *TagRepo) ListAll() ([]model.Tag, error) {
	var tags []model.Tag
	err := r.db.Order("id ASC").Find(&tags).Error
	return tags, err
}

func (r *TagRepo) FindByName(name string) (*model.Tag, error) {
	var t model.Tag
	if err := r.db.Where("name = ?", name).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

// Create 新建标签（同名冲突返回 ErrDuplicate）。
func (r *TagRepo) Create(name, color string) error {
	if _, err := r.FindByName(name); err == nil {
		return ErrDuplicate
	}
	return r.db.Create(&model.Tag{Name: name, Color: color}).Error
}

// Update 更新标签。
func (r *TagRepo) Update(id uint, name, color string) error {
	var dup model.Tag
	err := r.db.Where("name = ? AND id <> ?", name, id).First(&dup).Error
	if err == nil {
		return ErrDuplicate
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Model(&model.Tag{}).Where("id = ?", id).
		Updates(map[string]any{"name": name, "color": color}).Error
}

// CountLinked 该标签被项目引用的数量。
func (r *TagRepo) CountLinked(id uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.ProjectTag{}).Where("tag_id = ?", id).Count(&n).Error
	return n, err
}

// DeleteByID 删除标签（仍被项目引用时返回 ErrRefInUse）。
func (r *TagRepo) DeleteByID(id uint) error {
	n, err := r.CountLinked(id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrRefInUse
	}
	return r.db.Where("id = ?", id).Delete(&model.Tag{}).Error
}
