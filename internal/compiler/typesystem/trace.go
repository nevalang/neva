package typesystem

import (
	"strings"

	"github.com/nevalang/neva/pkg/core"
)

// Trace records the type references visited during recursive type resolution.
// cur is the newest reference; prev links to the reference that led to it.
//
// For example, while resolving the inner list in list<list<int>>, the trace is
// represented newest-first in memory as list -> list, while String prints it
// in source traversal order as [list, list].
type Trace struct {
	prev *Trace
	cur  core.EntityRef
}

// String prints a trace from the oldest reference to the newest reference.
func (t Trace) String() string {
	lastToFirst := []core.EntityRef{}

	tmp := &t
	for tmp != nil {
		lastToFirst = append(lastToFirst, tmp.cur)
		tmp = tmp.prev
	}

	var firstToLast strings.Builder
	firstToLast.WriteString("[")
	for i := len(lastToFirst) - 1; i >= 0; i-- {
		firstToLast.WriteString(lastToFirst[i].String())
		if i > 0 {
			firstToLast.WriteString(", ")
		}
	}

	return firstToLast.String() + "]"
}

// NewTrace appends cur to prev and returns the new end of the trace.
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func NewTrace(prev *Trace, cur core.EntityRef) Trace {
	return Trace{
		prev: prev,
		cur:  cur,
	}
}
