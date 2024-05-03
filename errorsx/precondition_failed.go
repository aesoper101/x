package errorsx

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	_ PreconditionFailed = (*PreconditionFailedError)(nil)
	_ Error              = (*PreconditionFailedError)(nil)
)

// PreconditionFailed 是前置条件失败的错误接口
// 对应 HTTP 状态码 412
// 用于表示前置条件不满足，导致当前操作无法继续执行
// 例如：资源 xxx 是非空目录，因此无法删除。
type PreconditionFailed interface {
	error
	IsPreconditionFailed()
}

type PreconditionFailedError struct {
	*XError
}

func ThrowPreconditionFailed(parent error, id, message string) error {
	return &PreconditionFailedError{Wrap(parent, id, message)}
}

func ThrowPreconditionFailedF(parent error, id, format string, a ...interface{}) error {
	return ThrowPreconditionFailed(parent, id, fmt.Sprintf(format, a...))
}

func (err *PreconditionFailedError) IsPreconditionFailed() {}

func IsPreconditionFailed(err error) bool {
	var preconditionFailed PreconditionFailed
	ok := errors.As(err, &preconditionFailed)
	return ok
}

func (err *PreconditionFailedError) Is(target error) bool {
	var t *PreconditionFailedError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return err.XError.Is(t.XError)
}

func (err *PreconditionFailedError) Unwrap() error {
	return err.XError
}

func (err *PreconditionFailedError) HttpStatusCode() int {
	return http.StatusPreconditionFailed
}
