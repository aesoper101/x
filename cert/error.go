package cert

type InvalidTLSConfigError struct {
	msg string
}

func newInvalidTLSConfigError(msg string) error {
	return &InvalidTLSConfigError{msg: msg}
}

func (e *InvalidTLSConfigError) Error() string {
	return e.msg
}
