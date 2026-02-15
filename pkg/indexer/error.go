package indexer

import (
	"github.com/nevalang/neva/internal/compiler"
	"github.com/nevalang/neva/pkg/core"
)

// Error is an indexing error intended for tooling consumers (for example LSP).
//
// It wraps compiler errors internally while exposing stable diagnostic data.
type Error struct {
	Meta    *core.Meta
	cause   *compiler.Error
	Message string
}

// Error returns a human-readable error message with location context when available.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	if e.Meta == nil {
		return e.Message
	}
	return (&compiler.Error{
		Message: e.Message,
		Meta:    e.Meta,
	}).Error()
}

// Unwrap exposes the wrapped compiler error for internal integrations.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// wrapCompilerError adapts internal compiler errors to indexer.Error.
func wrapCompilerError(err *compiler.Error) *Error {
	if err == nil {
		return nil
	}

	deepest := err.Unwrap()
	return &Error{
		Message: deepest.Message,
		Meta:    deepest.Meta,
		cause:   deepest,
	}
}
