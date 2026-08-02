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

// Terminator decides whether type resolution should continue expanding the
// current expression, stop at a valid recursive back-edge, or reject an
// invalid recursive definition.
type Terminator struct{}

// ShouldTerminate classifies the newest reference in a type-resolution trace.
// Its results have three meanings:
//
//   - false, nil: this is finite nesting; continue normal resolution;
//   - true, nil: this is a permitted recursive back-edge; keep the current type
//     reference without expanding it again;
//   - false, err: the type contains invalid direct or indirect recursion.
//
// For example, list<list<int>> returns false because both list occurrences are
// finite base-type instantiations. A recursive type such as
// `type Node struct { children list<Node> }` eventually returns true for the
// second Node because the path returns through list. In contrast, `type T T`
// returns ErrDirectRecursion because expanding T immediately produces T again.
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (r Terminator) ShouldTerminate(cur Trace, scope Scope) (bool, error) {
	return r.shouldTerminate(cur, scope, 0)
}

// shouldTerminate performs ShouldTerminate's bounded secondary check for
// indirect aliases. swapChecks prevents that internal re-check from recursing
// indefinitely; it does not represent the depth of the user's type.
//
//nolint:cyclop,gocognit,gocritic,gocyclo // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (r Terminator) shouldTerminate(cur Trace, scope Scope, swapChecks int) (bool, error) {
	if swapChecks > 1 {
		return false, ErrCounter
	}

	// A one-element trace cannot contain a cycle.
	if cur.prev == nil {
		return false, nil
	}

	if sameRefs(cur.cur, cur.prev.cur) {
		// The same reference may appear twice when a base type is nested. For
		// list<list<int>>, resolution visits the outer list and then the inner
		// list, producing [list, list]. A base type has no body to expand, so this
		// is finite nesting and resolution must continue into the inner list's
		// int argument.
		//
		// A user-defined alias does have a body. With `type T T`, resolving the
		// body immediately produces [T, T], which can never terminate. Only a
		// bodyless definition found in scope proves the finite base-type case;
		// every other adjacent repetition is rejected as direct recursion.
		if def, _, err := scope.GetType(cur.cur); err == nil && def.BodyExpr == nil {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", ErrDirectRecursion, cur)
	}

	// A repeated reference is valid only when the path back to it crosses a
	// type that can contain the recursive value without immediately expanding
	// it. For `type Node struct { children list<Node> }`, the relevant trace is
	// [Node, struct, list, Node], and list is that boundary.
	//
	// A reference may come from the resolver's generic frame instead of scope,
	// so a failed scope lookup is not itself an error here. maybe is recognized
	// explicitly because it is the built-in recursive wrapper.
	previousAllowsRecursion := isRecursiveWrapper(cur.prev.cur)
	if prevRef, _, err := scope.GetType(cur.prev.cur); err == nil {
		previousAllowsRecursion = prevRef.BodyExpr == nil || isRecursiveWrapper(cur.prev.cur)
	}

	// Search the older part of the trace for the current reference. Finding it
	// means that resolution has closed a cycle.
	prev := cur.prev
	for prev != nil {
		if prev.cur != cur.cur {
			prev = prev.prev
			continue
		}

		if previousAllowsRecursion {
			return true, nil
		}

		// No recursive boundary appears immediately before the repeated
		// reference. Check the final aliases from the opposite entry point, so
		// t1 -> t2 -> t1 cannot be accepted merely because resolution started at
		// t1 instead of t2.
		swappedTerminates, err := r.shouldTerminate(
			r.getLast3AndSwap(cur),
			scope,
			swapChecks+1,
		)
		if err != nil {
			if errors.Is(err, ErrCounter) || errors.Is(err, ErrIndirectRecursion) {
				return false, fmt.Errorf("%w: %v", ErrIndirectRecursion, cur)
			}
			return false, fmt.Errorf("%w: %w", ErrSwapRun, err)
		}
		if swappedTerminates {
			return true, nil
		}

		return false, fmt.Errorf("could not classify recursive trace %v", cur)
	}

	return false, nil
}

// getLast3AndSwap turns [... a, b, a] into [b, a, b]. This asks the same
// indirect-recursion question starting from the other alias. For example,
// [t1, t2, t1] becomes [t2, t1, t2].
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func (Terminator) getLast3AndSwap(cur Trace) Trace {
	t1 := Trace{prev: nil, cur: cur.prev.cur}
	t2 := Trace{prev: &t1, cur: cur.cur}
	return Trace{prev: &t2, cur: cur.prev.cur}
}

//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func sameRefs(cur, prev core.EntityRef) bool {
	a := cur.String()
	b := prev.String()
	return a == b
}

// isRecursiveWrapper recognizes the built-in recursive wrapper even when its
// definition is supplied through the resolver's generic frame rather than the
// scope passed to Terminator.
//
//nolint:gocritic // TODO(strict-lint phase 1): temporary suppression; remove after strict cleanup.
func isRecursiveWrapper(ref core.EntityRef) bool {
	return ref.Name == "maybe" && (ref.Pkg == "" || ref.Pkg == "builtin")
}
