package errorsx

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

var _ Error = (*XError)(nil)

type XError struct {
	Parent  error
	Message string
	ID      string
	Details map[string]string
}

func New(code string, message string) *XError {
	return &XError{ID: code, Message: message}
}

func NewF(code string, format string, args ...interface{}) *XError {
	return &XError{ID: code, Message: fmt.Sprintf(format, args...)}
}

func Wrap(parent error, code string, message string) *XError {
	return &XError{Parent: parent, ID: code, Message: message}
}

func WrapF(parent error, code string, format string, args ...interface{}) *XError {
	return &XError{Parent: parent, ID: code, Message: fmt.Sprintf(format, args...)}
}

func IsXError(err error) bool {
	var e *XError
	return errors.As(err, &e)
}

func (e *XError) GetMessage() string {
	return e.Message
}

func (e *XError) GetID() string {
	return e.ID
}

func (e *XError) GetParent() error {
	return e.Parent
}

func (e *XError) Error() string {
	if e.Parent != nil {
		return fmt.Sprintf("ID=%d Message=%s Parent=(%v)", e.ID, e.Message, e.Parent)
	}

	return fmt.Sprintf("ID=%d Message=%s", e.ID, e.Message)
}

func (e *XError) Unwrap() error {
	return e.Parent
}

func (e *XError) Is(target error) bool {
	var t *XError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	if t.ID != "" && t.ID != e.ID {
		return false
	}
	if t.Message != "" && t.Message != e.Message {
		return false
	}
	if t.Parent != nil && !errors.Is(e.Parent, t.Parent) {
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

func (e *XError) ResetMessage(message string) {
	e.Message = message
}

func (e *XError) SetDetails(details map[string]string) {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	for k, v := range details {
		e.Details[k] = v
	}
}

func (e *XError) GetDetails() map[string]string {
	details := make(map[string]string)
	for k, v := range e.Details {
		details[k] = v
	}
	return details
}

func (e *XError) HttpStatusCode() int {
	if e.Parent != nil {
		var p *XError
		if errors.As(e.Parent, &p) {
			return p.HttpStatusCode()
		}
	}
	return http.StatusInternalServerError
}
