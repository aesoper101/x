package app

import (
	"errors"
	"fmt"
)

type appError struct {
	exitCode int
	err      error
}

func newAppError(exitCode int, err error) *appError {
	if exitCode == 0 {
		err = fmt.Errorf(
			"got invalid exit code %d when constructing appError (original error was %w)",
			exitCode,
			err,
		)
		exitCode = 1
	}
	if err == nil {
		err = errors.New("got nil error when constructing appError")
	}
	return &appError{
		exitCode: exitCode,
		err:      err,
	}
}

func (e *appError) Error() string {
	if e == nil {
		return ""
	}
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *appError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func printError(container StderrContainer, err error) {
	if errString := err.Error(); errString != "" {
		_, _ = fmt.Fprintln(container.Stderr(), errString)
	}
}
