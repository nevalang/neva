package typesystem

import "fmt"

// Linked-list to handle recursive types
type Trace struct {
	prev *Trace
	ref  fmt.Stringer
}

// O(2n)
func (t Trace) String() string {
	lastToFirst := []fmt.Stringer{}

	tmp := &t
	for tmp != nil {
		lastToFirst = append(lastToFirst, tmp.ref)
		tmp = tmp.prev
	}

	firstToLast := "["
	for i := len(lastToFirst) - 1; i >= 0; i-- {
		if lastToFirst[i] == nil { // every trace starts with nil
			continue
		}
		firstToLast += lastToFirst[i].String()
		if i > 0 {
			firstToLast += ", "
		}
	}

	return firstToLast + "]"
}

func NewTrace(prev *Trace, v fmt.Stringer) Trace {
	return Trace{
		prev: prev,
		ref:  v,
	}
}
