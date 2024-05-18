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
// 对应 HTTP 状态码：403
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

func ThrowPermissionDenied(cause error, reason, message string) error {
	return &PermissionDeniedError{Wrap(cause, reason, message)}
}

func ThrowPermissionDeniedF(cause error, reason, format string, a ...interface{}) error {
	return ThrowPermissionDenied(cause, reason, fmt.Sprintf(format, a...))
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

func (err *PermissionDeniedError) Cause() error {
	return err.XError
}
