package types

import (
	"errors"
	"fmt"
)

type RecursionChecker struct{}

// checkRecursion returns true and nil error for recursive expressions that should not go on next step of resolving.
// It returns false and nil err for non-recursive expressions with valid trace
// and false with non-nil err for bad recursion cases.
func (RecursionChecker) Check(cur Trace, scope map[string]Def) (bool, error) {
	if cur.prev == nil {
		return false, nil
	}

	if cur.v == cur.prev.v {
		return false, fmt.Errorf("%w: trace: %v", ErrDirectRecursion, cur)
	}

	prevDef, ok := scope[cur.prev.v]
	if !ok {
		return false, fmt.Errorf("%w: %v", errors.New("prev def not found"), cur)
	}

	isPrevAllowRecursion := prevDef.RecursionAllowed

	prev := cur.prev
	for prev != nil {
		if prev.v == cur.v {
			if isPrevAllowRecursion {
				return true, nil
			}
			return false, fmt.Errorf("%w: %v", ErrIndirectRecursion, cur)
		}
		prev = prev.prev
	}

	return false, nil
}
