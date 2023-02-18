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
)

type RecursionTerminator struct{}

func (r RecursionTerminator) ShouldTerminate(cur Trace, scope map[string]Def) (bool, error) {
	fmt.Println(cur)

	if cur.prev == nil {
		return false, nil
	}

	if cur.ref == cur.prev.ref {
		return false, fmt.Errorf("%w: %v", ErrDirectRecursion, cur)
	}

	prevDef, ok := scope[cur.prev.ref]
	if !ok {
		return false, fmt.Errorf("%w: %v", ErrPrevDefNotFound, cur)
	}

	fmt.Println(prevDef.BodyExpr)

	isPrevAllowRecursion := prevDef.IsRecursionAllowed

	prev := cur.prev
	for prev != nil {
		if prev.ref == cur.ref { // same ref found, that's a loop
			if isPrevAllowRecursion { // but that's ok if prev ref allow recursion
				return true, nil
			}

			// or maybe prev ref is recursive itself
			isPrevRecursive, err := r.ShouldTerminate(r.swapTrace(cur), scope)
			if err != nil {
				return false, fmt.Errorf("%w: %v", ErrSwapRun, err)
			} else if isPrevRecursive {
				return true, nil
			}

			return false, fmt.Errorf("%w: %v", ErrIndirectRecursion, cur)
		}

		prev = prev.prev
	}

	return false, nil
}

// swapTrace turns [... a b a] into [b a b]
func (RecursionTerminator) swapTrace(cur Trace) Trace {
	t1 := Trace{prev: nil, ref: cur.prev.ref}
	t2 := Trace{prev: &t1, ref: cur.ref}
	return Trace{prev: &t2, ref: cur.prev.ref}
}
