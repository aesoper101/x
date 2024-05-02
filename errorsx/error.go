package errorsx

type Error interface {
	GetParent() error
	GetMessage() string
	ResetMessage(message string)
	GetID() string
}
