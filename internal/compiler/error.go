package compiler

import (
	"fmt"

	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Error struct {
	Message string
	Meta    *core.Meta
	child   *Error
}

func (e Error) Wrap(child *Error) *Error {
	e.child = child
	return &e
}

func (e Error) unwrap() *Error {
	for e.child != nil {
		e = *e.child
	}
	return &e
}

func (e *Error) Error() string {
	var s string

	current := e.unwrap()
	hasMeta := current.Meta != nil

	if hasMeta {
		s = fmt.Sprintf("%v:%v: %v", current.Meta.Location, current.Meta.Start, current.Message)
	} else {
		s = current.Message
	}

	return s
}
