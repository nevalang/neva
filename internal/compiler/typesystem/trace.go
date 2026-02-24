package typesystem

import (
	"strings"

	"github.com/nevalang/neva/pkg/core"
)

// Linked-list to handle recursive types
type Trace struct {
	prev *Trace
	cur  core.EntityRef
}

// O(2n)
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

func NewTrace(prev *Trace, cur core.EntityRef) Trace {
	return Trace{
		prev: prev,
		cur:  cur,
	}
}
