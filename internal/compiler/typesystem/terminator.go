package typesystem

import (
	"errors"
	"fmt"

	"github.com/nevalang/neva/pkg/core"
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

	if sameRefs(cur.cur, cur.prev.cur) {
		return false, fmt.Errorf("%w: %v", ErrDirectRecursion, cur)
	}

	// Get prev ref's CanBeUsedForRecursiveDefinitions if it exists.
	// Note that we don't care if it's not found. Not all types are in the scope, some of them are in the frame.
	canBeUsedForRecursiveDefinitions := isRecursiveWrapper(cur.prev.cur)
	if prevRef, _, err := scope.GetType(cur.prev.cur); err == nil {
		// we don't have to check if prev has params, it has because we're here
		canBeUsedForRecursiveDefinitions = prevRef.BodyExpr == nil || isRecursiveWrapper(cur.prev.cur)
	}

	prev := cur.prev
	for prev != nil {
		if prev.cur != cur.cur {
			prev = prev.prev
			continue
		}

		if canBeUsedForRecursiveDefinitions {
			return true, nil
		}

		isPrevRecursive, err := r.shouldTerminate(r.getLast3AndSwap(cur), scope, counter+1)
		if err != nil {
			if errors.Is(err, ErrCounter) || errors.Is(err, ErrIndirectRecursion) {
				return false, fmt.Errorf("%w: %v", ErrIndirectRecursion, cur)
			}
			return false, fmt.Errorf("%w: %w", ErrSwapRun, err)
		} else if isPrevRecursive {
			return true, nil
		}

		return false, errors.New("unknown")
	}

	return false, nil
}

// getLast3AndSwap turns [... a b a] into [b a b]
func (Terminator) getLast3AndSwap(cur Trace) Trace {
	t1 := Trace{prev: nil, cur: cur.prev.cur}
	t2 := Trace{prev: &t1, cur: cur.cur}
	return Trace{prev: &t2, cur: cur.prev.cur}
}

func sameRefs(cur, prev core.EntityRef) bool {
	a := cur.String()
	b := prev.String()
	return a == b
}

func isRecursiveWrapper(ref core.EntityRef) bool {
	return ref.Name == "maybe" && (ref.Pkg == "" || ref.Pkg == "builtin")
}
