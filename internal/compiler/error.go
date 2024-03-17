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
		var loc *src.Location
		if e.child.Location != nil {
			loc = e.child.Location
		} else {
			loc = e.Location
		}

		var meta *core.Meta
		if e.child.Meta != nil {
			meta = e.child.Meta
		} else {
			meta = e.Meta
		}

		e = Error{
			Err:      e.child.Err,
			Location: loc,
			Meta:     meta,
			child:    e.child.child,
		}
	}

	return e
}

func (e Error) Error() string {
	e = e.unwrap()

	hasMeta := e.Meta != nil
	hasLocation := e.Location != nil
	hasErr := e.Err != nil

	if _, ok := e.Err.(*Error); ok {
		panic("internal error")
	}

	switch {
	case hasLocation && hasMeta:
		return fmt.Sprintf("%v:%v %v", *e.Location, e.Meta.Start, e.Err)
	case hasLocation:
		return fmt.Sprintf("%v %v", *e.Location, e.Err)
	case hasMeta:
		return fmt.Sprintf("%v %v", e.Meta.Start, e.Err)
	case hasErr:
		return e.Err.Error()
	}

	panic(e)
}
