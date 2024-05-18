package httpx

import (
	"errors"
	"net/http"
	"time"

	"github.com/aesoper101/x/errorsx"
)

var defaultConvertErrorFunc ConvertErrorFunc = func(err error) errorsx.Error {
	return errorsx.ThrowUnknown(err, "", err.Error()).(errorsx.Error)
}

func SetDefaultConvertErrorFunc(f ConvertErrorFunc) {
	defaultConvertErrorFunc = f
}

type ConvertErrorFunc func(err error) errorsx.Error

type JSONResponseRender interface {
	JSON(code int, data interface{})
}

type SuccessResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func Success(render JSONResponseRender, message string, data ...interface{}) {
	resp := &SuccessResponse{
		Code:      http.StatusOK,
		Message:   message,
		Timestamp: time.Now(),
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	render.JSON(http.StatusOK, resp)
}

type FailureResponse struct {
	// 错误码，跟 http-status 一致，并且在 grpc 中可以转换成 grpc-status
	Code int `json:"code"`
	// 错误原因，定义为业务判定错误码
	Reason string `json:"reason"`
	// 错误信息，为用户可读的信息，可作为用户提示内容
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

func failure(render JSONResponseRender, statusCode int, reason, message string, details map[string]interface{}) {
	resp := &FailureResponse{
		Code:      statusCode,
		Reason:    reason,
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}

	render.JSON(statusCode, resp)
}

func Failure(render JSONResponseRender, cause error) {
	var code int
	var reason, message string
	var details map[string]interface{}

	var causeX errorsx.Error
	if errors.As(cause, &causeX) {
		reason = causeX.Reason()
		message = causeX.Message()
		details = causeX.Details()
	} else {
		causeX = defaultConvertErrorFunc(cause)
	}

	switch {
	case errorsx.IsErrorAlreadyExists(cause):
		code = http.StatusConflict
	case errorsx.IsDeadlineExceeded(cause):
		code = http.StatusRequestTimeout
	case errorsx.IsErrorInvalidArgument(cause):
		code = http.StatusBadRequest
	case errorsx.IsNotFound(cause):
		code = http.StatusNotFound
	case errorsx.IsPreconditionFailed(cause):
		code = http.StatusPreconditionFailed
	case errorsx.IsResourceExhausted(cause):
	case errorsx.IsUnauthenticated(cause):
		code = http.StatusUnauthorized
	case errorsx.IsUnavailable(cause):
		code = http.StatusServiceUnavailable
	case errorsx.IsUnimplemented(cause):
		code = http.StatusNotImplemented
	case errorsx.IsPermissionDenied(cause):
		code = http.StatusForbidden
	case errorsx.IsInternal(cause):
		code = http.StatusInternalServerError
		// 为了安全，不应该将内部错误暴露给用户
		reason = "InternalError"
		message = "oops, something went wrong"
	default:
		code = http.StatusInternalServerError
		// 为了安全，不应该将内部错误暴露给用户
		reason = "UnknownError"
		message = "oops, something went wrong"
	}

	failure(render, code, reason, message, details)
}
