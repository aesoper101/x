package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ ResourceExhausted = (*ResourceExhaustedError)(nil)
	_ Error             = (*ResourceExhaustedError)(nil)
)

// ResourceExhausted 是资源耗尽的错误
// 例如：数据库连接池已满，内存不足等
type ResourceExhausted interface {
	error
	IsResourceExhausted()
}

type ResourceExhaustedError struct {
	*XError
}

func ThrowResourceExhausted(parent error, id, message string) error {
	return &ResourceExhaustedError{Wrap(parent, id, message)}
}

func ThrowResourceExhaustedF(parent error, id, format string, a ...interface{}) error {
	return ThrowResourceExhausted(parent, id, fmt.Sprintf(format, a...))
}

func (err *ResourceExhaustedError) IsResourceExhausted() {}

func IsResourceExhausted(err error) bool {
	var resourceExhausted ResourceExhausted
	ok := errors.As(err, &resourceExhausted)
	return ok
}

func (err *ResourceExhaustedError) Is(target error) bool {
	var t *ResourceExhaustedError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *ResourceExhaustedError) Unwrap() error {
	return err.XError
}
