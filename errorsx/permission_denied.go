package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ PermissionDenied = (*PermissionDeniedError)(nil)
	_ Error            = (*PermissionDeniedError)(nil)
)

// PermissionDenied 是权限拒绝的错误
//
// 权限拒绝的错误，通常表示用户没有权限执行某些操作。
//
// 例如：
//   - 用户没有权限查看某个资源
//   - 用户没有权限修改某个资源
//   - 用户没有权限删除某个资源
type PermissionDenied interface {
	error
	IsPermissionDenied()
}

type PermissionDeniedError struct {
	*XError
}

func ThrowPermissionDenied(parent error, id, message string) error {
	return &PermissionDeniedError{Wrap(parent, id, message)}
}

func ThrowPermissionDeniedF(parent error, id, format string, a ...interface{}) error {
	return ThrowPermissionDenied(parent, id, fmt.Sprintf(format, a...))
}

func (err *PermissionDeniedError) IsPermissionDenied() {}

func IsPermissionDenied(err error) bool {
	var permissionDenied PermissionDenied
	ok := errors.As(err, &permissionDenied)
	return ok
}

func (err *PermissionDeniedError) Is(target error) bool {
	var t *PermissionDeniedError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *PermissionDeniedError) Unwrap() error {
	return err.XError
}
