package errorsx

import (
	"errors"
	"fmt"
	"net/http"
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

func ThrowInternal(parent error, id, message string) error {
	return &InternalError{Wrap(parent, id, message)}
}

func ThrowInternalF(parent error, id, format string, a ...interface{}) error {
	return ThrowInternal(parent, id, fmt.Sprintf(format, a...))
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

func (err *InternalError) Unwrap() error {
	return err.XError
}

func (err *InternalError) HttpStatusCode() int {
	return http.StatusInternalServerError
}
