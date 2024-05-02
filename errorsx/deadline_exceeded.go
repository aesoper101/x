package errorsx

import (
	"errors"
	"fmt"
)

var (
	_ DeadlineExceeded = (*DeadlineExceededError)(nil)
	_ Error            = (*DeadlineExceededError)(nil)
)

// DeadlineExceeded 是操作超时错误。
// 操作超时错误，通常是因为资源被其他进程占用，或者操作本身需要很长时间。
type DeadlineExceeded interface {
	error
	IsDeadlineExceeded()
}

type DeadlineExceededError struct {
	*XError
}

func ThrowDeadlineExceeded(parent error, id, message string) error {
	return &DeadlineExceededError{Wrap(parent, id, message)}
}

func ThrowDeadlineExceededF(parent error, id, format string, a ...interface{}) error {
	return ThrowDeadlineExceeded(parent, id, fmt.Sprintf(format, a...))
}

func (err *DeadlineExceededError) IsDeadlineExceeded() {}

func IsDeadlineExceeded(err error) bool {
	var deadlineExceeded DeadlineExceeded
	ok := errors.As(err, &deadlineExceeded)
	return ok
}

func (err *DeadlineExceededError) Is(target error) bool {
	var t *DeadlineExceededError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *DeadlineExceededError) Unwrap() error {
	return err.XError
}
