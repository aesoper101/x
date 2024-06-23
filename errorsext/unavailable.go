package errorsext

import (
	"errors"
	"fmt"
)

var (
	_ Unavailable = (*UnavailableError)(nil)
	_ Error       = (*UnavailableError)(nil)
)

// Unavailable 是不可用错误
// 对应于 HTTP 状态码 503
// 例如：服务不可用，资源不可用等
type Unavailable interface {
	error
	IsUnavailable()
}

type UnavailableError struct {
	*XError
}

func ThrowUnavailable(cause error, reason, message string) error {
	return &UnavailableError{Wrap(cause, reason, message)}
}

func ThrowUnavailableF(cause error, reason, format string, a ...interface{}) error {
	return ThrowUnavailable(cause, reason, fmt.Sprintf(format, a...))
}

func (err *UnavailableError) IsUnavailable() {}

func IsUnavailable(err error) bool {
	var unavailable Unavailable
	ok := errors.As(err, &unavailable)
	return ok
}

func (err *UnavailableError) Is(target error) bool {
	var t *UnavailableError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnavailableError) Unwrap() error {
	return err.XError
}

func (err *UnavailableError) Cause() error {
	return err.XError
}
