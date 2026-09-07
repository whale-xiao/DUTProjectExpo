package repository

import "errors"

// 仓储层哨兵错误：service 层经 errors.Is 映射为业务错误码。
var (
	// ErrDuplicate 唯一约束冲突（同名分类/标签/项目已存在）。
	ErrDuplicate = errors.New("duplicate resource")
	// ErrRefInUse 仍被引用，禁止删除。
	ErrRefInUse = errors.New("resource in use")
)
