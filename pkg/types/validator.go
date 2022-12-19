package types

import (
	"errors"
	"fmt"
)

type Validator struct{}

var (
	ErrInvalidExprType = errors.New("expr must be ether literal or instantiation, not both and not neither")
	ErrUnknownLit      = errors.New("expr literal must be known")
	ErrArrSize         = errors.New("arr size must be >= 2")
	ErrEnumLen         = errors.New("enum len must be >= 2")
	ErrUnionLen        = errors.New("union len must be >= 2")
	ErrEnumDupl        = errors.New("enum contains duplicate elements")
)

// Checks that expression is either an instantiation or a literal, not both and not neither.
// All instantiations are valid because it's ok to have "" type with 0 arguments.
// For arr, union and enum it checks that size is >= 2. For enum it also ensures there no duplicate elements.
func (v Validator) Validate(expr Expr) error {
	if expr.Lit.Empty() == expr.Inst.Empty() {
		return ErrInvalidExprType
	}

	if expr.Lit.Empty() { // if it's inst
		return nil // then nothing to validate, resolving needed
	} // by now we know it's not empty lit

	// we don't check recs and empty lits
	switch expr.Lit.Type() { // because we know lit isn't empty and because empty recs are fine
	case ArrLitType:
		if expr.Lit.ArrLit.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.ArrLit.Size)
		}
	case UnionLitType:
		if l := len(expr.Lit.UnionLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
	case EnumLitType:
		if l := len(expr.Lit.EnumLit); l < 2 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
		set := make(map[string]struct{}, len(expr.Lit.EnumLit))
		for _, el := range expr.Lit.EnumLit { // look for duplicate
			if _, ok := set[el]; ok {
				return fmt.Errorf("%w: %s", ErrEnumDupl, el)
			}
			set[el] = struct{}{}
		}
	}

	return nil // valid lit
}
