package service

// Error 业务层错误：携带统一错误码，handler 经 errors.As 转响应。
type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string { return e.Msg }

// NewError 构造业务错误（code 取 pkg/errcode 常量）。
func NewError(code int, msg string) *Error { return &Error{Code: code, Msg: msg} }
