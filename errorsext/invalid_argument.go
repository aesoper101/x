package errorsext

import (
	"errors"
	"fmt"
)

var (
	_ InvalidArgument = (*InvalidArgumentError)(nil)
	_ Error           = (*InvalidArgumentError)(nil)
)

// InvalidArgument 参数错误
// 对应 HTTP 状态码：400 Bad Request
// 参数错误，例如：参数为空、参数格式不正确等
// 请求字段 x.y.z 是 xxx，预期为 [yyy, zzz] 内的一个。
type InvalidArgument interface {
	error
	IsInvalidArgument()
}

type InvalidArgumentError struct {
	*XError
}

func ThrowInvalidArgument(cause error, reason, message string) error {
	return &InvalidArgumentError{Wrap(cause, reason, message)}
}

func ThrowInvalidArgumentF(cause error, reason, format string, a ...interface{}) error {
	return ThrowInvalidArgument(cause, reason, fmt.Sprintf(format, a...))
}

func (err *InvalidArgumentError) IsInvalidArgument() {}

func IsErrorInvalidArgument(err error) bool {
	var invalidArgument InvalidArgument
	ok := errors.As(err, &invalidArgument)
	return ok
}

func (err *InvalidArgumentError) Is(target error) bool {
	var t *InvalidArgumentError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *InvalidArgumentError) Unwrap() error {
	return err.XError
}

func (err *InvalidArgumentError) Cause() error {
	return err.XError
}
