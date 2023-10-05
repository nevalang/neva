package typesystem

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

	// Get prev ref's CanBeUsedForRecursiveDefinitions if it exists.
	// Note that we don't care if it's not found. Not all types are in the scope, some of them are in the frame.
	var canBeUsedForRecursiveDefinitions bool
	if prevRef, _, err := scope.GetType(cur.prev.ref); err == nil {
		canBeUsedForRecursiveDefinitions = prevRef.CanBeUsedForRecursiveDefinitions
	}

	prev := cur.prev
	for prev != nil {
		if prev.ref != cur.ref {
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
