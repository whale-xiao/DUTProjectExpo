package service

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	"showcase-backend/repository"
)

// MetaService 分类/标签等元数据列表业务。
type MetaService struct {
	catRepo *repository.CategoryRepo
	tagRepo *repository.TagRepo
}

func NewMetaService(db *gorm.DB) *MetaService {
	return &MetaService{
		catRepo: repository.NewCategoryRepo(db),
		tagRepo: repository.NewTagRepo(db),
	}
}

// Categories 分类列表。
func (s *MetaService) Categories() ([]dto.CategoryItemDTO, error) {
	cats, err := s.catRepo.ListAll()
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询分类失败")
	}
	items := make([]dto.CategoryItemDTO, 0, len(cats))
	for _, c := range cats {
		items = append(items, dto.CategoryItemDTO{ID: c.ID, Name: c.Name, SortOrder: c.SortOrder})
	}
	return items, nil
}

// Tags 标签列表。
func (s *MetaService) Tags() ([]dto.TagItemDTO, error) {
	tags, err := s.tagRepo.ListAll()
	if err != nil {
		return nil, NewError(errcode.InternalError, "查询标签失败")
	}
	items := make([]dto.TagItemDTO, 0, len(tags))
	for _, t := range tags {
		items = append(items, dto.TagItemDTO{ID: t.ID, Name: t.Name, Color: t.Color})
	}
	return items, nil
}

// --- 后台：分类 / 标签管理（写） ---

// SaveCategory 新建(id=0)或编辑分类。
func (s *MetaService) SaveCategory(id uint, req dto.CategorySaveReq) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return NewError(errcode.InvalidParam, "分类名不能为空")
	}
	var err error
	if id > 0 {
		err = s.catRepo.Update(id, req.Name, req.Description, req.SortOrder)
	} else {
		err = s.catRepo.Create(req.Name, req.Description, req.SortOrder)
	}
	return s.mapMetaErr(err, "分类")
}

func (s *MetaService) DeleteCategory(id uint) error {
	return s.mapMetaErr(s.catRepo.DeleteByID(id), "分类")
}

// SaveTag 新建(id=0)或编辑标签。
func (s *MetaService) SaveTag(id uint, req dto.TagSaveReq) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return NewError(errcode.InvalidParam, "标签名不能为空")
	}
	var err error
	if id > 0 {
		err = s.tagRepo.Update(id, req.Name, req.Color)
	} else {
		err = s.tagRepo.Create(req.Name, req.Color)
	}
	return s.mapMetaErr(err, "标签")
}

func (s *MetaService) DeleteTag(id uint) error {
	return s.mapMetaErr(s.tagRepo.DeleteByID(id), "标签")
}

// mapMetaErr 仓储层哨兵错误 → 业务错误码。
func (s *MetaService) mapMetaErr(err error, name string) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, repository.ErrDuplicate):
		return NewError(errcode.Conflict, name+"名已存在")
	case errors.Is(err, repository.ErrRefInUse):
		return NewError(errcode.Conflict, "该"+name+"仍被引用，无法删除")
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewError(errcode.NotFound, name+"不存在")
	default:
		return NewError(errcode.InternalError, "操作失败")
	}
}
