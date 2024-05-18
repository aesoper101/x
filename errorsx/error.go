package errorsx

// Error 错误接口
// inspire by https://cloud.google.com/apis/design/errors?hl=zh-cn
type Error interface {
	error
	// Cause 返回原始错误
	Cause() error
	// Unwrap 返回原始错误 兼容 Go 1.13
	Unwrap() error
	// Message 错误信息，为用户可读的信息，可作为用户提示内容
	Message() string
	WithMessage(message string)
	WithDetails(details map[string]interface{})
	// Details 错误元信息，为错误添加附加可扩展信息
	Details() map[string]interface{}
	// Reason 错误原因，定义为业务判定错误码
	Reason() string
}
