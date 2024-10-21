package compiler

import (
	"fmt"

	src "github.com/nevalang/neva/internal/compiler/sourcecode"
	"github.com/nevalang/neva/internal/compiler/sourcecode/core"
)

type Error struct {
	Err      error
	Location *src.Location
	Range    *core.Meta
	child    *Error
}

func NewError(err error, meta *core.Meta, location *src.Location) *Error {
	return &Error{
		Err:      err,
		Range:    meta,
		Location: location,
	}
}

func (e Error) Wrap(child *Error) *Error {
	e.child = child
	return &e
}

// FIXME: it doesn't make sense to wrap if we don't use anything from parent
func (e Error) unwrap() Error {
	for e.child != nil {
		e = *e.child
	}
	return e
}

func (e Error) Error() string {
	e = e.unwrap()

	hasErr := e.Err != nil
	hasMeta := e.Range != nil
	hasLocation := e.Location != nil

	switch {
	case hasLocation && hasMeta:
		return fmt.Sprintf("%v:%v: %v", *e.Location, e.Range.Start, e.Err)
	case hasLocation:
		return fmt.Sprintf("%v: %v", *e.Location, e.Err)
	case hasMeta:
		return fmt.Sprintf("%v: %v", e.Range.Start, e.Err)
	case hasErr:
		return e.Err.Error()
	}

	panic(e)
}
