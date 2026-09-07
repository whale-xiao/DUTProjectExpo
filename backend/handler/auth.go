package handler

import (
	"github.com/gin-gonic/gin"

	"showcase-backend/dto"
	"showcase-backend/pkg/errcode"
	pkgjwt "showcase-backend/pkg/jwt"
	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// AuthHandler 认证接口。
type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc: svc} }

// Login POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errcode.InvalidParam, "用户名和密码不能为空")
		return
	}
	res, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.OK(c, res)
}

// Me GET /api/auth/me（需登录，返回当前用户身份）
func (h *AuthHandler) Me(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		response.Fail(c, errcode.Unauthorized, "")
		return
	}
	cl := claims.(*pkgjwt.Claims)
	response.OK(c, gin.H{
		"id":       cl.UserID,
		"username": cl.Username,
		"role":     cl.Role,
	})
}
