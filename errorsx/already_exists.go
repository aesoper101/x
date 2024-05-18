package errorsx

import (
	"errors"
)

var (
	_ AlreadyExists = (*AlreadyExistsError)(nil)
	_ Error         = (*AlreadyExistsError)(nil)
)

// AlreadyExists 是已经存在的错误
// 对应状态码 409
// 例如：用户已经存在，数据库中已经存在等
type AlreadyExists interface {
	error
	IsAlreadyExists()
}

type AlreadyExistsError struct {
	*XError
}

func ThrowAlreadyExists(cause error, reason, message string) error {
	return &AlreadyExistsError{Wrap(cause, reason, message)}
}

func ThrowAlreadyExistsF(cause error, reason, format string, a ...interface{}) error {
	return &AlreadyExistsError{WrapF(cause, reason, format, a...)}
}

func (err *AlreadyExistsError) IsAlreadyExists() {}

func (err *AlreadyExistsError) Is(target error) bool {
	var t *AlreadyExistsError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func IsErrorAlreadyExists(err error) bool {
	var alreadyExists AlreadyExists
	ok := errors.As(err, &alreadyExists)
	return ok
}

func (err *AlreadyExistsError) Unwrap() error {
	return err.XError
}

func (err *AlreadyExistsError) Cause() error {
	return err.XError
}
