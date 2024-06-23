package appcmd

// invalidArgumentError is used to indicate that an error was
// caused by argument validation.
type invalidArgumentError struct {
	err error
}

func newInvalidArgumentError(err error) *invalidArgumentError {
	return &invalidArgumentError{err: err}
}

func (a *invalidArgumentError) Error() string {
	if a == nil {
		return ""
	}
	return a.err.Error()
}

func (a *invalidArgumentError) Unwrap() error {
	if a == nil {
		return nil
	}
	return a.err
}
