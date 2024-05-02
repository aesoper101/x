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
// 用于表示用户未登录、未授权等。
type Unauthenticated interface {
	error
	IsUnauthenticated()
}

type UnauthenticatedError struct {
	*XError
}

func ThrowUnauthenticated(parent error, id, message string) error {
	return &UnauthenticatedError{Wrap(parent, id, message)}
}

func ThrowUnauthenticatedF(parent error, id, format string, a ...interface{}) error {
	return ThrowUnauthenticated(parent, id, fmt.Sprintf(format, a...))
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
