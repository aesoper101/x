package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ Unknown = (*UnknownError)(nil)
	_ Error   = (*UnknownError)(nil)
)

// Unknown 未知错误
// 对应Http 500
// 未知错误，通常用于表示无法识别的错误
type Unknown interface {
	error
	IsUnknown()
}

type UnknownError struct {
	*XError
}

func ThrowUnknown(cause error, reason, message string) error {
	return &UnknownError{Wrap(cause, reason, message)}
}

func ThrowUnknownF(cause error, reason, format string, a ...interface{}) error {
	return ThrowUnknown(cause, reason, fmt.Sprintf(format, a...))
}

func (err *UnknownError) IsUnknown() {}

func IsUnknown(err error) bool {
	var unknown Unknown
	ok := errors.As(err, &unknown)
	return ok
}

func (err *UnknownError) Is(target error) bool {
	var t *UnknownError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnknownError) Unwrap() error {
	return err.XError
}

func (err *UnknownError) Cause() error {
	return err.XError
}
