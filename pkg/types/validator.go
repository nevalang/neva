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
	if expr.Lit.Empty() == expr.Inst.Empty() { // we don't use expr.Empty() because constraint can be empty
		return ErrInvalidExprType
	}

	if expr.Lit.Empty() { // it's non-empty inst, nothing to validate
		return nil
	}

	switch expr.Lit.Type() { // by now we know it's not empty literal
	case RecLitType:
		return nil // nothing to check here, records with 0 fields are ok also
	case ArrLitType:
		if expr.Lit.Arr.Size < 2 {
			return fmt.Errorf("%w: got %d", ErrArrSize, expr.Lit.Arr.Size)
		}
	case UnionLitType:
		if l := len(expr.Lit.Union); l < 2 {
			return fmt.Errorf("%w: got %d", ErrUnionLen, l)
		}
	case EnumLitType:
		if l := len(expr.Lit.Enum); l < 2 {
			return fmt.Errorf("%w: got %d", ErrEnumLen, l)
		}
		set := make(map[string]struct{}, len(expr.Lit.Enum))
		for _, el := range expr.Lit.Enum { // look for duplicate
			if _, ok := set[el]; ok {
				return fmt.Errorf("%w: %s", ErrEnumDupl, el)
			}
			set[el] = struct{}{}
		}
	}

	return nil
}
