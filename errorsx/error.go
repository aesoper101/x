package errorsx

// Error 错误接口
// inspire by https://cloud.google.com/apis/design/errors?hl=zh-cn
type Error interface {
	GetParent() error
	GetMessage() string
	ResetMessage(message string)
	SetDetails(details map[string]string)
	GetDetails() map[string]string
	GetID() string
	HttpStatusCode() int
}
