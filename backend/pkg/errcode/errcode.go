package errcode

// 统一业务错误码。约定：code=0 成功；非 0 为业务错误。
// 分位约定：4xxxx 客户端/业务，5xxxx 服务端。
const (
	OK            = 0
	InvalidParam  = 40000 // 参数不合法
	Unauthorized  = 40100 // 未登录 / Token 失效
	Forbidden     = 40300 // 无权限
	NotFound      = 40400 // 资源不存在
	Conflict      = 40900 // 冲突（如名称重复）
	InternalError = 50000 // 服务端内部错误
)

// Text 返回错误码默认文案（后端优先用具体 message，前端可据此兜底）。
var Text = map[int]string{
	OK:            "ok",
	InvalidParam:  "参数不合法",
	Unauthorized:  "请先登录",
	Forbidden:     "没有权限执行该操作",
	NotFound:      "资源不存在",
	Conflict:      "资源冲突",
	InternalError: "服务内部错误",
}
