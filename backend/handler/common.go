package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"showcase-backend/pkg/errcode"
	"showcase-backend/pkg/response"
	"showcase-backend/service"
)

// writeServiceError 把 service 层错误统一转成响应：
// 业务错误(service.Error)按其 code 返回；其余视为内部错误。
func writeServiceError(c *gin.Context, err error) {
	var se *service.Error
	if errors.As(err, &se) {
		response.Fail(c, se.Code, se.Msg)
		return
	}
	response.Fail(c, errcode.InternalError, "")
}
