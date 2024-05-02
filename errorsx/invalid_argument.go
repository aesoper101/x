package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ InvalidArgument = (*InvalidArgumentError)(nil)
	_ Error           = (*InvalidArgumentError)(nil)
)

// InvalidArgument 参数错误
// 参数错误，例如：参数为空、参数格式不正确等
type InvalidArgument interface {
	error
	IsInvalidArgument()
}

type InvalidArgumentError struct {
	*XError
}

func ThrowInvalidArgument(parent error, id, message string) error {
	return &InvalidArgumentError{Wrap(parent, id, message)}
}

func ThrowInvalidArgumentF(parent error, id, format string, a ...interface{}) error {
	return ThrowInvalidArgument(parent, id, fmt.Sprintf(format, a...))
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
