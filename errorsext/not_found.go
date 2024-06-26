package errorsext

import (
	"errors"
	"fmt"
)

// NotFound 是未找到的错误
// 对应 HTTP 状态码：404 Not Found
// 未找到的错误，通常表示请求的资源不存在。
// 例如：用户不存在，或者请求的资源不存在。
type NotFound interface {
	error
	IsNotFound()
}

type NotFoundError struct {
	*XError
}

func ThrowNotFound(cause error, reason, message string) error {
	return &NotFoundError{Wrap(cause, reason, message)}
}

func ThrowNotFoundF(cause error, reason, format string, a ...interface{}) error {
	return ThrowNotFound(cause, reason, fmt.Sprintf(format, a...))
}

func (err *NotFoundError) IsNotFound() {}

func IsNotFound(err error) bool {
	var notFound NotFound
	ok := errors.As(err, &notFound)
	return ok
}

func (err *NotFoundError) Is(target error) bool {
	var t *NotFoundError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *NotFoundError) Unwrap() error {
	return err.XError
}

func (err *NotFoundError) Cause() error {
	return err.XError
}
