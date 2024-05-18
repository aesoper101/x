package errorsx

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Error = (*XError)(nil)

type XError struct {
	cause   error
	message string
	reason  string
	details map[string]interface{}
}

func New(reason string, message string) *XError {
	return &XError{reason: reason, message: message}
}

func NewF(reason string, format string, args ...interface{}) *XError {
	return &XError{reason: reason, message: fmt.Sprintf(format, args...)}
}

func Wrap(cause error, reason string, message string) *XError {
	return &XError{cause: cause, reason: reason, message: message}
}

func WrapF(cause error, reason string, format string, args ...interface{}) *XError {
	return &XError{cause: cause, reason: reason, message: fmt.Sprintf(format, args...)}
}

func IsXError(err error) bool {
	var e *XError
	return errors.As(err, &e)
}

func (e *XError) Message() string {
	return e.message
}

func (e *XError) Reason() string {
	return e.reason
}

func (e *XError) Cause() error {
	return e.cause
}

func (e *XError) Error() string {
	if e.Cause() != nil {
		return fmt.Sprintf("reason=%s message=%s cause=%s", e.reason, e.message, e.Cause())
	}
	return fmt.Sprintf("reason=%s message=%s", e.reason, e.message)
}

func (e *XError) Unwrap() error {
	return e.Cause()
}

func (e *XError) Is(target error) bool {
	var t *XError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	if t.reason != "" && t.reason != e.reason {
		return false
	}
	if t.message != "" && t.message != e.message {
		return false
	}
	if t.Cause() != nil && !errors.Is(e.Cause(), t.Cause()) {
		return false
	}

	return true
}

func (e *XError) As(target interface{}) bool {
	_, ok := target.(**XError)
	if !ok {
		return false
	}
	reflect.Indirect(reflect.ValueOf(target)).Set(reflect.ValueOf(e))
	return true
}

func (e *XError) WithMessage(message string) {
	e.message = message
}

func (e *XError) WithDetails(details map[string]interface{}) {
	if e.details == nil {
		e.details = make(map[string]interface{})
	}
	for k, v := range details {
		e.details[k] = v
	}
}

func (e *XError) Details() map[string]interface{} {
	details := make(map[string]interface{})
	for k, v := range e.details {
		details[k] = v
	}
	return details
}
