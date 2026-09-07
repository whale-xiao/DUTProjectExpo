package handler

import (
	"github.com/gin-gonic/gin"

	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// UploadHandler 图片上传（本地存储）。
type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

// Save POST /api/admin/upload（multipart form 字段 file）
func (h *UploadHandler) Save(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 40000, "请选择要上传的图片")
		return
	}
	url, err := h.svc.Save(file)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, gin.H{"url": url})
}
