package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// ProjectHandler 前台项目接口。
type ProjectHandler struct {
	svc *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

// List GET /api/projects（全部 / 搜索 / 分类标签筛选 / 排序 / 分页）
func (h *ProjectHandler) List(c *gin.Context) {
	var q dto.ListQuery
	_ = c.ShouldBindQuery(&q) // 查询参数宽松绑定，非法值回落默认
	data, err := h.svc.List(q)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// Suggest GET /api/search/suggest?keyword=（联想，名称优先）
// 说明：09 契约写 /api/projects/suggest，与 /projects/:id 同层会被 Gin 判定通配冲突，
// 故独立为 /search/suggest；前端 api 层与契约文档已同步此差异。
func (h *ProjectHandler) Suggest(c *gin.Context) {
	data, err := h.svc.Suggest(c.Query("keyword"))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}

// Detail GET /api/projects/:id
func (h *ProjectHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id < 1 {
		response.Fail(c, errcode.InvalidParam, "项目 ID 不合法")
		return
	}
	data, err := h.svc.Detail(uint(id))
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}
