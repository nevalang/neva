package types

import (
	"errors"
	"fmt"
)

var (
	ErrDirectRecursion   = errors.New("type definition's body must not be directly referenced to itself")
	ErrPrevDefNotFound   = errors.New("prev def not found")
	ErrIndirectRecursion = errors.New("type definition's body must not be indirectly referenced to itself")
	ErrSwapRun           = errors.New("couldn't do test run with swapped trace")
	ErrCounter           = errors.New("recursive calls counter limit exceeded")
)

type Terminator struct{}

func (r Terminator) ShouldTerminate(cur Trace, scope map[string]Def) (bool, error) {
	return r.shouldTerminate(cur, scope, 0)
}

func (r Terminator) shouldTerminate(cur Trace, scope map[string]Def, counter int) (bool, error) {
	if counter > 1 {
		return false, ErrCounter
	}

	if cur.prev == nil {
		return false, nil
	}

	if cur.ref == cur.prev.ref {
		return false, fmt.Errorf("%w: %v", ErrDirectRecursion, cur)
	}

	isRecursionAllowed := false
	if _, ok := scope[cur.prev.ref]; ok { // could not be there in case of type parameter
		isRecursionAllowed = scope[cur.prev.ref].IsRecursionAllowed
	}

	prev := cur.prev
	for prev != nil {
		if prev.ref == cur.ref {
			if isRecursionAllowed {
				return true, nil
			}

			isPrevRecursive, err := r.shouldTerminate(r.getLast3AndSwap(cur), scope, counter+1)
			if err != nil {
				if errors.Is(err, ErrCounter) || errors.Is(err, ErrIndirectRecursion) {
					return false, fmt.Errorf("%w: %v", ErrIndirectRecursion, cur)
				}
				return false, fmt.Errorf("%w: %v", ErrSwapRun, err)
			} else if isPrevRecursive {
				return true, nil
			}

			return false, errors.New("unknown")
		}

		prev = prev.prev
	}

	return false, nil
}

// getLast3AndSwap turns [... a b a] into [b a b]
func (Terminator) getLast3AndSwap(cur Trace) Trace {
	t1 := Trace{prev: nil, ref: cur.prev.ref}
	t2 := Trace{prev: &t1, ref: cur.ref}
	return Trace{prev: &t2, ref: cur.prev.ref}
}
