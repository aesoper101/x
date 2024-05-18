package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ Unauthenticated = (*UnauthenticatedError)(nil)
	_ Error           = (*UnauthenticatedError)(nil)
)

// Unauthenticated 是未认证的错误。
// 对应于 HTTP 401 Unauthorized。
// 用于表示用户未登录、未授权等。
type Unauthenticated interface {
	error
	IsUnauthenticated()
}

type UnauthenticatedError struct {
	*XError
}

func ThrowUnauthenticated(cause error, reason, message string) error {
	return &UnauthenticatedError{Wrap(cause, reason, message)}
}

func ThrowUnauthenticatedF(cause error, reason, format string, a ...interface{}) error {
	return ThrowUnauthenticated(cause, reason, fmt.Sprintf(format, a...))
}

func (err *UnauthenticatedError) IsUnauthenticated() {}

func IsUnauthenticated(err error) bool {
	var unauthenticated Unauthenticated
	ok := errors.As(err, &unauthenticated)
	return ok
}

func (err *UnauthenticatedError) Is(target error) bool {
	var t *UnauthenticatedError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *UnauthenticatedError) Unwrap() error {
	return err.XError
}

func (err *UnauthenticatedError) Cause() error {
	return err.XError
}
