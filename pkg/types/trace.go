package types

// Linked-list to to handle recursive types
type Trace struct {
	prev *Trace
	ref  string
}

// O(2n)
func (t Trace) String() string {
	ss := []string{}

	var tmp *Trace = &t
	for tmp != nil {
		ss = append(ss, tmp.ref)
		tmp = tmp.prev
	}

	s := "["
	for i := len(ss) - 1; i >= 0; i-- {
		s += ss[i]
		if i > 0 {
			s += ", "
		}
	}

	return s + "]"
}

func NewTrace(prev *Trace, v string) Trace {
	return Trace{
		prev: prev,
		ref:  v,
	}
}
