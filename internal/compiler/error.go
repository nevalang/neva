package compiler

import (
	"fmt"
	"strings"

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

func (e *Error) Error() string {
	var builder strings.Builder

	current := e
	for current != nil {
		hasLocation := current.Location != nil
		hasMeta := current.Meta != nil

		if hasLocation && hasMeta {
			fmt.Fprintf(&builder, "%v:%v: %v\n", *current.Location, current.Meta.Start, current.Message)
		} else if hasLocation {
			fmt.Fprintf(&builder, "%v: %v\n", *current.Location, current.Message)
		} else if hasMeta {
			fmt.Fprintf(&builder, "%v: %v\n", current.Meta.Start, current.Message)
		} else {
			builder.WriteString(current.Message + "\n")
		}

		current = current.child
	}

	return builder.String()
}
