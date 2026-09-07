package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"showcase-backend/config"
	"showcase-backend/pkg/errcode"
	pkgjwt "showcase-backend/pkg/jwt"
	"showcase-backend/pkg/response"
)

// JWT 鉴权中间件：解析 Authorization: Bearer <token>，成功后注入 uid/role/claims。
func JWT(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		const prefix = "Bearer "
		if !strings.HasPrefix(header, prefix) {
			response.Fail(c, errcode.Unauthorized, "")
			c.Abort()
			return
		}
		claims, err := pkgjwt.Parse(cfg.JWTSecret, strings.TrimPrefix(header, prefix))
		if err != nil {
			response.Fail(c, errcode.Unauthorized, "登录已过期，请重新登录")
			c.Abort()
			return
		}
		c.Set("uid", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("claims", claims)
		c.Next()
	}
}
