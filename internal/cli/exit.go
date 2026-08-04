package cli

import "errors"

// exitError carries a process status without exposing the command framework to
// cmd/neva.
type exitError struct {
	cause error
	code  int
}

// Error returns the underlying message, if any.
func (err *exitError) Error() string {
	if err.cause == nil {
		return ""
	}
	return err.cause.Error()
}

// Unwrap exposes the underlying command error.
func (err *exitError) Unwrap() error {
	return err.cause
}

// exitWithCode returns an error that makes the command exit with code.
func exitWithCode(cause error, code int) error {
	return &exitError{cause: cause, code: code}
}

// ExitCode returns the process status associated with err, or one for an
// ordinary command error.
func ExitCode(err error) int {
	var commandExit *exitError
	if errors.As(err, &commandExit) {
		return commandExit.code
	}
	return 1
}
