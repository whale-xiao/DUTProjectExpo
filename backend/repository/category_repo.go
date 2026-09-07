package repository

import (
	"gorm.io/gorm"

	"showcase-backend/model"
)

// CategoryRepo 分类数据访问。
type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db: db} }

// ListAll 全部分类（按 sort_order）。
func (r *CategoryRepo) ListAll() ([]model.Category, error) {
	var cats []model.Category
	err := r.db.Order("sort_order ASC, id ASC").Find(&cats).Error
	return cats, err
}

// FindByID 查单个分类（供详情内嵌）。
func (r *CategoryRepo) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepo) FindByName(name string) (*model.Category, error) {
	var c model.Category
	if err := r.db.Where("name = ?", name).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

// Create 新建分类（同名冲突返回 ErrDuplicate）。
func (r *CategoryRepo) Create(name, description string, sortOrder int) error {
	if _, err := r.FindByName(name); err == nil {
		return ErrDuplicate
	}
	return r.db.Create(&model.Category{Name: name, Description: description, SortOrder: sortOrder}).Error
}

// Update 更新分类（改名撞名返回 ErrDuplicate）。
func (r *CategoryRepo) Update(id uint, name, description string, sortOrder int) error {
	var dup model.Category
	err := r.db.Where("name = ? AND id <> ?", name, id).First(&dup).Error
	if err == nil {
		return ErrDuplicate
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Model(&model.Category{}).Where("id = ?", id).
		Updates(map[string]any{"name": name, "description": description, "sort_order": sortOrder}).Error
}

// CountProjects 该分类下项目数。
func (r *CategoryRepo) CountProjects(id uint) (int64, error) {
	var n int64
	err := r.db.Model(&model.Project{}).Where("category_id = ?", id).Count(&n).Error
	return n, err
}

// DeleteByID 删除分类（仍有项目引用时返回 ErrRefInUse）。
func (r *CategoryRepo) DeleteByID(id uint) error {
	n, err := r.CountProjects(id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrRefInUse
	}
	return r.db.Where("id = ?", id).Delete(&model.Category{}).Error
}
