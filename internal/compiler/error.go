package compiler

import (
	"fmt"

	src "github.com/nevalang/neva/pkg/sourcecode"
)

type Error struct {
	Err      error
	Location *src.Location
	Meta     *src.Meta
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

	hasMeta := e.Meta != nil
	hasLocation := e.Location != nil
	hasErr := e.Err != nil

	if _, ok := e.Err.(*Error); ok {
		panic("")
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
