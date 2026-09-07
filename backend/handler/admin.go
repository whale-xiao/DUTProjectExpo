package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// AdminHandler 后台项目接口（整组走 JWT 中间件）。
type AdminHandler struct {
	svc *service.AdminService
}

func NewAdminHandler(svc *service.AdminService) *AdminHandler {
	return &AdminHandler{svc: svc}
}

func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response.Fail(c, errcode.InvalidParam, "ID 不合法")
		return 0, false
	}
	return uint(id), true
}

// List GET /api/admin/projects
func (h *AdminHandler) List(c *gin.Context) {
	var q dto.AdminProjectQuery
	_ = c.ShouldBindQuery(&q)
	data, err := h.svc.ListProjects(q)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// Detail GET /api/admin/projects/:id（编辑回显）
func (h *AdminHandler) Detail(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	data, err := h.svc.Detail(id)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// Create POST /api/admin/projects
func (h *AdminHandler) Create(c *gin.Context) {
	var req dto.SaveProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	id, err := h.svc.Save(0, req, currentUID(c))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "创建成功", gin.H{"id": id})
}

// Update PUT /api/admin/projects/:id
func (h *AdminHandler) Update(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.SaveProjectReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	id, err := h.svc.Save(id, req, currentUID(c))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "保存成功", gin.H{"id": id})
}

// Delete DELETE /api/admin/projects/:id
func (h *AdminHandler) Delete(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已删除", nil)
}

// SetStatus PUT /api/admin/projects/:id/status（上架/下架/转草稿）
func (h *AdminHandler) SetStatus(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.SetStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "状态值不合法")
		return
	}
	if err := h.svc.SetStatus(id, req.Status); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "状态已更新", nil)
}

// SetRecommend PUT /api/admin/projects/:id/recommend
func (h *AdminHandler) SetRecommend(c *gin.Context) {
	id, ok := parseIDParam(c)
	if !ok {
		return
	}
	var req dto.RecommendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "请求参数有误")
		return
	}
	if err := h.svc.SetRecommend(id, req); err != nil {
		writeServiceError(c, err)
		return
	}
	response.OKWithMessage(c, "已更新", nil)
}

// Overview GET /api/admin/overview
func (h *AdminHandler) Overview(c *gin.Context) {
	data, err := h.svc.Overview()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

func currentUID(c *gin.Context) uint {
	if v, ok := c.Get("uid"); ok {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}
