package typesystem

// Linked-list to handle recursive types
type Trace struct {
	prev *Trace
	ref  string
}

// O(2n)
func (t Trace) String() string {
	lastToFirst := []string{}

	tmp := &t
	for tmp != nil {
		lastToFirst = append(lastToFirst, tmp.ref)
		tmp = tmp.prev
	}

	firstToLast := "["
	for i := len(lastToFirst) - 1; i >= 0; i-- {
		firstToLast += lastToFirst[i]
		if i > 0 {
			firstToLast += ", "
		}
	}

	return firstToLast + "]"
}

func NewTrace(prev *Trace, v string) Trace {
	return Trace{
		prev: prev,
		ref:  v,
	}
}
