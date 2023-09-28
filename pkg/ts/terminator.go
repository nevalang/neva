package ts

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

func (r Terminator) ShouldTerminate(cur Trace, scope Scope) (bool, error) {
	return r.shouldTerminate(cur, scope, 0)
}

func (r Terminator) shouldTerminate(cur Trace, scope Scope, counter int) (bool, error) {
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
	// TODO refactor to use errors.Is(notfound), TODO call scope.Update()?
	if prevRef, err := scope.GetType(cur.prev.ref); err == nil {
		isRecursionAllowed = prevRef.IsRecursionAllowed
	}

	prev := cur.prev
	for prev != nil {
		if prev.ref != cur.ref {
			prev = prev.prev
			continue
		}

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

	return false, nil
}

// getLast3AndSwap turns [... a b a] into [b a b]
func (Terminator) getLast3AndSwap(cur Trace) Trace {
	t1 := Trace{prev: nil, ref: cur.prev.ref}
	t2 := Trace{prev: &t1, ref: cur.ref}
	return Trace{prev: &t2, ref: cur.prev.ref}
}
