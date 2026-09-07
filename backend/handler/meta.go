package handler

import (
	"github.com/gin-gonic/gin"

	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// MetaHandler 分类/标签接口（前台列表公开；管理写操作在 /api/admin 下）。
type MetaHandler struct {
	svc *service.MetaService
}

func NewMetaHandler(svc *service.MetaService) *MetaHandler {
	return &MetaHandler{svc: svc}
}

// Categories GET /api/categories（前台/后台列表共用）
func (h *MetaHandler) Categories(c *gin.Context) {
	data, err := h.svc.Categories()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// Tags GET /api/tags
func (h *MetaHandler) Tags(c *gin.Context) {
	data, err := h.svc.Tags()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// CategoryCreate POST /api/admin/categories
func (h *MetaHandler) CategoryCreate(c *gin.Context) {
	var req dto.CategorySaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	if err := h.svc.SaveCategory(0, req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "创建成功", nil)
}

// CategoryUpdate PUT /api/admin/categories/:id
func (h *MetaHandler) CategoryUpdate(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.CategorySaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	if err := h.svc.SaveCategory(id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已更新", nil)
}

// CategoryDelete DELETE /api/admin/categories/:id
func (h *MetaHandler) CategoryDelete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteCategory(id); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已删除", nil)
}

// TagCreate POST /api/admin/tags
func (h *MetaHandler) TagCreate(c *gin.Context) {
	var req dto.TagSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	if err := h.svc.SaveTag(0, req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "创建成功", nil)
}

// TagUpdate PUT /api/admin/tags/:id
func (h *MetaHandler) TagUpdate(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.TagSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	if err := h.svc.SaveTag(id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已更新", nil)
}

// TagDelete DELETE /api/admin/tags/:id
func (h *MetaHandler) TagDelete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteTag(id); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已删除", nil)
}
