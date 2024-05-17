package httpx

import "time"

const OK = 200

type SuccessResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TraceId   string      `json:"traceId,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func Success(message string, data ...interface{}) *SuccessResponse {
	resp := &SuccessResponse{
		Code:      OK,
		Message:   message,
		Timestamp: time.Now(),
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	return resp
}

type ErrorResponse struct {
	// 错误码，跟 http-status 一致，并且在 grpc 中可以转换成 grpc-status
	Code int `json:"code"`
	// 错误原因，定义为业务判定错误码
	Reason string `json:"reason"`
	// 错误信息，为用户可读的信息，可作为用户提示内容
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	TraceId   string                 `json:"traceId,omitempty"`
	RequestID string                 `json:"requestId,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}
