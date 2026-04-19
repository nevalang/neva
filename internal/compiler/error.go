//nolint:all // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
package compiler

import (
	"fmt"

	"github.com/nevalang/neva/pkg/core"
)

type Error struct {
	Meta    *core.Meta
	child   *Error
	Message string
}

func (e Error) Wrap(child *Error) *Error {
	e.child = child
	return &e
}

func (e Error) Unwrap() *Error {
	for e.child != nil {
		e = *e.child
	}
	return &e
}

func (e *Error) Error() string {
	var s string

	current := e.Unwrap()
	hasMeta := current.Meta != nil

	if hasMeta {
		s = fmt.Sprintf("%v:%v: %v", current.Meta.Location, current.Meta.Start, current.Message)
	} else {
		s = current.Message
	}

	return s
}
