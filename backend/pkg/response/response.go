package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"showcase-backend/pkg/errcode"
)

// Body 统一响应结构 {code, message, data}（对齐 02-基本架构 §3.2）。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// OK 成功响应：data 直接可用。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: errcode.OK, Message: "ok", Data: data})
}

// OKWithMessage 成功但携带自定义提示。
func OKWithMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Body{Code: errcode.OK, Message: message, Data: data})
}

// Fail 业务失败：HTTP 用 200 + 业务码（前端拦截器统一解析），
// 也允许调用方传具体 HTTP 状态（如 404 页面语义）。
func Fail(c *gin.Context, code int, message string) {
	if message == "" {
		if m, ok := errcode.Text[code]; ok {
			message = m
		} else {
			message = "error"
		}
	}
	c.JSON(http.StatusOK, Body{Code: code, Message: message})
}

// FailHTTP 携带 HTTP 状态码的失败（用于 401/404 等需页面语义的场景）。
func FailHTTP(c *gin.Context, status, code int, message string) {
	if message == "" {
		if m, ok := errcode.Text[code]; ok {
			message = m
		} else {
			message = "error"
		}
	}
	c.JSON(status, Body{Code: code, Message: message})
}
