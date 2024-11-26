package compiler

import (
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Error struct {
	Message  string
	Location *src.Location
	Meta     *core.Meta

	child *Error
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
	hasLocation := current.Location != nil
	hasMeta := current.Meta != nil

	if hasLocation && hasMeta {
		s = fmt.Sprintf("%v:%v: %v", *current.Location, current.Meta.Start, current.Message)
	} else if hasLocation {
		s = fmt.Sprintf("%v: %v", *current.Location, current.Message)
	} else if hasMeta {
		s = fmt.Sprintf("%v: %v", current.Meta.Start, current.Message)
	} else {
		s = current.Message
	}

	return s
}
