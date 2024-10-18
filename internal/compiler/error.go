package compiler

import (
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Error struct {
	Err      error
	Location *src.Location
	Meta     *core.Meta
	child    *Error
}

func (e Error) Wrap(child *Error) *Error {
	e.child = child
	return &e
}

func (e Error) unwrap() Error {
	for e.child != nil {
		e = *e.child
	}
	return e
}

func (e Error) Error() string {
	e = e.unwrap()

	hasErr := e.Err != nil
	hasMeta := e.Meta != nil
	hasLocation := e.Location != nil

	switch {
	case hasLocation && hasMeta:
		return fmt.Sprintf("%v:%v: %v", *e.Location, e.Meta.Start, e.Err)
	case hasLocation:
		return fmt.Sprintf("%v: %v", *e.Location, e.Err)
	case hasMeta:
		return fmt.Sprintf("%v: %v", e.Meta.Start, e.Err)
	case hasErr:
		return e.Err.Error()
	}

	panic(e)
}
