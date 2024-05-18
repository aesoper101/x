package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ Unimplemented = (*UnimplementedError)(nil)
	_ Error         = (*UnimplementedError)(nil)
)

// Unimplemented 是未实现错误
// 对应状态码 501
// 用于表示未实现的功能，比如未实现的接口方法
type Unimplemented interface {
	error
	IsUnimplemented()
}

type UnimplementedError struct {
	*XError
}

func ThrowUnimplemented(cause error, reason, message string) error {
	return &UnimplementedError{Wrap(cause, reason, message)}
}

func ThrowUnimplementedF(cause error, reason, format string, a ...interface{}) error {
	return ThrowUnimplemented(cause, reason, fmt.Sprintf(format, a...))
}

func (err *UnimplementedError) IsUnimplemented() {}

func IsUnimplemented(err error) bool {
	var unimplemented Unimplemented
	ok := errors.As(err, &unimplemented)
	return ok
}

func (err *UnimplementedError) Is(target error) bool {
	var t *UnimplementedError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnimplementedError) Unwrap() error {
	return err.XError
}

func (err *UnimplementedError) Cause() error {
	return err.XError
}
