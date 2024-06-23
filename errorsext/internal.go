package errorsext

import (
	"errors"
	"fmt"
)

var (
	_ Internal = (*InternalError)(nil)
	_ Error    = (*InternalError)(nil)
)

// Internal 接口用于标识内部错误，
// 对应状态码 500
// 内部错误通常表示程序逻辑错误，
// 内部错误通常不应该被用户看到，
type Internal interface {
	error
	IsInternal()
}

type InternalError struct {
	*XError
}

func ThrowInternal(cause error, reason, message string) error {
	return &InternalError{Wrap(cause, reason, message)}
}

func ThrowInternalF(cause error, reason, format string, a ...interface{}) error {
	return ThrowInternal(cause, reason, fmt.Sprintf(format, a...))
}

func (err *InternalError) IsInternal() {}

func IsInternal(err error) bool {
	var internal Internal
	ok := errors.As(err, &internal)
	return ok
}

func (err *InternalError) Is(target error) bool {
	var t *InternalError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}

	return err.XError.Is(t.XError)
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("InternalError: reason=%s, message=%s", err.Reason(), err.Message())
}

func (err *InternalError) Unwrap() error {
	return err.XError
}

func (err *InternalError) Cause() error {
	return err.XError
}
