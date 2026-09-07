package handler

import (
	"github.com/gin-gonic/gin"

	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// HomeHandler 首页聚合接口。
type HomeHandler struct {
	svc *service.HomeService
}

func NewHomeHandler(svc *service.HomeService) *HomeHandler {
	return &HomeHandler{svc: svc}
}

// Home GET /api/home
func (h *HomeHandler) Home(c *gin.Context) {
	data, err := h.svc.Home()
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, data)
}
